package tcp

import (
	"bufio"
	"fmt"
	"net"
	"tcpstore/service"
)

func EnableTcpServer() {
	fmt.Println("Starting TCP Key Store Server")
	listener, err := net.Listen("tcp4", ":8000")
	fmt.Println("Accessible at ", listener.Addr())

	if err != nil {
		fmt.Printf("connection failure: %s", err.Error())
		//log.Fatalf("connection failure: %s", err.Error())
	}

	defer func() { _ = listener.Close() }()

	go service.Store.RequestMonitor()

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
	}
}

func handle(connection net.Conn) {
	defer func() { _ = connection.Close() }()
	defer fmt.Println("Closed Connection")

	go func() {
		for response := range service.Store.ResponseChannel {
			fmt.Printf("Sending Message to Client: [%s]\n", fmt.Sprint(response.ClientString()))
			connection.Write([]byte(fmt.Sprint(response.ClientString())))
		}
	}()

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		fmt.Printf("\nReceived message from Client: [%s]\n", scanner.Text())
		command, _ := service.ParseCommand(scanner.Text())
		if command.Valid() {
			fmt.Println("Command is Valid")
			service.Store.QueueRequest(command.ToRequest())
		} else {
			service.Store.ResponseChannel <- service.NewResponse("err", "", "")
		}
	}
}
