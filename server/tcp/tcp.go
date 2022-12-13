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

	service.ResponseStreams = map[int]chan service.Response{}

	connectionId := 1
	for {
		fmt.Println("Waiting for Client connection")
		connection, err := listener.Accept()
		service.ResponseStreams[connectionId] = make(chan service.Response)
		if err != nil {
			fmt.Printf("connection error: %s", err.Error())
			break
		}
		fmt.Println("Established Connection")
		fmt.Println("===============Handling Client Request on Connection ", connectionId)
		go handle(connection, connectionId)
		connectionId++
	}
}

func handle(connection net.Conn, connectionId int) {
	defer func() { _ = connection.Close() }()
	defer fmt.Println("Closed Connection ", connectionId)

	go func() {
		for response := range service.ResponseStreams[connectionId] {
			fmt.Printf("Sending Message to Client over connection %d: [%s]\n", connectionId, fmt.Sprint(response.ClientString()))
			connection.Write([]byte(fmt.Sprint(response.ClientString())))
		}
	}()

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		fmt.Printf("\nReceived message from Client over connection %d: [%s]\n", connectionId, scanner.Text())
		command, _ := service.ParseCommand(scanner.Text())
		if command.Valid() {
			fmt.Println("Command is Valid")
			service.Store.QueueRequest(command.ToRequest())
		} else {
			service.ResponseStreams[connectionId] <- service.NewResponse("err", "", "", connectionId)
		}
	}
}
