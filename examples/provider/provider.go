// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/phques/mppq"
)

func main() {

	fmt.Println("Starting provider!")

	// start mppq provider
	prov := mppq.NewProvider()
	prov.Start()

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname()
	prov.AddService(mppq.ServiceDef{
		ServiceName:  "androidPush",
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	})

	//test,  wait
	//	select {}

	delay := time.Second * 10
	fmt.Println("providing for", delay)
	time.Sleep(delay)
	prov.Stop()
	time.Sleep(1500 * time.Millisecond)
}
