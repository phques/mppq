package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/phques/mppq"
)

type RpcServerExample struct {
	data int
}

type RpcArgs struct {
	Arg string
}
type RpcReply struct {
	Reply string
}

func (s *RpcServerExample) Test(args RpcArgs, reply *RpcReply) error {
	println("RpcServerExample, received request:", args.Arg, ", s.data:", s.data)
	reply.Reply = "hello " + args.Arg
	return nil
}

func main() {
	// start an rpc server for RpcServerExample
	println("starting rpc server...")
	rpcServerExample := &RpcServerExample{data: 12}
	rpc.Register(rpcServerExample)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	// create service def for rpc server
	serviceDefs := []mppq.ServiceDef{
		mppq.NewLocalServiceDef("RpcServerExample", 1234, "http rpc"),
	}

	// start mppq server for 30 seconds
	println("starting mppq server, for 30 seconds...")

	ctx, cancel := context.WithCancel(context.Background())
	mppqServer := mppq.NewMppqServer(serviceDefs)
	go mppqServer.Serve(ctx)

	time.Sleep(60 * time.Second)

	println("stopping mppq server after 60 seconds...")
	cancel()

	// wait for server to stop
	time.Sleep(time.Second)
}
