// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	//"fmt"
	"log"
	"net"
)

const (
	// we will listen for udp messages on this multicast address
	multicastUdpAddrStr = "239.255.0.13:1440"

	whosthereStr = "mppq.whosthere?"
)

var (
	// setup int init()
	multicastUdpAddr *net.UDPAddr = nil
)

func init() {
	// resolve multicast udp address
	var err error
	multicastUdpAddr, err = net.ResolveUDPAddr("udp4", multicastUdpAddrStr)
	if err != nil {
		log.Fatal("failed to resolve multicast udp address. ", err)
	}
}

//---------

// this is what we answer back from a 'whos there' msg
// nb: caller should get our IP from the returned udp packet
type ServiceDef struct {
	ServiceName  string
	ProviderName string
	HostPort     int
	Protocol     string // ie "jsonrpcv1"
}

// holds the data / remote address when a udp msg is received
type UDPPacket struct {
	remoteAddr *net.UDPAddr
	data       []byte
}

//------------

// waits for incoming UDP message/data
// sends it on msgChan
func udpReadLoop(conn *net.UDPConn, msgChan chan *UDPPacket) {

	// wait for msg, send it on channel
	data := make([]byte, 4*1024)
	for {
		//fmt.Println("wait for msg")
		nbRead, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Print("error reading udp socket. ", err)
			return
		}

		// create new *copy* array/slice of proper length = nb bytes nbRead
		// (avoid keeping ref to original buffer)
		udpData := make([]byte, nbRead)
		copy(udpData, data[:nbRead])

		// create & fill a new UDPPacket
		udpPacket := &UDPPacket{remoteAddr: remoteAddr, data: udpData}

		// send it on the channel
		msgChan <- udpPacket

	}
}
