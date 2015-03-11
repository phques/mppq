// Package hi provides a function for saying hello.
package provider

import (
	"code.google.com/p/mppq"
	"fmt"
	"os"
)

func Start(serviceName string) {

	fmt.Println("provider.Start")

	// start mppq provider

	prov := mppq.NewProvider()
	go prov.MarcoPoloLoop()

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname()
	prov.AddService <- mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	}

	// wait
	select {}
}
