package database

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"gobackup/helper"
	"gobackup/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDatabase struct {
	DatabaseName string `gorm:"column:Database"`
}

type Mysql struct {
	h        *helper.Helper
	host     string
	port     string
	username string
	password string
	args     []string
	command  string // 命令二进制程序
}

func NewMysql(h *helper.Helper, host, port, username, password, commandPath string, args []string) (*Mysql, error) {
	os.Setenv("MYSQL_PWD", password)
	command := "mysqldump"
	if runtime.GOOS == "windows" {
		gormdb, err := getMysqlClient(username, password, host, port)
		if err != nil {
			return nil, err
		}

		var version string
		if err = gormdb.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
			return nil, err
		}
		if strings.HasPrefix(version, "5") {
			command = path.Join(commandPath, "mysqldump.exe")
		} else {
			command = path.Join(commandPath, "mysqldump8.exe")
		}
	}
	if len(args) <= 0 {
		args = []string{"--single-transaction", "--quick"}
	}
	return &Mysql{
		h:        h,
		host:     host,
		port:     port,
		username: username,
		password: password,
		args:     args,
		command:  command,
	}, nil
}

func (db *Mysql) BackupAll() (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	dumpArgs = append(dumpArgs, "-u", db.username)
	// dumpArgs = append(dumpArgs, `-p`+db.password)
	dumpArgs = append(dumpArgs, "--all-databases")

	sqlfile := fmt.Sprintf("%s-%s-%s-%s.sql", "mysql", db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-r", sqlfile)

	dumpArgs = append(dumpArgs, db.args...)
	_, err := db.h.Exec(db.command, dumpArgs...)
	return sqlfile, err
}

func (db *Mysql) Backup(database string, excludeTables []string) (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	dumpArgs = append(dumpArgs, "-u", db.username)
	// dumpArgs = append(dumpArgs, `-p`+db.password)

	dumpArgs = append(dumpArgs, database)
	for _, table := range excludeTables {
		dumpArgs = append(dumpArgs, "--ignore-table="+database+"."+table)
	}

	// sql file name
	// render_farm-127.0.0.1-20230901031418-4b7d5e34f3034195955e186e43270a45.sql
	fileName := fmt.Sprintf("%s-%s-%s-%s-%s.sql", "mysql", database, db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-r", fileName)

	dumpArgs = append(dumpArgs, db.args...)

	_, err := db.h.Exec(db.command, dumpArgs...)
	return fileName, err
}

func (db *Mysql) GetDatabases() ([]string, error) {
	gormdb, err := getMysqlClient(db.username, db.password, db.host, db.port)
	if err != nil {
		return nil, err
	}
	var databases []MysqlDatabase
	if err := gormdb.Raw("SHOW DATABASES").Scan(&databases).Error; err != nil {
		return nil, err
	}

	var ds []string
	for _, v := range databases {
		// 过滤information_schema表
		if v.DatabaseName == "information_schema" {
			continue
		}
		ds = append(ds, v.DatabaseName)
	}

	return ds, err
}

func (db *Mysql) GetTables(database string) ([]string, error) {
	gormdb, err := getMysqlClient(db.username, db.password, db.host, db.port)
	if err != nil {
		return nil, err
	}
	tables, err := gormdb.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (db *Mysql) Version() (string, error) {
	gormdb, err := getMysqlClient(db.username, db.password, db.host, db.port)
	if err != nil {
		return "", err
	}

	var version string
	err = gormdb.Raw("SELECT VERSION()").Scan(&version).Error
	return version, err
}

func getMysqlClient(username, password, host, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormdb, nil
}
