package server

import (
	"io"
	"net"
	"time"

	"github.com/ninet33n19/HyprKV/internal/resp"
	"github.com/ninet33n19/HyprKV/internal/storage"
	"github.com/rs/zerolog"
)

type Server struct {
	storage *storage.Storage
	logger  zerolog.Logger
}

func New(storage *storage.Storage, logger zerolog.Logger) *Server {
	s := &Server{
		storage: storage,
		logger:  logger,
	}
	s.storage.StartCleaner(time.Minute)
	return s
}

func (s *Server) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	s.logger.Info().Str("address", address).Msg("server listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			s.logger.Error().Err(err).Msg("accept failed")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	logger := s.logger.With().Str("remote_addr", conn.RemoteAddr().String()).Logger()
	logger.Debug().Msg("connection accepted")
	buf := make([]byte, 2048) // A larger buffer for commands with payloads

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Error().Err(err).Msg("read failed")
			}
			logger.Debug().Msg("connection closed")
			return
		}

		val, err := resp.Decode(buf[:n])
		if err != nil {
			logger.Warn().Err(err).Msg("failed to decode request")
			continue
		}

		logger.Debug().Int("bytes", n).Msg("request received")
		response := s.routeCommand(val)

		encoded, _ := resp.Encode(response)
		conn.Write(encoded)
	}
}
