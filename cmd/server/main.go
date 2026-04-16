package main

import (
	"log"

	"github.com/ninet33n19/HyprKV/internal/server"
	"github.com/ninet33n19/HyprKV/internal/storage"
)

func main() {
	db := storage.New()

	srv := server.New(db)

	if err := srv.Start(":7379"); err != nil {
		log.Fatal(err)
	}
}
