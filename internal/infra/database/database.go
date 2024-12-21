package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/felipeversiane/go-boiterplate/internal/infra/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once sync.Once
	db   DatabaseInterface
)

type database struct {
	db     *pgxpool.Pool
	config config.DatabaseConfig
}

type DatabaseInterface interface {
	GetDB() *pgxpool.Pool
	Ping(ctx context.Context) error
	Close()
}

func NewDatabaseConnection(ctx context.Context, config config.DatabaseConfig) DatabaseInterface {
	once.Do(func() {
		dsn := getConnectionString(config)
		poolConfig, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			log.Printf("failed to parse database config. DSN: %s, Error: %v", dsn, parseErr)
			return
		}

		pool, poolErr := pgxpool.NewWithConfig(ctx, poolConfig)
		if poolErr != nil {
			log.Printf("failed to create database pool. Error: %v", poolErr)
			return
		}

		db = &database{
			db:     pool,
			config: config,
		}

		if pingErr := db.Ping(ctx); pingErr != nil {
			log.Printf("database ping failed. Error: %v", pingErr)
		}
	})

	return db
}

func (database *database) Ping(ctx context.Context) error {
	return database.db.Ping(ctx)
}

func (database *database) Close() {
	database.db.Close()
}

func (database *database) GetDB() *pgxpool.Pool {
	return database.db
}

func getConnectionString(config config.DatabaseConfig) string {
	user := config.User
	password := config.Password
	dbname := config.Name
	dbport := config.Port
	dbhost := config.Host
	sslmode := config.SslMode

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s", user, password, dbname, dbport, dbhost, sslmode)
	return dsn
}
