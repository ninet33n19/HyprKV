package resp_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ninet33n19/XiaoKV/internal/resp"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		value, err := resp.Decode([]byte(k))
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
		value, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if value != v {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
		}
	}
}

func TestInteger(t *testing.T) {
	cases := map[string]int64{
		":-123\r\n": -123,
		":0\r\n":    0,
		":1000\r\n": 1000,
		":+12\r\n":  12,
	}

	for k, v := range cases {
		value, err := resp.Decode([]byte(k))
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
		"$0\r\n\r\n":      []byte(""),
		"$-1\r\n":         nil,
	}

	for k, v := range cases {
		value, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if !bytes.Equal(value.([]byte), v) {
			t.Errorf("Decode(%s) = %v, want %v", k, value, v)
		}
	}
}

func TestArray(t *testing.T) {
	cases := map[string][]any{
		"*3\r\n:1\r\n:2\r\n:3\r\n":             {int64(1), int64(2), int64(3)},
		"*0\r\n":                               {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n": {[]byte("hello"), []byte("world")},
	}

	for k, v := range cases {
		value, err := resp.Decode([]byte(k))
		if err != nil {
			t.Errorf("Decode(%s) failed: %v", k, err)
		}
		if !reflect.DeepEqual(value.([]any), v) {
			t.Errorf(`Got = %v, want %v`, value, v)
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
