package controller

import (
	"gobackup/config"
	database "gobackup/databases"
	"gobackup/helper"
	"log/slog"
	"os"
	"path/filepath"
)

type Controller struct {
	cfg         config.GobackupConfig
	pwd         string
	commandPath string
	h           *helper.Helper
	a           *helper.Archive
	c           *helper.CronTask
}

func NewController(cfg config.GobackupConfig, pwd string, commandPath string, h *helper.Helper, a *helper.Archive, c *helper.CronTask) *Controller {
	return &Controller{
		cfg:         cfg,
		pwd:         pwd,
		commandPath: commandPath,
		h:           h,
		a:           a,
		c:           c,
	}
}

func (con *Controller) Start() error {
	databaseTask := make([]database.Database, 0)
	for _, v := range con.cfg.ConfigDatabase {
		task, err := database.GetDatabase(v.DBType, con.h, con.commandPath, database.BaseConfig{
			Host:     v.Host,
			Port:     v.Port,
			Username: v.User,
			Password: v.Password,
		})
		if err != nil {
			return err
		}

		databaseTask = append(databaseTask, task)
	}

	// 增加定时任务
	for i := 0; i < len(databaseTask); i++ {
		task := databaseTask[i]
		con.c.Add(con.cfg.Cron, func() {
			// 备份
			fileName, err := con.BackupController(task)
			if err != nil {
				slog.Error(err.Error())
			}

			// 压缩
			if err := con.CompressController(con.a, fileName); err != nil {
				slog.Error(err.Error())
			}

		})
	}

	con.c.Start()
	return nil
}

// 备份每一个数据库
func (con *Controller) BackupController(task database.Database) ([]string, error) {
	databaseNames, err := task.GetDatabases()
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0)
	for _, databaseName := range databaseNames {
		fileName, err := task.Backup(databaseName, nil)
		if err != nil {
			os.Remove(fileName)
			slog.Error(err.Error())
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil

}

// 压缩后删掉
func (con *Controller) CompressController(ar *helper.Archive, filesName []string) error {
	if _, err := ar.Compressor(filesName); err != nil {
		return err
	}

	for _, fileName := range filesName {
		fileAbs := filepath.Join(con.pwd, fileName)
		os.RemoveAll(fileAbs)
	}
	return nil
}
