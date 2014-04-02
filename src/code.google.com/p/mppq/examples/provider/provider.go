// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"code.google.com/p/mppq"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")

	// start mppq provider
	prov := mppq.NewProvider()
	go prov.MarcoPoloLoop()

	// register a service (provider main loop must be running)
	prov.AddService <- mppq.ServiceDef{
		ServiceName:  "androidPush",
		ProviderName: "moue",
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	}

	// wait
	select {}
}
