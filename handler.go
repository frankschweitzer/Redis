package main

import (
	"sync"
)

var MAP = map[string]string{}
var MAPmu = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING" : ping,
	"SET" : set,
	"GET" : get,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ:"error", str:"ERROR: WRONG NUMBER OF ARGUMENTS"}
	}

	key := args[0].bulk
	value := args[1].bulk

	MAPmu.Lock()
	MAP[key] = value
	MAPmu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ:"error", str:"ERROR: WRONG NUMBER OF ARGUMENTS"}
	}

	key := args[0].bulk

	MAPmu.Lock()
	value, ok := MAP[key]
	MAPmu.Unlock()
	if !ok {
		return Value{typ:"error", str:"ERROR: COULD NOT GET VALUE"}
	}

	return Value{typ:"string", str:value}
}