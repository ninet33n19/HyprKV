package server

import (
	"log"
	"net/http"

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

	return mux
}

func (s *Server) ListenAndServe() {
	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", s.routes()); err != nil {
		log.Fatal(err)
	}
}
