package database

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"time"

	"gobackup/helper"
	"gobackup/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	h        *helper.Helper
	host     string
	port     string
	username string
	password string
	args     []string
	command  string // 命令二进制程序
}

func NewMongo(h *helper.Helper, host, port, username, password, commandPath string, args []string) *Mongo {
	command := "mongodump"
	if runtime.GOOS == "windows" {
		command = path.Join(commandPath, "mongodump.exe")
	}
	return &Mongo{
		h:        h,
		host:     host,
		port:     port,
		username: username,
		password: password,
		args:     args,
		command:  command,
	}
}

func (db *Mongo) BackupAll() (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	if len(db.username) > 0 {
		dumpArgs = append(dumpArgs, "--username", db.username)
	}
	if len(db.password) > 0 {
		dumpArgs = append(dumpArgs, "--password", db.password)
	}
	dumpArgs = append(dumpArgs, "--gzip")

	sqlfile := fmt.Sprintf("%s-%s-%s-%s", "mongo", db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "--out", sqlfile)

	_, err := db.h.Exec(db.command, dumpArgs...)
	return sqlfile, err
}

func (db *Mongo) GetDatabases() ([]string, error) {
	client, err := getMongoClient(db.host, db.port)
	if err != nil {
		return nil, err
	}

	names, err := client.ListDatabaseNames(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (db *Mongo) GetTables(database string) ([]string, error) {
	return nil, nil
}

func (db *Mongo) Backup(dbName string, excludeTables []string) (string, error) {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	if len(db.username) > 0 {
		dumpArgs = append(dumpArgs, "--username", db.username)
	}
	if len(db.password) > 0 {
		dumpArgs = append(dumpArgs, "--password", db.password)
	}
	if len(dbName) > 0 {
		dumpArgs = append(dumpArgs, "--db", dbName)
	}
	dumpArgs = append(dumpArgs, "--gzip")

	sqlfile := fmt.Sprintf("%s-%s-%s-%s", "mongo", db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "--out", sqlfile)

	_, err := db.h.Exec(db.command, dumpArgs...)
	return sqlfile, err
}

func getMongoClient(host, port string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hostAndPort := fmt.Sprintf("mongodb://%s:%s", host, port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(hostAndPort))
	if err != nil {
		return nil, err
	}
	return client, nil
}
