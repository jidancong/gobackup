package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// config
type (
	GobackupConfig struct {
		Level          string            `yaml:"level" validate:"required"`
		Store          string            `yaml:"store"`
		Cron           string            `yaml:"cron" validate:"required,cron"`
		ConfigDatabase []*ConfigDatabase `yaml:"backup" validate:"required,dive,required"`
	}

	ConfigDatabase struct {
		Port     string `yaml:"port" validate:"required,number"`
		Host     string `yaml:"host" validate:"required,ipv4"`
		DBType   string `yaml:"type" validate:"required"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
)

func NewConfig(path string) (GobackupConfig, error) {
	var cfg GobackupConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return GobackupConfig{}, fmt.Errorf("读取配置失败 Error: %w", err)
	}
	if err := validateStruct(cfg); err != nil {
		return GobackupConfig{}, fmt.Errorf("配置文件错误 Error: %w", err)
	}

	return cfg, nil
}

func validateStruct(config GobackupConfig) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&config); err != nil {
		return err
	}
	return nil
}
