package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ninet33n19/XiaoKV/internal/config"
	"github.com/ninet33n19/XiaoKV/internal/resp"
)

func main() {
	cfg := config.NewConfig("127.0.0.1", 7379)

	log.Println("Starting synchronous TCP server on ", cfg.Addr, cfg.Port)

	var concurrent_clients int64

	listener, err := net.Listen("tcp", cfg.Addr+":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		atomic.AddInt64(&concurrent_clients, 1)
		log.Println("Concurrent clients:", concurrent_clients)

		go handleConnection(conn, &concurrent_clients)
	}
}

func handleConnection(conn net.Conn, concurrent_clients *int64) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				atomic.AddInt64(concurrent_clients, -1)
				log.Println("Client disconnected")
				return
			}
			log.Println(err)
			return
		}

		value, err := resp.Decode(buf[:n])
		if err != nil {
			fmt.Println("Decode error:", err)
			continue
		}

		response := handleCommand(value)

		encoded, _ := resp.Encode(response)
		conn.Write(encoded)
	}
}

func handleCommand(value any) any {
	args, ok := value.([]any)
	if !ok || len(args) == 0 {
		return errors.New("ERR unknown command")
	}

	command := strings.ToUpper(string(args[0].([]byte)))

	switch command {
	case "PING":
		if len(args) > 1 {
			msg, ok := args[1].([]byte)
			if !ok {
				return errors.New("ERR ping argument must be bulk string")
			}
			return msg
		}
		return "PONG"
	case "ECHO":
		if len(args) > 1 {
			msg, ok := args[1].([]byte)
			if !ok {
				return errors.New("ERR echo argument must be bulk string")
			}
			return msg
		}
		return errors.New("ERR echo requires an argument")
	case "COMMAND":
		return []any{}
	default:
		return fmt.Errorf("ERR unknown command '%s'", command)
	}
}

func dispatchCommand(val any) (any, error) {
	parts, ok := val.([]any)
	if !ok {
		return nil, errors.New("ERR expected array command")
	}
	if len(parts) == 0 {
		return nil, errors.New("ERR empty command")
	}

	cmdRaw, ok := parts[0].([]byte)
	if !ok {
		return nil, errors.New("ERR command name must be bulk string")
	}
	cmd := strings.ToUpper(string(cmdRaw))

	switch cmd {
	case "PING":
		if len(parts) > 1 {
			msg, ok := parts[1].([]byte)
			if !ok {
				return nil, errors.New("ERR ping argument must be bulk string")
			}
			return msg, nil
		}
		return "PONG", nil
	case "COMMAND":
		return []any{}, nil
	default:
		return nil, fmt.Errorf("ERR unknown command '%s'", cmd)
	}
}

func writeResp(conn net.Conn, val any) error {
	encoded, err := resp.Encode(val)
	if err != nil {
		return err
	}

	_, err = conn.Write(encoded)
	return err
}
