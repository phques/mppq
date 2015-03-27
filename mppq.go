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
	udpPort      = 1440
	whosthereStr = "mppq.whosthere?" // query sent to find services
	ImhereStr    = "mppq.imHere!"    // response with serviceDef
)

var (
	// setup in init()
	multicastUdpAddr net.UDPAddr
	broadcastUdpAddr net.UDPAddr
)

func init() {
	//RFC 2365 - Administratively Scoped IP Multicast
	//  The IPv4 Local Scope -- 239.255.0.0/16
	multicastUdpAddr = net.UDPAddr{IP: net.IPv4(239, 255, 0, 13), Port: udpPort}

	broadcastUdpAddr = net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: udpPort}
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
	RemoteIP *net.IP `json:"RemoteIP,omitempty"`
}

// holds the data / remote address when a udp msg is received
type UDPPacket struct {
	remoteAddr *net.UDPAddr
	data       []byte
}

//------------

// waits for incoming UDP message/data
// sends it on msgChan,
// send a value to stopChan before closing conn (stopChan should be buffered so we can use len())
func udpReadLoop(conn *net.UDPConn, msgChan chan<- *UDPPacket, stopChan <-chan bool) {

	// wait for msg, send it on channel
	data := make([]byte, 4*1024)
	for {
		// can't know if err is caused by closed connection (?)
		nbRead, remoteAddr, err := conn.ReadFromUDP(data)

		if err != nil {
			// if we get a value on stopChan we'll assume that the connection was closed
			lenChan := len(stopChan)
			if lenChan > 0 {
				// ok, dont display error, we were asked to stop
				return
			}
			// nothing on stop channel, display error before quiting
			log.Print("udpReadLoop, error reading udp socket. ", err)
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
