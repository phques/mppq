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

	// create/start mppq provider
	prov, err := mppq.NewProvider()

	if err == nil {
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

		// 'provide' for 10secs
		delay := time.Second * 10
		fmt.Println("providing for", delay)
		time.Sleep(delay)
		// then stop
		prov.Stop()

		// wait (debug) for goroutines to stop
		time.Sleep(500 * time.Millisecond)
		// debug
		//panic(nil)
	} else {
		fmt.Printf("error %v\n", err)
	}
}
