package tcp

import (
	"bufio"
	"fmt"
	"net"
	"tcpstore/logg"
	"tcpstore/service"
)

var dh service.DataHandler

func EnableTcpServer(datasource service.Datasource) {
	logg.Info.Println("Starting TCP Server")
	listener, err := net.Listen("tcp4", ":8000")
	logg.Info.Println("TCP Server accessible at", listener.Addr().String())

	if err != nil {
		logg.Error.Fatalf("Connection failure: %s", err.Error())
	}

	defer func() { _ = listener.Close() }()

	logg.Info.Println("Setting up Data Handler with datasource:", datasource.Name())
	dh = service.NewDataHandler(datasource)
	go dh.RequestMonitor()

	for {
		logg.Info.Println("Waiting for Client connection")
		connection, err := listener.Accept()
		if err != nil {
			logg.Error.Printf("connection error: %s", err.Error())
			break
		}
		logg.Info.Printf("Established Connection to %s", connection.RemoteAddr().String())
		go handle(connection)
	}
}

func handle(connection net.Conn) {
	logg.Info.Println("Handling Requests from Client", connection.RemoteAddr().String())
	defer func() { _ = connection.Close() }()
	defer logg.Info.Println("Closed Connection")

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		if connection == nil {
			break
		}
		command, _ := service.ParseMessage(scanner.Text())
		if command.Valid() {
			request := command.AsRequest()
			dh.QueueRequest(request)
			response := <-request.ResponseChannel
			logg.Response.Printf("Responded [%s] to Request [%s]", service.TrimMessage(response.ClientString()), service.TrimMessage(scanner.Text()))
			_, err := connection.Write([]byte(fmt.Sprint(response.ClientString())))
			if err != nil {
				logg.Error.Println("Unable to send response to Client", connection.RemoteAddr().String())
			}
		} else {
			logg.Info.Printf("Received invalid command: %s", scanner.Text())
			logg.Response.Printf("Responded [%s] to Request [%s]", "err", service.TrimMessage(scanner.Text()))
			_, err := connection.Write([]byte("err"))
			if err != nil {
				logg.Error.Println("Unable to send response to Client", connection.RemoteAddr().String())
			}
		}
	}
}
