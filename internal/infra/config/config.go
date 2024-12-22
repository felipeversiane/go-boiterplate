package config

import (
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
)

type config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Log      LogConfig
}

type ConfigInterface interface {
	GetDatabaseConfig() DatabaseConfig
	GetServerConfig() ServerConfig
	GetLogConfig() LogConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

type ServerConfig struct {
	Port string
}

type LogConfig struct {
	Level string
}

func NewConfig() ConfigInterface {
	var cfg *config
	once.Do(func() {
		cfg = &config{
			Database: DatabaseConfig{
				Host:     getEnv("POSTGRES_HOST"),
				Port:     getEnv("POSTGRES_PORT"),
				User:     getEnv("POSTGRES_USER"),
				Password: getEnv("POSTGRES_PASSWORD"),
				Name:     getEnv("POSTGRES_DB"),
				SslMode:  getEnv("POSTGRES_SSL"),
			},
			Server: ServerConfig{
				Port: getEnv("SERVER_PORT"),
			},
			Log: LogConfig{
				Level: getEnv("LOG_LEVEL"),
			},
		}
	})
	return cfg
}

func (config *config) GetDatabaseConfig() DatabaseConfig {
	return config.Database
}

func (config *config) GetServerConfig() ServerConfig {
	return config.Server
}

func (c *config) GetLogConfig() LogConfig {
	return c.Log
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return value
}
