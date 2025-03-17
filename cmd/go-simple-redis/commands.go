package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fabinzne/go-simple-redis/internal/storage"
)

func processCommand(parts []string, store *storage.DataStore) string {
	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) < 3 {
			return "ERROR: Syntax: SET key value"
		}
		store.Set(parts[1], parts[2])
		return "OK"

	case "GET":
		if len(parts) < 2 {
			return "ERROR: Syntax: GET key"
		}
		val, ok := store.Get(parts[1])
		if !ok {
			return "NULL"
		}
		return fmt.Sprintf("%v", val)

	case "INCR":
		if len(parts) < 2 {
			return "ERROR: Syntax: INCR key"
		}
		val, ok := store.Get(parts[1])
		if !ok {
			store.Set(parts[1], 1)
			return "1"
		}
		if num, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
			store.Set(parts[1], num+1)
			return fmt.Sprintf("%d", num+1)
		}
		return "ERROR: Value is not an integer"

	case "DEL":
		if len(parts) < 2 {
			return "ERROR: Syntax: DEL key"
		}
		store.Delete(parts[1])
		return "OK"

	case "MGET":
		if len(parts) < 2 {
			return "ERROR: Syntax: MGET key [key ...]"
		}
		var values []string
		for _, key := range parts[1:] {
			val, ok := store.Get(key)
			if !ok {
				values = append(values, "NULL")
			} else {
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
		return strings.Join(values, " ")

	case "MSET":
		if len(parts) < 3 || len(parts)%2 != 1 {
			return "ERROR: Syntax: MSET key value [key value ...]"
		}

		for i := 1; i < len(parts); i += 2 {
			store.Set(parts[i], parts[i+1])
		}

		return "OK"

	case "FLUSH":
		store.Flush()
		return "OK"

	case "EXPIRE":
		if len(parts) < 3 {
			return "ERROR: Syntax: EXPIRE key seconds"
		}

		sec, err := strconv.Atoi(parts[2])
		if err != nil {
			return "ERROR: Invalid TTL"
		}

		if store.Expire(parts[1], time.Duration(sec)*time.Second) {
			return "1"
		}

		return "0"

	case "TTL":
		if len(parts) < 2 {
			return "ERRO: Sintaxe: TTL chave"
		}

		_, exists, exp := store.SafeGetKeyInfo(parts[1])
		if !exists {
			return "-2"
		}

		if exp.IsZero() {
			return "-1"
		}

		remaining := time.Until(exp).Seconds()
		if remaining < 0 {
			return "-2"
		}
		return fmt.Sprintf("%d", int(remaining))

	case "SAVE":
		if err := store.SaveToFile(*dumpFile); err != nil {
			return fmt.Sprintf("ERROR: %v", err)
		}
		return "OK"

	default:
		return "ERROr: Unknown command"
	}
}
