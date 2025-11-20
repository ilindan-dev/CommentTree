package config

import (
	"errors"
	"time"

	"github.com/wb-go/wbf/config"
)

var ErrConfigNotFound = errors.New("config not found")

// Config holds the application configuration.
type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

// AppConfig holds the application-specific configuration.
type AppConfig struct {
	Port string `yaml:"port"`
}

// DBConfig holds the database configuration.
type DBConfig struct {
	DSN  string     `yaml:"dsn"`
	Pool PoolConfig `yaml:"pool"`
}

// PoolConfig holds the database connection pool configuration.
type PoolConfig struct {
	MaxOpenConns int32         `yaml:"max_open_conns"`
	MaxIdleConns int32         `yaml:"max_idle_conns"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
}

func LoadConfig() (*Config, error) {
	cfgLoader := config.New()
	_ = cfgLoader.LoadConfigFiles("configs/config.yaml")

	cfgLoader.EnableEnv("APP")

	var cfg Config
	if err := cfgLoader.Unmarshal(&cfg); err != nil {
		return nil, ErrConfigNotFound
	}
	return &cfg, nil
}
