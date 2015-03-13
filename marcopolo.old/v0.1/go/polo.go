// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// main.go
package main

import (
	"fmt"
	"net"
)

func main() {

	// --- listen for marco msg on udp  ---
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 4444,
	})
	if err != nil {
		fmt.Println("error open udp socket ", err)
		return
	}

	// wait for marco
	fmt.Println("wait for marco")
	data := make([]byte, 1024)
	read, remoteAddr, err := udpConn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("error reading udp socket ", err)
		return
	}
	fmt.Printf("read %d bytes '%s' from udp (remote %s)\n", read, data[:read], remoteAddr.String())

	// --- send back polo ----

	polo := "polo|testMarcoPolo|1234"
	fmt.Println("sending back ", polo)
	_, err = udpConn.WriteToUDP([]byte(polo), remoteAddr)
	//fmt.Println("nbbytes, err : ", nbBytes, err)
}
