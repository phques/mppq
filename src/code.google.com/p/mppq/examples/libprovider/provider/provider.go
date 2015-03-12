// Package hi provides a function for saying hello.
package provider

import (
	"code.google.com/p/mppq"
	"fmt"
	"os"
)

var (
    provider *mppq.Provider
)

func Start() {

	fmt.Println("provider.Start")

	// create/start mppq provider
	provider = mppq.NewProvider()
	go provider.MarcoPoloLoop()

}

func Register(serviceName string) {

	fmt.Println("provider.Register")

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname() // returns 'localhost' on my Nexus 7
	hostname = "PQ Nexus 7"
	provider.AddService <- mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	}
}

