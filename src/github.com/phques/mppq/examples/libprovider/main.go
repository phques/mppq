package main

import (
	_ "github.com/phques/mppq/examples/libprovider/provider/gen"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
)

func main() {
	app.Run(app.Callbacks{Start: java.Init})
}
