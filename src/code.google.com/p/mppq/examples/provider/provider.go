// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"code.google.com/p/mppq"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	mppq := mppq.NewProvider()
	mppq.MarcoPoloLoop()
}
