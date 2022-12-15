package main

import (
	"tcpstore/logg"
	"tcpstore/persistence"
	"tcpstore/tcp"
)

func main() {
	logg.SetupLogging()
	defer logg.CloseLogFiles()

	// Create a datastore
	ks := persistence.NewKeyStore()

	// Setup Server and Request Handler
	tcp.EnableTcpServer(&ks)

}
