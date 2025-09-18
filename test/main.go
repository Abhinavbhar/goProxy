package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	server := "tcpbin.com:4242" // replace with your proxy server IP:Port
	connections := 100          // number of concurrent connections
	requestsPerConn := 50       // how many requests per connection
	payload := []byte("ping\n") // data to send

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < connections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", server)
			if err != nil {
				fmt.Println("Connection error:", err)
				return
			}
			defer conn.Close()

			for j := 0; j < requestsPerConn; j++ {
				t0 := time.Now()
				_, err := conn.Write(payload)
				if err != nil {
					fmt.Println("Write error:", err)
					return
				}

				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}
				latency := time.Since(t0)
				fmt.Printf("Latency: %v\n", latency)
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	totalRequests := connections * requestsPerConn
	throughput := float64(totalRequests) / duration.Seconds()
	fmt.Printf("\nTotal requests: %d\nTotal time: %v\nThroughput: %.2f req/s\n",
		totalRequests, duration, throughput)
}
