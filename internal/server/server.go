package server

import (
	"log"
	"net/http"
	"time"

	"github.com/ninet33n19/HyprKV/internal/store"
)

type kvStore interface {
	Get(key string) (string, bool)
	Set(key, value string) error
}

type Server struct {
	store kvStore
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{key}", s.Get)
	mux.HandleFunc("POST /set/{key}", s.Set)
	mux.HandleFunc("GET /healthz", s.HealthCheck)

	return mux
}

func (s *Server) ListenAndServe() {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      s.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server listening on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
