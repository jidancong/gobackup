package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// config
type (
	GobackupConfig struct {
		Level          string           `yaml:"level" validate:"required"`
		Store          string           `yaml:"store"`
		Cron           string           `yaml:"cron" validate:"required,cron"`
		ConfigDatabase []ConfigDatabase `yaml:"backup"`
	}

	ConfigDatabase struct {
		Port     string `yaml:"port" validate:"required,number"`
		Host     string `yaml:"host" validate:"required,ipv4"`
		DBType   string `yaml:"type" validate:"required"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		FromPath string `yaml:"fromPath"`
	}
)

func NewConfig(path string) (GobackupConfig, error) {
	var cfg GobackupConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return GobackupConfig{}, err
	}

	return cfg, nil
}

// 校验配置
func (config *GobackupConfig) Validate() error {
	if len(config.ConfigDatabase) <= 0 {
		return fmt.Errorf("没有配置需备份数据库信息")
	}

	if config.Cron == "" {
		return fmt.Errorf("定时任务为空")
	}

	if config.Store == "" {
		return fmt.Errorf("存储目录为空")
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return err
	}
	return nil
}
