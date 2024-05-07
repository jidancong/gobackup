package dbs

import (
	"fmt"
	"slices"

	"gobackup/utils"
)

type Mongo struct {
	host     string
	port     string
	username string
	password string
}

func NewMongo(host, port, username, password string) *Mongo {
	return &Mongo{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (db *Mongo) BackupAll() ([][]string, error) {
	execCommand := [][]string{}
	names, err := getMongoDBs(db.username, db.password, db.host, db.port)
	if err != nil {
		return nil, err
	}

	// 过滤
	names = slices.DeleteFunc(names, func(name string) bool { return name == "local" })

	for _, name := range names {
		execCommand = append(execCommand, db.backup(name))
	}
	return execCommand, nil
}

func (db *Mongo) backup(dbname string) []string {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "--host", db.host)
	dumpArgs = append(dumpArgs, "--port", db.port)
	if len(db.username) > 0 {
		dumpArgs = append(dumpArgs, "--username", db.username)
	}
	if len(db.password) > 0 {
		dumpArgs = append(dumpArgs, "--password", db.password)
	}
	dumpArgs = append(dumpArgs, "--db", dbname)
	dumpArgs = append(dumpArgs, "--gzip")
	sqlfile := fmt.Sprintf("%s-%s-%s-%s", "mongo", db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "--out", sqlfile)

	return dumpArgs
}
