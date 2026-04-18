package server

import (
	"errors"
	"fmt"
	"strings"
)

func (s *Server) routeCommand(value any) any {
	args, ok := value.([]any)
	if !ok || len(args) == 0 {
		return errors.New("ERR unknown command")
	}

	command := strings.ToUpper(string(args[0].([]byte)))

	switch command {
	case "PING":
		if len(args) > 2 {
			return errors.New("ERR wrong number of arguments for 'ping' command")
		}
		if len(args) > 1 {
			msg, ok := args[1].([]byte)
			if !ok {
				return errors.New("ERR ping argument must be bulk string")
			}
			return msg
		}
		return "PONG"
	case "SET":
		if len(args) < 3 {
			return errors.New("ERR set requires key and value")
		}

		key := string(args[1].([]byte))
		val := args[2].([]byte)

		s.storage.Set(key, val, 0)
		return "OK"
	case "GET":
		if len(args) != 2 {
			return errors.New("ERR wrong number of arguments for 'get' command")
		}
		key := string(args[1].([]byte))
		val, exists := s.storage.Get(key)
		if !exists {
			return nil // Encode will turn this into the Null Bulk String ($-1\r\n)
		}
		return val
	case "DEL":
		if len(args) < 2 {
			return errors.New("ERR wrong number of arguments for 'DEL' command")
		}

		keys := make([]string, 0, len(args)-1)
		for _, arg := range args[1:] {
			key, ok := arg.([]byte)
			if !ok {
				return errors.New("ERR DEL arguments must be bulk strings")
			}
			keys = append(keys, string(key))
		}

		deleted, err := s.storage.Delete(keys...)
		if err != nil {
			return err
		}
		return deleted
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
