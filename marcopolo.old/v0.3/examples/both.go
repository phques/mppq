// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phques/mppq"
)

var useBroadcast = flag.Bool("useBroadcast", false, "use broadcast instead of multicast (Win8)")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// ------ start provider -------
	fmt.Println("Starting provider!")

	// start mppq provider
	prov, err := mppq.NewProvider()
	check(err)

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname()
	prov.AddService(mppq.ServiceDef{
		ServiceName:  "androidPush",
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	})

	//------ ('client') query for service -----
	flag.Parse()

	waitFor, err := time.ParseDuration("3s")
	check(err)

	// query for provided service
	serviceDefs, err := mppq.QueryService("androidPush", waitFor, *useBroadcast)
	check(err)

	// print found serviceDefs
	fmt.Printf("-----\ngot %d services definitions\n", len(serviceDefs))
	for _, sdef := range serviceDefs {
		fmt.Printf("%v\n", sdef)
	}

	//## test debug, wait
	time.Sleep(700 * time.Millisecond)

	// stop provider
	prov.Stop()
}
