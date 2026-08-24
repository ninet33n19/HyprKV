package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ninet33n19/HyprKV/internal/store"
)

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	value, ok := s.store.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	payload := json.RawMessage(value)
	if !json.Valid(payload) {
		http.Error(w, "stored value is not valid JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (s *Server) Set(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	value := string(body)
	err = s.store.Set(key, value)
	if err != nil {
		http.Error(w, "failed to set value", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (s *Server) ListenAndServe() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{key}", s.Get)
	mux.HandleFunc("POST /set/{key}", s.Set)
	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
