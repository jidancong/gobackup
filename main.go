package main

import (
	_ "embed"
	"flag"
	"fmt"
	"gobackup/config"
	"gobackup/controller"
	"gobackup/helper"
	"gobackup/utils"
	"os"
	"path/filepath"
)

//go:embed tools.rar
var dumpTools []byte

func main() {
	initConfig := flag.Bool("config", false, "生成默认配置")
	run := flag.Bool("run", false, "执行定时任务")

	flag.Parse()
	if *initConfig {
		// 创建初始化配置文件
		configName := "config.yaml"
		if err := CreateConfig(configName); err != nil {
			fmt.Printf("初始化配置失败 error: %v", err)
			return
		}
	}
	if *run {
		Run()
	}

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func Run() {
	// 实例配置
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始化日志
	utils.NewSlog(cfg.Level)
	// 实例压缩
	archive := helper.NewArchive(cfg.Store)
	// 实例cmd
	h := helper.NewHelper(cfg.Store)
	// 实例定时任务
	c := helper.NewCronTask()

	// 创建数据存储目录
	if _, err := os.Stat(cfg.Store); os.IsNotExist(err) {
		os.Mkdir(cfg.Store, os.ModePerm)
	}

	// 每次都对bin目录删
	toolPath := filepath.Join(cfg.Store, "bin")
	if err := os.RemoveAll(toolPath); err != nil {
		fmt.Println("删bin目录失败 error: ", err)
		return
	}
	if err := upload(); err != nil {
		fmt.Println("上传压缩包失败 error: ", err)
		return
	}

	if err := archive.DeCompressor("tools-tmp.rar"); err != nil {
		fmt.Println("解压失败 error: ", err)
		return
	}

	// 实例控制器
	con := controller.NewController(cfg, cfg.Store, toolPath, h, archive, c)
	// 启动
	if err := con.Start(); err != nil {
		fmt.Printf("new controller error: %v", err)
		return
	}

	select {}

}

func upload() error {
	file, err := os.Create("tools-tmp.rar")
	if err != nil {
		return err
	}

	if _, err := file.Write(dumpTools); err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func CreateConfig(fileName string) error {
	config := `level: "debug"
store: "/app/data"
cron: "* * * * *"
backup:
- type: mysql
  host: 192.168.52.146
  port: 3306
  user: root
  password: root
- type: pg
  host: 192.168.52.146
  port: 5432
  user: postgres
  password: root
- type: mongo
  host: 192.168.52.146
  port: 27017
  user: ""
  password: ""
`
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(config))
	return err
}
