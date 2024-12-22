package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

func (s *server) SetupRouter() {
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "healthy"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}

func (s *server) Start() {
	port := s.config.Port

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), s.mux)
	if err != nil {
		slog.Error("Failed to start server", "port", port, "error", err)
	}
}
