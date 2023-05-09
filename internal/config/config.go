package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"table10/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	IsProd  *bool `yaml:"is_prod" env-default:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
	Keys    KeysConfig    `yaml:"keys"`
}

type StorageConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	UploadPath string `yaml:"upload_path"`
}

type KeysConfig struct {
	Telegram string `yaml:"telegram"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
