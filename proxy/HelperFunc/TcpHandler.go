package helperfunc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func TcpHandler(conn net.Conn) {
	// Do not close immediately, wait until both goroutines finish
	// defer conn.Close()

	// Get client IP
	ip := ReturnIp(conn)
	if ip == "" {
		conn.Close()
		return
	}

	// Check IP whitelist/blacklist
	if !CheckIp(ip) {
		conn.Close()
		return
	}

	// Parse HTTP CONNECT request
	reader := bufio.NewReader(conn)
	host := ReturnHost(reader)
	if host == "" {
		conn.Close()
		return
	}

	// Send HTTP 200 Connection Established response
	firstResponse := "HTTP/1.1 200 Connection Established\r\n\r\n"
	_, err := conn.Write([]byte(firstResponse))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		conn.Close()
		return
	}

	// Connect to destination with timeout
	dest, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		fmt.Printf("Error connecting to destination %s: %v\n", host, err)
		conn.Close()
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Client -> Destination
	go func() {
		defer wg.Done()
		io.Copy(dest, conn)
	}()

	// Destination -> Client
	go func() {
		defer wg.Done()
		bandwidth, _ := io.Copy(conn, dest)
		IpMutex.Lock()
		IpBandwidth[ip] += bandwidth / 1000
		IpMutex.Unlock()
	}()

	// Wait until both copies finish
	wg.Wait()

	// Now safe to close both ends
	conn.Close()
	dest.Close()

	IpMutex.RLock()
	fmt.Printf("Connection closed for IP %s. Total bandwidth: %d KB\n", ip, IpBandwidth[ip])
	IpMutex.RUnlock()
}

// Your existing ReturnIp function with minor improvements
func ReturnIp(conn net.Conn) string {
	if conn == nil {
		return ""
	}

	remoteAddr := conn.RemoteAddr()
	if remoteAddr == nil {
		return ""
	}

	host, _, err := net.SplitHostPort(remoteAddr.String())
	if err != nil {
		// If SplitHostPort fails, just return the address as-is
		return remoteAddr.String()
	}

	// Normalize IPv6 loopback to IPv4
	if host == "::1" {
		return "127.0.0.1"
	}

	return host
}
