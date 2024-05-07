package dbs

import (
	"fmt"
	"gobackup/utils"
	"os"
)

type Postgresql struct {
	host     string
	port     string
	username string
	password string
}

func NewPostgresql(host, port, username, password string) *Postgresql {
	os.Setenv("PGPASSWORD", password)
	return &Postgresql{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (db *Postgresql) BackupAll() ([][]string, error) {
	dbNames, err := getPGDBs(db.host, db.username, db.password, db.port)
	if err != nil {
		return nil, err
	}

	execComand := [][]string{}
	for _, name := range dbNames {
		execComand = append(execComand, db.backup(name))
	}
	return execComand, nil
}

func (db *Postgresql) backup(database string) []string {
	dumpArgs := []string{}
	dumpArgs = append(dumpArgs, "-h", db.host)
	dumpArgs = append(dumpArgs, "-p", db.port)
	dumpArgs = append(dumpArgs, "-U", db.username)
	dumpArgs = append(dumpArgs, "-d", database)

	sqlFile := fmt.Sprintf("%s-%s-%s-%s-%s.sql", "pgsql", database, db.host, utils.TimeFormat(), utils.UUID())
	dumpArgs = append(dumpArgs, "-f", sqlFile)

	return dumpArgs
}
