package main

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
	CheckIp(Ip)
	// refuse the connection if checkIp returns false
	reader := bufio.NewReader(conn)
	Host := ReturnHost(reader)
	fmt.Println(Host)
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
		io.Copy(Dest, conn)
		done <- struct{}{}
	}()

	<-done // wait for one direction to fin
}

// This is a really basic function it can give a common ip if some users are using ip
// from the same public wifi SO that can conflict for bandwidth
func ReturnIp(conn net.Conn) string {
	clientIp := conn.RemoteAddr().String()
	return clientIp
}
