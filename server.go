package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// start server on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("Error starting server: %s",err)
		return
	}
	fmt.Println("Starting server on port 6379")

	// listent to connections on the server
	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("Error accepting connections to server: %s", err)
		return
	}

	// defer will wait until this main function is finished before closing the connection
	defer conn.Close()

	// continuously scan for input
	for {
		resp := NewResp(conn)

		// use our own reader to parse the input
		value, err := resp.Read()
		if err != nil {
			fmt.Printf("ERROR READING INPUT: %s\n", err)
			return
		}

		if value.typ != "array" {
			fmt.Println("INVALID REQUEST: EXPECTED ARRAY")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("INVALID REQUEST: EXPECTED ARRAY LENGTH > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("INVALID COMMAND: ", command)
			writer.Write(Value{typ:"string", str:""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}