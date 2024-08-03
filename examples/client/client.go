package main

import (
	"fmt"
	"log"

	"github.com/phques/mppq"
)

func main() {
	println("starting mppq client... for 10 seconds...")
	mppqResponse, err := mppq.QueryServiceDefs(10)

	if err != nil {
		println("failed to query service defs. ", err)
	} else {
		fmt.Printf("got valid MppqResponse, header: %s\nserviceDefs:\n", mppqResponse.Header)
		for _, s := range mppqResponse.ServiceDefs {
			log.Printf("service: %v\n", s)
		}
	}

}
