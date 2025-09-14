package helperfunc

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func TcpHandler(conn net.Conn) {
	defer conn.Close()
	//returns the ip of user
	Ip := ReturnIp(conn)
	//fmt.Println(Ip)
	if !CheckIp(Ip) {

		conn.Close()
		return
	}
	// refuse the connection if checkIp returns false
	reader := bufio.NewReader(conn)
	Host := ReturnHost(reader)
	//fmt.Println(Host)
	FirstResponse := "HTTP/1.1 200 Connection Established\r\n\r\n"
	conn.Write([]byte(FirstResponse))
	Dest, err := net.Dial("tcp", Host)
	if err != nil {
		fmt.Println("error connecting to the destination")
		return
	}
	defer Dest.Close()

	done := make(chan struct{})
	go func() {
		io.Copy(conn, Dest)
		done <- struct{}{}
	}()

	go func() {
		n, err := io.Copy(Dest, conn)
		if err != nil {
			fmt.Println("error copying the tcp connection")
		}
		fmt.Println("copied bytes are", n)
		IpBandwidth[Ip] = IpBandwidth[Ip] + n/1000
		fmt.Println(IpBandwidth)
		done <- struct{}{}
	}()

	<-done // wait for one direction to fin
}

// This is a really basic function it can give a common ip if some users are using ip
// from the same public wifi SO that can conflict for bandwidth
func ReturnIp(conn net.Conn) string {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return ""
	}

	// Normalize IPv6 loopback to IPv4
	if host == "::1" {
		host = "127.0.0.1"
	}

	return host
}
