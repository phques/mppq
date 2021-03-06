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

// ServiceDef describes a service the provider offers.
// this is what we answer back from a 'whos there' msg
type ServiceDef struct {
	ServiceName  string // the name of the service provided
	ProviderName string // the name of the provider
	HostPort     int    // port on which the service is available
	Protocol     string // app arbitrary protocol name, ie "jsonrpcv1"
	// filled by Client lib when recving response from query
	RemoteIP *net.IP `json:"RemoteIP,omitempty"`
}

// UDPPacket holds the data & remote address of a received udp msg
type UDPPacket struct {
	remoteAddr *net.UDPAddr
	data       []byte
}

//------------

// readFromUDP starts a goroutine that feeds the returned channel
// with data read from the UDP connection
// the data channel will be closed automatically when the goroutine stops
func readFromUDP(conn *net.UDPConn, quitChan <-chan struct{}) chan *UDPPacket {
	ch := make(chan *UDPPacket)

	go func() {
		defer close(ch) // close channel when we stop
		data := make([]byte, 4*1024)

		for {
			// read from UDP connection
			// can't know if err is caused by closed connection (?)
			nbRead, remoteAddr, err := conn.ReadFromUDP(data)

			select {
			case <-quitChan:
				// ok, dont display error, we were asked to stop
				log.Println("readFromUDP recvd quit (conn closed)")
				return

			default:
				if err != nil {
					// nothing on quit channel, display error before quiting
					log.Print("readFromUDP, error reading udp socket. ", err)
					return
				}
			}

			// create new *copy* array/slice of proper length = nb bytes nbRead
			// (avoid keeping ref to original buffer)
			udpData := make([]byte, nbRead)
			copy(udpData, data[:nbRead])

			// create an UDPPacket
			udpPacket := &UDPPacket{remoteAddr: remoteAddr, data: udpData}

			// send on out channel
			// here we assume that udpReadLoop will read quickly
			ch <- udpPacket
		}
	}()

	return ch
}

// udpReadLoop waits for incoming UDP message/data
// and sends it on msgChan.
// Close conn to stop, close quitChan before though!
func udpReadLoop(conn *net.UDPConn, msgChan chan<- *UDPPacket, quitChan <-chan struct{}) {

	dataCh := readFromUDP(conn, quitChan)

	// wait for msg, send it on channel
	for {
		// check if we were asked to quit
		select {
		case <-quitChan:
			log.Println("udpReadLoop recvd quit")
			return

		case udpPacket := <-dataCh:
			if udpPacket == nil {
				// dataCh was closed,
				// something must have happened with UDPConn
				return

			} else {
				// send it on the channel
				msgChan <- udpPacket
			}
		}
	}
}
