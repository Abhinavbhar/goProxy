package main

import (
	"fmt"
	"net"
	helperfunc "proxy/HelperFunc"
	"time"
)

func main() {
	helperfunc.InitMongo()
	go helperfunc.StartBandwidthUpdater()
	start := time.Now()
	var Ips []string = helperfunc.LoadAllowedIp()

	for i := 0; i < len(Ips); i++ {
		helperfunc.IpBandwidth[Ips[i]] = int64(0)
	}

	elapsed := time.Since(start)
	fmt.Println("time take:", elapsed, "to load", len(Ips), "Ips in memory")
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
		go helperfunc.TcpHandler(conn)
	}
}
