package dbs

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "pg"
	MONGO    = "mongo"
)

type Database interface {
	BackupAll() ([][]string, error)
}

func GetDatabase(dbType string, Username, Password, Host, Port string) Database {
	switch dbType {
	case MYSQL:
		return NewMysql(Host, Port, Username, Password)
	case POSTGRES:
		return NewPostgresql(Host, Port, Username, Password)
	case MONGO:
		return NewMongo(Host, Port, Username, Password)
	// case "redis":
	// 	return NewRedis(config.Command, config.Host, config.Port, config.Username, config.Password, config.Args), nil
	// case "scp":
	// 	return NewScp(config.Host, config.Port, config.Username, config.Password, config.Pwd, config.SrcDir)
	default:
		return nil
	}
}

func GetCommand(dbType string, username, passwd, host, port string) (string, error) {
	if runtime.GOOS == "windows" {
		switch dbType {
		case MYSQL:
			version, err := version(username, passwd, host, port)
			if err != nil {
				return "", fmt.Errorf("无法获取到mysql版本")
			}
			if strings.HasPrefix(version, "5") {
				return "mysqldump.exe", nil
			}
			return "mysqldump8.exe", nil
		case POSTGRES:
			return "pg_dump.exe", nil
		case MONGO:
			return "mongodump.exe", nil
		}
	}
	if runtime.GOOS == "linux" {
		switch dbType {
		case MYSQL:
			return "mysqldump", nil
		case POSTGRES:
			return "pg_dump", nil
		case MONGO:
			return "mongodump", nil
		}
	}
	return "", fmt.Errorf("不支持的系统类型")
}
