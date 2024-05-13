package main

import (
	_ "embed"
	"flag"
	"fmt"
	"gobackup/config"
	"gobackup/utils"
	"os"
	"path/filepath"

	"gobackup/dbs"

	"github.com/duke-git/lancet/v2/fileutil"
)

//go:embed tools.zip
var dumpTools []byte

// 压缩包名
const zipname = "tools.zip"

// 提取压缩包
func ExtractZip() error {
	name := "tools.zip"
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(dumpTools); err != nil {
		return err
	}
	return err
}

// 解压工具
func UnzipTool() error {
	dirname := "bin"
	// zipname := "tools.zip"
	if fileutil.IsExist(dirname) {
		return nil
	}

	if err := ExtractZip(); err != nil {
		return err
	}

	err := fileutil.UnZip(zipname, ".")
	return err
}

// 导出配置
func ExportConfig(path string) error {
	if fileutil.IsExist(path) {
		return nil
	}
	CreateConfig(path)
	return os.ErrNotExist
}

func CreateConfig(path string) error {
	config := `level: "debug"
store: "/app/data"
cron: "* * * * *"
backup:
- type: mysql
  host: 192.168.52.147
  port: 3306
  user: root
  password: root
- type: pg
  host: 192.168.52.147
  port: 5432
  user: postgres
  password: root
- type: mongo
  host: 192.168.52.147
  port: 27017
  user: ""
  password: ""
- type: scp
  user: root 
  password: root
  host: 192.168.52.147
  port: 22
  fromPath: /tmp/scp
`
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(config))
	return err
}

// 创建存储目录
func CreateStoreDir(path string) error {
	if fileutil.IsExist(path) {
		return nil
	}
	return os.Mkdir(path, os.ModePerm)
}

func Run(cfg config.GobackupConfig) {
	helper := utils.NewHelper(cfg.Store)

	storePath := cfg.Store

	// 需要压缩的文件名
	zipNames := make([]string, 0)

	// 数据库配置列表(mysql, mongo, postgres)
	dbConfigs := make([]config.ConfigDatabase, 0)
	for _, db := range cfg.ConfigDatabase {
		if db.DBType == dbs.MONGO || db.DBType == dbs.MYSQL || db.DBType == dbs.POSTGRES {
			dbConfigs = append(dbConfigs, db)
		}
	}

	// 工具配置列表(scp)
	toolConfigs := make([]config.ConfigDatabase, 0)
	for _, db := range cfg.ConfigDatabase {
		if db.DBType == dbs.SCP {
			toolConfigs = append(toolConfigs, db)
		}
	}

	// 数据库备份
	for _, db := range dbConfigs {
		dbType := db.DBType
		dbHost := db.Host
		dbPort := db.Port
		dbUser := db.User
		dbPasswd := db.Password

		// 获取数据库实例
		client := dbs.GetDatabase(dbType, dbUser, dbPasswd, dbHost, dbPort)
		backupCommands, err := client.BackupAll()
		if err != nil {
			utils.Error("backup error", err)
			continue
		}

		// 获取当前程序路径
		pwd, err := os.Getwd()
		if err != nil {
			utils.Error("get pwd error", err)
			return
		}

		// 根据系统, 获取命令行工具
		toolDir := filepath.Join(pwd, "bin")
		command, err := dbs.GetCommand(dbType, toolDir, dbUser, dbPasswd, dbHost, dbPort)
		if err != nil {
			utils.Error("get command error", err)
			continue
		}

		for _, dumpArgs := range backupCommands {
			_, err := helper.Exec(command, dumpArgs...)
			if err != nil {
				utils.Error("exec error", err)
				continue
			}

			// 增加到压缩列表
			name := filepath.Join(storePath, dumpArgs[len(dumpArgs)-1])
			zipNames = append(zipNames, name)
		}
	}

	// scp备份
	for _, db := range toolConfigs {
		dbType := db.DBType
		dbHost := db.Host
		dbPort := db.Port
		dbUser := db.User
		dbPasswd := db.Password
		tool := dbs.GetTools(dbType, dbUser, dbPasswd, dbHost, dbPort)
		name, err := tool.Backup(db.FromPath, storePath)
		if err != nil {
			utils.Error("scp error", err)
			continue
		}
		zipNames = append(zipNames, name)
	}

	// 压缩
	if len(zipNames) == 0 {
		return
	}
	if err := utils.Compressor(zipNames, cfg.Store); err != nil {
		utils.Error("compress error", err)
		return
	}

	// 删
	for _, filename := range zipNames {
		if filename == storePath {
			panic("无法删存储目录")
		}
		if err := os.RemoveAll(filename); err != nil {
			utils.Error("delete error", err)
		}
	}
}

func main() {
	// 解压工具
	if err := UnzipTool(); err != nil {
		fmt.Println("解压工具失败", err)
		return
	}

	cfgPath := "config.yaml"
	// 导出配置
	if err := ExportConfig(cfgPath); err != nil {
		fmt.Println("重新导出配置文件", err)
		return
	}

	// 读取配置
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		fmt.Println("读取配置失败", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println("校验配置信息失败", err)
		return
	}

	// 创建存储目录
	if err := CreateStoreDir(cfg.Store); err != nil {
		fmt.Println("创建存储目录失败", err)
		return
	}

	// 命令行
	if cmd(cfg) {
		return
	}

	con := utils.NewCronTask()
	con.Add(cfg.Cron, func() {
		Run(cfg)
	})
	con.Start()
	select {}
}

func cmd(cfg config.GobackupConfig) bool {
	run := flag.Bool("run", false, "执行任务")

	flag.Parse()

	if *run {
		Run(cfg)
		return true
	}
	return false
}
