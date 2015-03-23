package main

import (
	"flag"
	"fmt"
	"github.com/phques/mppq/examples/libprovider/provider"
	_ "github.com/phques/mppq/examples/libprovider/provider/gen"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
	//	"net/http"
)

var standalone = flag.Bool("standalone", false, "running as a stand alone pgm")

func main() {
	fmt.Println("main")
	flag.Parse()
	app.Run(app.Callbacks{Start: start, Stop: stop})
}

func start() {
	fmt.Println("main.start")
	java.Init()

	if *standalone {
		provider.InitAppFilesDir("files")
		provider.Start()
		provider.Register("androidPush")
	}
}

func stop() {
	fmt.Println("main.stop")
}
