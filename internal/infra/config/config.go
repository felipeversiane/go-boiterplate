package config

import (
	"log"
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

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return value
}