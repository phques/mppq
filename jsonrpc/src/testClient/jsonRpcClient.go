package main

import (
	"arith"
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

func main() {

	client, err := jsonrpc.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	args := &arith.Args{7, 8}
	var reply int

	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)
}
