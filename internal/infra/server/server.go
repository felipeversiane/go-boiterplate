package server

import (
	"fmt"
	"log/slog"

	"github.com/felipeversiane/go-boiterplate/internal/domain/user"
	"github.com/felipeversiane/go-boiterplate/internal/infra/config"
	"github.com/felipeversiane/go-boiterplate/internal/infra/database"
	"github.com/gin-gonic/gin"
)

type ServerInterface interface {
	SetupRouter()
	Start()
}

type server struct {
	router   *gin.Engine
	config   config.ServerConfig
	database database.DatabaseInterface
}

func NewServer(
	cfg config.ServerConfig,
	database database.DatabaseInterface,
) ServerInterface {
	server := &server{
		router:   gin.New(),
		config:   cfg,
		database: database,
	}
	return server
}

func (server *server) SetupRouter() {
	server.router.Use(gin.Recovery())
	v1 := server.router.Group("/api/v1")
	{
		user.UserRouter(v1, server.database)
	}
	server.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
}

func (server *server) Start() {
	port := server.config.Port
	err := server.router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
	}
}
