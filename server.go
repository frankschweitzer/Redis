package main

import (
	"fmt"
	"net"
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
			fmt.Printf("ERROR READIND INPUT: %s\n", err)
			return
		}

		fmt.Print("INPUT: %s", value)
		conn.Write([]byte("+OK\r\n"))
	}
}