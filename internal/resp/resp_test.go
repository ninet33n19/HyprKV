package resp_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/ninet33n19/XiaoKV/internal/resp"
)

func TestRESPDecode(t *testing.T) {
	t.Run("simple string", func(t *testing.T) {
		cases := map[string]string{
			"+OK\r\n": "OK",
		}

		for k, v := range cases {
			value, err := resp.Decode([]byte(k))
			if err != nil {
				t.Errorf("Decode(%s) failed: %v", k, err)
			}
			if value != v {
				t.Errorf("Got = %v, want %v", value, v)
				t.Fail()
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		cases := map[string]string{
			"-Error message\r\n": "Error message",
		}

		for k, v := range cases {
			value, err := resp.Decode([]byte(k))
			if err != nil {
				t.Errorf("Decode(%s) failed: %v", k, err)
			}
			if value != v {
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("integer", func(t *testing.T) {
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
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("bulk string", func(t *testing.T) {
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
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("array", func(t *testing.T) {
		cases := map[string][]any{
			"*3\r\n:1\r\n:2\r\n:3\r\n":              {int64(1), int64(2), int64(3)},
			"*0\r\n":                                {},
			"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":  {[]byte("hello"), []byte("world")},
			"*2\r\n$7\r\nCOMMAND\r\n$4\r\nDOCS\r\n": {[]byte("COMMAND"), []byte("DOCS")},
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
	})
}

func TestRespEncode(t *testing.T) {
	t.Run("null values", func(t *testing.T) {
		cases := map[any][]byte{
			nil: []byte("$-1\r\n"),
		}

		for k, v := range cases {
			value, err := resp.Encode(k)
			if err != nil {
				t.Errorf("Encode(%s) failed: %v", k, err)
			}
			if !bytes.Equal(value, v) {
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("simple string", func(t *testing.T) {
		cases := map[string][]byte{
			"OK": []byte("+OK\r\n"),
		}

		for k, v := range cases {
			value, err := resp.Encode(k)
			if err != nil {
				t.Errorf("Encode(%s) failed: %v", k, err)
			}
			if !bytes.Equal(value, v) {
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("bulk string", func(t *testing.T) {
		cases := []struct {
			input []byte
			want  []byte
		}{
			{
				input: []byte("hello"),
				want:  []byte("$5\r\nhello\r\n"),
			},
			{
				input: []byte(""),
				want:  []byte("$0\r\n\r\n"),
			},
		}

		for _, tc := range cases {
			value, err := resp.Encode(tc.input)
			if err != nil {
				t.Errorf("Encode(%s) failed: %v", tc.input, err)
			}
			if !bytes.Equal(value, tc.want) {
				t.Errorf("Got = %v, want %v", value, tc.want)
			}
		}
	})

	t.Run("integer", func(t *testing.T) {
		cases := map[int64][]byte{
			-123: []byte(":-123\r\n"),
			0:    []byte(":0\r\n"),
			1000: []byte(":1000\r\n"),
			12:   []byte(":12\r\n"),
		}

		for k, v := range cases {
			value, err := resp.Encode(k)
			if err != nil {
				t.Errorf("Encode(%d) failed: %v", k, err)
			}
			if !bytes.Equal(value, v) {
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})

	t.Run("array", func(t *testing.T) {
		cases := []struct {
			input []any
			want  []byte
		}{
			{
				input: []any{int64(1), int64(2), int64(3)},
				want:  []byte("*3\r\n:1\r\n:2\r\n:3\r\n"),
			},
			{
				input: []any{},
				want:  []byte("*0\r\n"),
			},
			{
				input: []any{[]byte("hello"), []byte("world")},
				want:  []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			},
		}

		for _, tc := range cases {
			value, err := resp.Encode(tc.input)
			if err != nil {
				t.Errorf("Encode(%v) failed: %v", tc.input, err)
			}
			if !bytes.Equal(value, tc.want) {
				t.Errorf("Got = %v, want %v", value, tc.want)
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		cases := map[error][]byte{
			errors.New("error message"): []byte("-error message\r\n"),
		}

		for k, v := range cases {
			value, err := resp.Encode(k)
			if err != nil {
				t.Errorf("Encode(%v) failed: %v", k, err)
			}
			if !bytes.Equal(value, v) {
				t.Errorf("Got = %v, want %v", value, v)
			}
		}
	})
}

// ---- BENCHMARKS ----
func BenchmarkBulkParse(b *testing.B) {
	for b.Loop() {
		resp.Decode([]byte("$5\r\nhello\r\n"))
	}
}

func BenchmarkBulkEncode(b *testing.B) {
	for b.Loop() {
		resp.Encode([]byte("hello"))
	}
}

func BenchmarkArrayEncode(b *testing.B) {
	for b.Loop() {
		resp.Encode([]any{int64(1), int64(2), int64(3)})
	}
}

func BenchmarkArrayParse(b *testing.B) {
	for b.Loop() {
		resp.Decode([]byte("*3\r\n:1\r\n:2\r\n:3\r\n"))
	}
}
