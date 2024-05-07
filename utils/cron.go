package utils

import (
	"github.com/robfig/cron/v3"
)

type CronTask struct {
	co *cron.Cron
}

func NewCronTask() *CronTask {
	co := cron.New()
	return &CronTask{co: co}
}

func (c *CronTask) Add(spec string, cmd func()) {
	c.co.AddFunc(spec, cmd)
}

func (c *CronTask) Start() {
	c.co.Start()
}
