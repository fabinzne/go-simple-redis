package main

import (
	"fmt"
	"strconv"
	"strings"

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

	default:
		return "ERROr: Unknown command"
	}
}
