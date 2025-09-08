package main

import (
	"fmt"
	"net"
	"sync"
)

var ipBandwidth = make(map[string]int64)

// RWMutex for synchronizing access
var mu sync.RWMutex

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting proxy: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("HTTP Forward Proxy listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Println("a tcp request came")
		go TcpHandler(conn)
	}
}
