package database

import (
	"fmt"

	"gobackup/helper"
)

type Database interface {
	GetDatabases() ([]string, error)
	GetTables(database string) ([]string, error)
	Backup(database string, excludeTables []string) (string, error)
	BackupAll() (string, error)
}

type BaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Args     []string
}

func GetDatabase(t string, h *helper.Helper, commandPath string, config BaseConfig) (Database, error) {
	switch t {
	case "mysql":
		return NewMysql(h, config.Host, config.Port, config.Username, config.Password, commandPath, config.Args)
	case "pg":
		return NewPostgresql(h, config.Host, config.Port, config.Username, config.Password, commandPath, config.Args), nil
	// case "redis":
	// 	return NewRedis(config.Command, config.Host, config.Port, config.Username, config.Password, config.Args), nil
	// case "mongo":
	// 	return NewMongo(config.Command, config.Host, config.Port, config.Username, config.Password, config.Args), nil
	// case "scp":
	// 	return NewScp(config.Host, config.Port, config.Username, config.Password, config.Pwd, config.SrcDir)
	default:
		return nil, fmt.Errorf("unknown type")
	}
}
