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
		//logg.Fatalf("connection failure: %s", err.Error())
	}

	defer func() { _ = listener.Close() }()

	go service.Store.RequestMonitor()

	service.ResponseStreams = map[string]chan service.Response{}

	for {
		fmt.Println("Waiting for Client connection")
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("connection error: %s", err.Error())
			break
		}
		fmt.Println("Established Connection")
		fmt.Println("===============Handling Client Request on Connection ", connection.RemoteAddr().String())
		go handle(connection)
	}
}

func handle(connection net.Conn) {
	connectionId := connection.RemoteAddr().String()
	defer func() { _ = connection.Close() }()
	defer fmt.Println("Closed Connection ", connectionId)
	defer delete(service.ResponseStreams, connectionId)
	
	service.ResponseStreams[connectionId] = make(chan service.Response)

	go func() {
		for response := range service.ResponseStreams[connectionId] {
			fmt.Printf("Sending Message to Client over connection %s: [%s]\n", connectionId, fmt.Sprint(response.ClientString()))
			connection.Write([]byte(fmt.Sprint(response.ClientString())))
		}
	}()

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		fmt.Printf("\nReceived message from Client over connection %s: [%s]\n", connectionId, scanner.Text())
		command, _ := service.ParseCommand(scanner.Text())
		if command.Valid() {
			fmt.Println("Command is Valid")
			request := command.ToRequest()
			request.ConnectionId = connectionId
			service.Store.QueueRequest(request)
		} else {
			service.ResponseStreams[connectionId] <- service.NewResponse("err", "", "", connectionId)
		}
	}
}
