package main

import (
	"arith"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func startServer() {
	arith := new(arith.Arith)

	server := rpc.NewServer()
	server.Register(arith)

	//server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	go startServer()

	conn, err := net.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	args := &arith.Args{7, 8}
	var reply int

	c := jsonrpc.NewClient(conn)

	err = c.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

}
