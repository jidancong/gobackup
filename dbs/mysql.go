package dbs

import (
	"fmt"
	"os"
	"slices"

	"gobackup/utils"
)

type Mysql struct {
	host     string
	port     string
	username string
	password string
}

func NewMysql(host, port, username, password string) *Mysql {
	os.Setenv("MYSQL_PWD", password)
	return &Mysql{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (db *Mysql) BackupAll() ([][]string, error) {
	dbnames, err := getDBs(db.username, db.password, db.host, db.port)
	if err != nil {
		return nil, err
	}

	// 不备份数据库
	dbnames = slices.DeleteFunc(dbnames, func(name string) bool {
		if name == "performance_schema" {
			return true
		}
		if name == "information_schema" {
			return true
		}
		if name == "mysql" {
			return true
		}
		if name == "sys" {
			return true
		}
		return false
	})

	execCommand := [][]string{}
	for _, dbname := range dbnames {
		execCommand = append(execCommand, db.backup(dbname))
	}
	return execCommand, nil
}

func (db *Mysql) backup(database string) []string {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	dumpArgs = append(dumpArgs, "-u", db.username)

	dumpArgs = append(dumpArgs, database)
	dumpArgs = append(dumpArgs, "--single-transaction", "--quick")
	// sql file name
	// render_farm-127.0.0.1-20230901031418-4b7d5e34f3034195955e186e43270a45.sql
	fileName := fmt.Sprintf("%s-%s-%s-%s-%s.sql", "mysql", database, db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-r", fileName)
	return dumpArgs
}
