package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/ninet33n19/XiaoKV/internal/config"
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

		log.Println(string(buf[:n]))
	}
}
