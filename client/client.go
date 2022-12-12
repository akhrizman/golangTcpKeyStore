package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Starting Client")
	message := "get12hello34world"
	bytes := []byte(fmt.Sprint(message, "\n"))
	fmt.Println("Connecting to TCP Key Store Server")
	connection, err := net.Dial("tcp4", ":8000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Sent Message: ", message)

	if _, err = connection.Write(bytes); err != nil {
		log.Fatalf("write failure: %s", err.Error())
	}

	fmt.Println("Waiting for Response")
	if _, err = connection.Read(bytes); err != nil {
		log.Fatalf("read failure: %s", err.Error())
	}

	fmt.Println("Received Response: ", string(bytes))
	_ = connection.Close()

}
