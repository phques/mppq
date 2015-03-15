package main

import (
	"fmt"
	"github.com/phques/mppq/examples/libprovider/provider"
	_ "github.com/phques/mppq/examples/libprovider/provider/gen"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
)

func main() {
	app.Run(app.Callbacks{Start: start, Stop: stop})
}

func start() {
	fmt.Println("main.start")
	java.Init()

	provider.Start()
	provider.Register("androidPush")
}

func stop() {
	fmt.Println("main.stop")
}
