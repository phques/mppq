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
	// we will listen for udp messages on this
	udpPort             = 1440
	udpPortStr          = ":1440"
	multicastUdpAddrStr = "239.255.0.13"
	broadcastUdpAddrStr = "255.255.255.255"

	whosthereStr = "mppq.whosthere?"
)

var (
	// setup in init()
	multicastUdpAddr *net.UDPAddr
	broadcastUdpAddr *net.UDPAddr
)

func init() {
	// resolve udp addresses
	var err error
	var addrstr string

	// resolve multicast
	addrstr = multicastUdpAddrStr + udpPortStr
	multicastUdpAddr, err = net.ResolveUDPAddr("udp4", addrstr)
	if err != nil {
		log.Fatal("failed to resolve multicast udp address. ", err)
	}

	// resolve broadcast
	addrstr = broadcastUdpAddrStr + udpPortStr
	broadcastUdpAddr, err = net.ResolveUDPAddr("udp4", addrstr)
	if err != nil {
		log.Fatal("failed to resolve broadcast udp address. ", err)
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
	// filled by Client lib when recving response from query
	RemoteAddr *net.UDPAddr `json:"RemoteAddr,omitempty"`
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
