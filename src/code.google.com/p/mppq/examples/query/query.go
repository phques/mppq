// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"code.google.com/p/mppq"
	"log"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	waitFor, err := time.ParseDuration("3s")
	check(err)

	serviceDefs, err := mppq.QueryService("androidPush", waitFor)
	check(err)

	log.Printf("got %d services definitions\n", len(serviceDefs))
	for sdef := range serviceDefs {
		log.Printf("%q\n", sdef)
	}
}
