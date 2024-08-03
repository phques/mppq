// Package provider
package provider

import (
	"log"
	"net"
	"net/http"
	"strconv"
)

var (
	httpListener   net.Listener
	httpListenPort int
)

func StartHTTP() error {
	// listen, on system assigned port
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Printf("Failed to HTTP Server : %v\n", err)
		return err
	}

	// save listener
	httpListener = ln

	// save listener port
	addr := ln.Addr()
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		httpListener.Close()
		httpListener = nil
		log.Printf("Failed to HTTP Server, error getting listener port : %v\n", err)
		return err
	}
	httpListenPort, err = strconv.Atoi(port)
	if err != nil {
		httpListener.Close()
		httpListener = nil
		log.Printf("Failed to HTTP Server, error getting listener port : %v\n", err)
	}

	// start serving
	go http.Serve(ln, nil)
	return nil
}
