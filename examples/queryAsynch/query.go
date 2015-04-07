// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/phques/mppq"
)

var useBroadcast = flag.Bool("useBroadcast", false, "use broadcast instead of multicast (Win8)")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loop(q *mppq.Query) {
	delay := time.Second * 5
	timeout := time.NewTimer(delay)
	fmt.Println("querying for ", delay)
	for {
		select {
		case service := <-q.ServiceCh:
			fmt.Println("got serfvice :", service)
		case <-timeout.C:
			fmt.Println("stopping")
			q.Stop()
			return
		}
	}
}

func main() {
	flag.Parse()

	query := mppq.NewQuery("androidPush", *useBroadcast)
	check(query.Start())

	loop(query)

	//## test debug
	time.Sleep(700 * time.Millisecond)
}
