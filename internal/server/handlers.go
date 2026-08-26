package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "key is required")
		return
	}

	value, ok := s.store.Get(key)
	if !ok {
		writeError(w, http.StatusNotFound, "key not found")
		return
	}

	writeJSON(w, http.StatusOK, value)
}

func (s *Server) Set(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "key is required")
		return
	}

	limitedBody := http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read body")
		return
	}

	value := string(body)
	if !json.Valid(body) {
		writeError(w, http.StatusBadRequest, "value must be valid JSON")
		return
	}

	if err := s.store.Set(key, value); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to set value")
		return
	}

	writeJSON(w, http.StatusOK, value)
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
