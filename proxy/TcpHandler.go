package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func TcpHandler(conn net.Conn) {
	defer conn.Close()

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
