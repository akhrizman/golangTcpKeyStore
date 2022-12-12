package main

import (
	"bufio"
	"fmt"
	"net"
	"tcpstore/keystore"
	"tcpstore/parser"
)

func main() {
	// input: [put][1][3]<key>[2][12]<stored value>
	// create a monitor which creates the keystore
	// create the tcp server and listener

	fmt.Println("Starting TCP Key Store Server")
	listener, err := net.Listen("tcp4", ":8000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection Accepted")
	defer func() { _ = listener.Close() }()

	for {
		fmt.Println("Waiting for Client connection")
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("connection error: %s", err.Error())
			break
		}
		fmt.Println("Established Connection")
		fmt.Println("Handling Client Request")
		go handle(connection)

		go func() {
			for response := range keystore.Store.ResponseChannel {
				connection.Write([]byte(response.ClientString()))
			}
		}()
	}
}

func handle(connection net.Conn) {
	defer fmt.Println("Closed Connection")
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		fmt.Println("Received message: ", scanner.Text())
		command := parser.ParseCommand(scanner.Text())
		if command.Valid() {
			keystore.Store.QueueRequest(command.ToRequest())
		}
	}
}
