package main

import (
	"context"
	"time"

	"github.com/phques/mppq"
)

func main() {
	serviceDefs := []mppq.ServiceDef{
		mppq.NewLocalServiceDef("marco", 1440, "jsonrpcv1"),
		mppq.NewLocalServiceDef("polo", 1441, "jsonrpcv1"),
		mppq.NewLocalServiceDef("pippo", 1442, "jsonrpcv1"),
		mppq.NewLocalServiceDef("pluto", 1443, "jsonrpcv1"),
		mppq.NewLocalServiceDef("paperino", 1444, "jsonrpcv1"),
		mppq.NewLocalServiceDef("topolino", 1445, "jsonrpcv1"),
	}

	// start mppq server
	ctx, cancel := context.WithCancel(context.Background())

	mppqServer := mppq.NewMppqServer(serviceDefs)

	println("starting mppq server, for 30 seconds...")
	go mppqServer.Serve(ctx)

	time.Sleep(60 * time.Second)

	println("stopping mppq server after 60 seconds...")
	cancel()

	// wait for server to stop
	time.Sleep(time.Second)
}
