package server

import (
	"fmt"
	"io"
	"net"

	"github.com/ninet33n19/HyprKV/internal/resp"
	"github.com/ninet33n19/HyprKV/internal/storage"
)

type Server struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Println("Server listening on", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048) // A larger buffer for commands with payloads

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			return
		}

		// Decode using your resp package
		val, err := resp.Decode(buf[:n])
		if err != nil {
			continue
		}

		// Route the command
		response := s.routeCommand(val)

		// Encode and reply
		encoded, _ := resp.Encode(response)
		conn.Write(encoded)
	}
}
