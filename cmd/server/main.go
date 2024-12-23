package main

import (
	"context"

	"github.com/felipeversiane/go-boiterplate/internal/infra/config"
	"github.com/felipeversiane/go-boiterplate/internal/infra/database"
	"github.com/felipeversiane/go-boiterplate/internal/infra/server"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := database.NewDatabaseConnection(ctx, cfg.GetDatabaseConfig())
	defer database.Close()

	server := server.NewServer(cfg.GetServerConfig(), database)
	server.SetupRouter()
	server.Start()
}
