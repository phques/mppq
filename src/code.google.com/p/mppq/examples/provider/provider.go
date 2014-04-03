// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"code.google.com/p/mppq"
	"fmt"
	"os"
)

func main() {

	fmt.Println("Hello World!")

	// start mppq provider
	prov := mppq.NewProvider()
	go prov.MarcoPoloLoop()

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname()
	prov.AddService <- mppq.ServiceDef{
		ServiceName:  "androidPush",
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	}

	// wait
	select {}
}
