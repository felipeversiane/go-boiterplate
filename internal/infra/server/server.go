package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/felipeversiane/go-boiterplate/internal/infra/config"
	"github.com/felipeversiane/go-boiterplate/internal/infra/database"
)

type ServerInterface interface {
	SetupRouter()
	Start()
}

type server struct {
	mux      *http.ServeMux
	config   config.ServerConfig
	database database.DatabaseInterface
}

func NewServer(
	cfg config.ServerConfig,
	database database.DatabaseInterface,
) ServerInterface {
	return &server{
		mux:      http.NewServeMux(),
		config:   cfg,
		database: database,
	}
}

func (server *server) SetupRouter() {
	server.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "healthy"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}

func (server *server) Start() {
	port := server.config.Port
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server.mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
