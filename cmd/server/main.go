package main

import (
	"github.com/ninet33n19/HyprKV/internal/server"
	"github.com/ninet33n19/HyprKV/internal/store"
)

func main() {
	store := store.NewStore()
	server := server.NewServer(store)

	server.ListenAndServe()
}
