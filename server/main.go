package main

import "tcpstore/tcp"

func main() {
	// input: [put][1][3]<key>[2][12]<stored value>
	tcp.EnableTcpServer()
}
