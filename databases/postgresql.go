package database

import (
	"fmt"
	"gobackup/helper"
	"gobackup/utils"
	"os"
	"path"
	"runtime"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase struct {
	DatabaseName string `gorm:"column:datname"`
}

type Postgresql struct {
	h        *helper.Helper
	host     string
	port     string
	username string
	password string
	args     []string
	command  string // 命令二进制工具
}

func NewPostgresql(h *helper.Helper, host, port, username, password string, commandPath string, args []string) *Postgresql {
	os.Setenv("PGPASSWORD", password)
	command := "pg_dump"
	if runtime.GOOS == "windows" {
		// command = "pg_dump.exe"
		command = path.Join(commandPath, "pg_dump.exe")
	}
	// commandFull := path.Join(commandPath, command)
	return &Postgresql{
		h:        h,
		host:     host,
		port:     port,
		username: username,
		password: password,
		args:     args,
		command:  command,
	}
}

func (db *Postgresql) Backup(database string, excludeTables []string) (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "-h", db.host)
	dumpArgs = append(dumpArgs, "-p", db.port)
	dumpArgs = append(dumpArgs, "-U", db.username)
	dumpArgs = append(dumpArgs, "-d", database)

	for _, table := range excludeTables {
		dumpArgs = append(dumpArgs, "-T", table)
	}

	sqlFile := fmt.Sprintf("%s-%s-%s-%s-%s.sql", "pgsql", database, db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-f", sqlFile)

	_, err := db.h.Exec(db.command, dumpArgs...)
	return sqlFile, err
}

func (db *Postgresql) BackupAll() (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "-h", db.host)
	dumpArgs = append(dumpArgs, "-p", db.port)
	dumpArgs = append(dumpArgs, "-U", db.username)

	sqlfile := fmt.Sprintf("%s-%s-%s-%s.sql", "pgsql", db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-f", sqlfile)

	_, err := db.h.Exec(db.command, dumpArgs...)
	return sqlfile, err
}

func (db *Postgresql) GetDatabases() ([]string, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=template1 port=%s sslmode=disable TimeZone=Asia/Shanghai", db.host, db.username, db.password, db.port)
	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var databases []PgDatabase
	if err := gormdb.Raw("SELECT datname FROM pg_database WHERE datistemplate = false").Scan(&databases).Error; err != nil {
		panic(err)
	}

	var ds []string
	for _, v := range databases {
		ds = append(ds, v.DatabaseName)
	}

	return ds, nil
}

func (db *Postgresql) GetTables(database string) ([]string, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", db.host, db.username, db.password, db.port, database)
	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	tables, err := gormdb.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}
