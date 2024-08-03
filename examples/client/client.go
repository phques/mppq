package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/phques/mppq"
)

type RpcArgs struct {
	Arg string
}
type RpcReply struct {
	Reply string
}

func main() {
	// query for provided services (want rpc server)
	println("starting mppq client... for 10 seconds...")
	mppqResponse, err := mppq.QueryServiceDefs(10)

	if err != nil {
		println("failed to query service defs. ", err)
	} else {
		fmt.Printf("got valid MppqResponse, header: %s\nserviceDefs:\n", mppqResponse.Header)
		for _, s := range mppqResponse.ServiceDefs {
			log.Printf("service: %v\n", s)
		}
	}

	// get rpc server info from service def (we expect only one)
	serviceDef := mppqResponse.ServiceDefs[0]
	if serviceDef.ServiceName != "RpcServerExample" {
		log.Fatal("expected service name 'RpcServerExample', got: ", serviceDef.ServiceName)
	}
	rpcServer := serviceDef.HostIP + ":" + fmt.Sprintf("%d", serviceDef.HostPort)

	// create rpc client
	client, err := rpc.DialHTTP("tcp", rpcServer)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// call rpc server
	println("\n\ncalling rpc server @", rpcServer, "(info from mppqResponse)")

	args := &RpcArgs{Arg: "Phil"}
	rpcReply := &RpcReply{Reply: "reply"}

	err = client.Call("RpcServerExample.Test", args, &rpcReply)

	// show results
	if err != nil {
		log.Fatal("arith error:", err)
	} else {
		println("RpcServerExample, received reply:", rpcReply.Reply)
	}
}
