package main

import (
	"os"
	"time"

	"github.com/ninet33n19/HyprKV/internal/server"
	"github.com/ninet33n19/HyprKV/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	logger := log.With().Str("service", "hyprkv").Logger()

	db := storage.New(logger.With().Str("component", "storage").Logger())

	srv := server.New(db, logger.With().Str("component", "server").Logger())

	if err := srv.Start(":7379"); err != nil {
		logger.Fatal().Err(err).Msg("server stopped")
	}
}
