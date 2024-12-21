package config

import (
	"fmt"
	"os"
	"sync"
)

var (
	once sync.Once
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type ConfigInterface interface {
	GetDatabaseConfig() DatabaseConfig
	GetServerConfig() ServerConfig
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

func NewConfig() ConfigInterface {
	var cfg *Config
	once.Do(func() {
		cfg = &Config{
			Database: DatabaseConfig{
				Host:     getEnvOrDie("POSTGRES_HOST"),
				Port:     getEnvOrDie("POSTGRES_PORT"),
				User:     getEnvOrDie("POSTGRES_USER"),
				Password: getEnvOrDie("POSTGRES_PASSWORD"),
				Name:     getEnvOrDie("POSTGRES_DB"),
				SslMode:  getEnvOrDie("POSTGRES_SSL"),
			},
			Server: ServerConfig{
				Port: getEnvOrDie("SERVER_PORT"),
			},
		}
	})
	return cfg
}

func (config *Config) GetDatabaseConfig() DatabaseConfig {
	return config.Database
}

func (config *Config) GetServerConfig() ServerConfig {
	return config.Server
}

func getEnvOrDie(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("missing environment variable %s", key))
	}
	return value
}
