package resp_test

import (
	"bytes"
	"testing"

	"github.com/ninet33n19/XiaoKV/internal/resp"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		value, _, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if value != v {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
			t.Fail()
		}
	}
}

func TestError(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}

	for k, v := range cases {
		value, _, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if value != v {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
		}
	}
}

func TestInteger(t *testing.T) {
	cases := map[string]int{
		":-123\r\n": -123,
		":0\r\n":    0,
		":1000\r\n": 1000,
	}

	for k, v := range cases {
		value, _, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if value != v {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
		}
	}
}

func TestBulk(t *testing.T) {
	cases := map[string][]byte{
		"$5\r\nhello\r\n": []byte("hello"),
		// "$0\r\n\r\n":      "",
	}

	for k, v := range cases {
		value, _, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if !bytes.Equal(value.([]byte), v) {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
		}
	}
}

func BenchmarkIntegerParse(b *testing.B) {
	for b.Loop() {
		resp.Decode([]byte(":1000\r\n"))
	}
}

func BenchmarkBulkParse(b *testing.B) {
	for b.Loop() {
		resp.Decode([]byte("$5\r\nhello\r\n"))
	}
}
