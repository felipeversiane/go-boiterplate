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
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Fatalf("failed to parse database config : %v", err)
			return
		}

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to create database pool : %v", err)
			return
		}

		db = &database{
			db:     pool,
			config: config,
		}

		if err := db.Ping(ctx); err != nil {
			log.Printf("database ping failed : %v", err)
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
