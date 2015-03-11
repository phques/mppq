package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
    _ "code.google.com/p/mppq/examples/libprovider/provider/go_provider"
)

func main() {
	app.Run(app.Callbacks{Start: java.Init})
}

