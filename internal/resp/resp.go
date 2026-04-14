package resp

import (
	"errors"
	"fmt"
)

func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeOne(data)

	return value, err
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '$':
		return readBulk(data)
	case '*':
		return readArray(data)
	// case '_':
	// 	return readNull(data)
	default:
		return nil, 0, errors.New("unknown type")
	}
}

func readSimpleString(data []byte) (string, int, error) {
	pos := 1

	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1

	var number int64
	var isNegative bool = false

	for ; data[pos] != '\r'; pos++ {
		if data[pos] == '-' {
			isNegative = true
			continue
		} else if data[pos] == '+' {
			continue
		}
		number = number*10 + int64(data[pos]-'0')
	}

	if isNegative {
		number = -number
	}

	return number, pos + 2, nil
}

func readBulk(data []byte) ([]byte, int, error) {
	length, nextPos, err := readInteger(data)
	if err != nil {
		return nil, 0, err
	}

	if length == -1 {
		return nil, nextPos, nil
	}

	start := nextPos
	end := start + int(length)

	if end > len(data) {
		return nil, 0, errors.New("bulk length exceeds data")
	}

	if data[end] != '\r' || data[end+1] != '\n' {
		return nil, 0, errors.New("invalid bulk string termination")
	}

	return data[start:end], end + 2, nil
}

func readArray(data []byte) ([]any, int, error) {
	length, nextPos, err := readInteger(data)
	if err != nil {
		return nil, 0, err
	}
	items := make([]any, length)
	for i := range items {
		item, delta, err := DecodeOne(data[nextPos:])
		if err != nil {
			return nil, 0, err
		}
		items[i] = item
		nextPos += delta
	}

	return items, nextPos, nil
}

func Encode(val any) ([]byte, error) {
	switch v := val.(type) {
	case nil:
		return encodeNull(), nil
	case string:
		return encodeSimpleString(v), nil
	case []byte:
		return encodeBulkString(v), nil
	case int:
		return encodeInteger(int64(v)), nil
	case int64:
		return encodeInteger(v), nil
	case error:
		return encodeError(v.Error()), nil
	case []any:
		return encodeArray(v)
	default:
		return nil, errors.New("unsupported type")
	}
}

func encodeNull() []byte {
	return []byte("$-1\r\n")
}

func encodeSimpleString(data string) []byte {
	return fmt.Appendf(nil, "+%s\r\n", data)
}

func encodeBulkString(data []byte) []byte {
	length := len(data)
	result := append(fmt.Appendf(nil, "$%d\r\n", length), data...)
	result = append(result, "\r\n"...)

	return result
}

func encodeInteger(n int64) []byte {
	return fmt.Appendf(nil, ":%d\r\n", n)
}

func encodeError(message string) []byte {
	return fmt.Appendf(nil, "-%s\r\n", message)
}

func encodeArray(items []any) ([]byte, error) {
	if items == nil {
		return []byte("*-1\r\n"), nil
	}

	result := fmt.Appendf(nil, "*%d\r\n", len(items))

	for _, item := range items {
		encoded, err := Encode(item)
		if err != nil {
			return nil, err
		}
		result = append(result, encoded...)
	}

	return result, nil
}
