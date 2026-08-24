package server

import (
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read body")
		return
	}

	value := string(body)
	if err := s.store.Set(key, value); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to set value")
		return
	}

	writeJSON(w, http.StatusOK, value)
}
