package main

import (
	"fmt"
	"log"
	"net"
)

func startServer() {

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	bytes := make([]byte, 1024)
	nbRecv, err := conn.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recvd %d bytes :\n%s\n", nbRecv, bytes[:nbRecv])
}

func main() {
	startServer()
}
