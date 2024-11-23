package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var configPath = "."

func SetConfigPath(p string) {
	configPath = p
}

type Config struct {
	Postgres
	Redis RedisConfig
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewConfig() (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	c := Config{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
