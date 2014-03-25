package main

import (
	"arith"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func startJSONServerOrig() {
	arith := new(arith.Arith)

	server := rpc.NewServer()
	server.Register(arith)

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// simpler version
func startJSONServer() {

	rpc.Register(new(arith.Arith))

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()
	fmt.Printf("listener address: %v\n", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		jsonrpc.ServeConn(conn)
	}
}

func listenUDP() {
	// --- listen for udp msg ---
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 4444,
	})
	if err != nil {
		log.Fatal("error open udp socket ", err)
	}
	defer udpConn.Close()

	// wait for msg
	data := make([]byte, 1024*2)
	for {
		fmt.Println("wait for UDP msg")
		nbRead, remoteAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			log.Fatal("error reading udp socket ", err)
		}
		udpData := data[:nbRead]

		fmt.Printf("read %d bytes '%s' from udp (remote %s)\n", nbRead, udpData, remoteAddr.String())
		if string(udpData) == "quit" {
			break
		}

		// did we get a mppq.whosthere for androidpush or * ?

		//fmt.Println("wait a sec...")
		//time.Sleep(1 * time.Second)
	}
}

func main() {
	startJSONServer()
}
