package helperfunc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

func TcpHandler(conn net.Conn) {
	defer conn.Close()

	// Get client IP
	ip := ReturnIp(conn)
	if ip == "" {
		return
	}

	// Check IP whitelist/blacklist
	if !CheckIp(ip) {
		return
	}

	// Parse HTTP CONNECT request
	reader := bufio.NewReader(conn)
	host := ReturnHost(reader)
	if host == "" {
		return
	}

	// Send HTTP 200 Connection Established response
	firstResponse := "HTTP/1.1 200 Connection Established\r\n\r\n"
	_, err := conn.Write([]byte(firstResponse))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		return
	}

	// Connect to destination with timeout
	dest, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		fmt.Printf("Error connecting to destination %s: %v\n", host, err)
		return
	}
	defer dest.Close()

	// Create context for coordinating both goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called when function exits

	// Channel to collect errors from both goroutines
	errChan := make(chan error, 2)

	// Track total bytes for bandwidth calculation
	var totalBytes int64

	// Goroutine 1: Copy from client to destination
	go func() {
		defer cancel() // Cancel context when this direction completes

		n, err := copyWithCancel(ctx, dest, conn)

		// Update bandwidth using your existing mutex
		IpMutex.Lock()
		IpBandwidth[ip] += n / 1000 // Convert to KB
		totalBytes += n
		IpMutex.Unlock()

		if err != nil && err != io.EOF {
			fmt.Printf("Error copying client->dest: %v\n", err)
		}

		errChan <- err
	}()

	// Goroutine 2: Copy from destination to client
	go func() {
		defer cancel() // Cancel context when this direction completes

		n, err := copyWithCancel(ctx, conn, dest)

		// Update bandwidth using your existing mutex
		IpMutex.Lock()
		IpBandwidth[ip] += n / 1000 // Convert to KB
		totalBytes += n
		IpMutex.Unlock()

		if err != nil && err != io.EOF {
			fmt.Printf("Error copying dest->client: %v\n", err)
		}

		errChan <- err
	}()

	// Wait for both goroutines to complete
	// This is better than your original <-done which only waited for one
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			// Context cancelled, both goroutines should be finishing
		case err := <-errChan:
			if err != nil && err != io.EOF && err != context.Canceled {
				fmt.Printf("Copy error: %v\n", err)
			}
		}
	}

	// Print final bandwidth info
	IpMutex.RLock()
	fmt.Printf("Connection closed for IP %s. Total bandwidth: %d KB\n", ip, IpBandwidth[ip])
	IpMutex.RUnlock()
}

// copyWithCancel is like io.Copy but can be cancelled via context
func copyWithCancel(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// Use a reasonable buffer size (32KB is good for network operations)
	buf := make([]byte, 32*1024)
	var written int64

	for {
		// Check if context was cancelled before each read
		select {
		case <-ctx.Done():
			return written, ctx.Err()
		default:
		}

		// Read with a timeout to make it cancellable
		if conn, ok := src.(net.Conn); ok {
			// Set a short read deadline so we can check context regularly
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		}

		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = fmt.Errorf("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}

		if er != nil {
			// Check if it's just a timeout (which is expected due to our deadline)
			if netErr, ok := er.(net.Error); ok && netErr.Timeout() {
				continue // Continue the loop to check context
			}
			// Real error or EOF
			return written, er
		}
	}
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
