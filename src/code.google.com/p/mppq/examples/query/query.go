// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"code.google.com/p/mppq"
	"flag"
	"log"
	"time"
)

var useBroadcast = flag.Bool("useBroadcast", false, "use broadcast instead of multicast (Win8)")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	waitFor, err := time.ParseDuration("3s")
	check(err)

	serviceDefs, err := mppq.QueryService("androidPush", waitFor, *useBroadcast)
	check(err)

	log.Printf("got %d services definitions\n", len(serviceDefs))
	for _, sdef := range serviceDefs {
		log.Printf("%v\n", sdef)
	}
}
