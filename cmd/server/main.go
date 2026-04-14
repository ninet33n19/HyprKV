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

		go handleClient(conn, &concurrent_clients)
	}
}

func handleClient(conn net.Conn, concurrent_clients *int64) {
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

		payload := buf[:n]
		log.Print(string(payload))

		val, err := resp.Decode(payload)
		if err != nil {
			log.Println(err)
			writeResp(conn, errors.New("ERR parse error"))
			continue
		}
		log.Printf("Received command: %v", val)

		reply, err := dispatchCommand(val)
		if err != nil {
			if writeErr := writeResp(conn, err); writeErr != nil {
				log.Println(writeErr)
				return
			}
			continue
		}

		if err := writeResp(conn, reply); err != nil {
			log.Println(err)
			return
		}
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
