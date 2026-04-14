package main

import (
	"net"
	"testing"
	"time"

	"github.com/ninet33n19/XiaoKV/internal/resp"
)

func TestDispatchCommand(t *testing.T) {
	t.Run("ping without args", func(t *testing.T) {
		got, err := dispatchCommand([]any{[]byte("PING")})
		if err != nil {
			t.Fatalf("dispatchCommand returned error: %v", err)
		}
		if got != "PONG" {
			t.Fatalf("got %v, want PONG", got)
		}
	})

	t.Run("command docs returns array", func(t *testing.T) {
		got, err := dispatchCommand([]any{[]byte("COMMAND"), []byte("DOCS")})
		if err != nil {
			t.Fatalf("dispatchCommand returned error: %v", err)
		}
		reply, ok := got.([]any)
		if !ok {
			t.Fatalf("got type %T, want []any", got)
		}
		if len(reply) != 0 {
			t.Fatalf("got %v, want empty array", reply)
		}
	})

	t.Run("invalid payload type", func(t *testing.T) {
		_, err := dispatchCommand([]byte("PING"))
		if err == nil {
			t.Fatalf("expected error for non-array payload")
		}
	})
}

func TestWriteResp(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	done := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 64)
		_ = client.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := client.Read(buf)
		if err != nil {
			done <- nil
			return
		}
		done <- buf[:n]
	}()

	if err := writeResp(server, "PONG"); err != nil {
		t.Fatalf("writeResp returned error: %v", err)
	}

	select {
	case got := <-done:
		want, err := resp.Encode("PONG")
		if err != nil {
			t.Fatalf("resp encode failed: %v", err)
		}
		if string(got) != string(want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out waiting for response write")
	}
}
