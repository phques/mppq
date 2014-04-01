// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	"fmt"
	"log"
	"net"
)

//---------

// a marcopolopq provider,
// listens on udp for 'whosthere' msgs,
// answering with info about service if we match the queried service
type Provider struct {
	Quit       chan bool
	AddService chan ImThere
	DelService chan ImThere
	services   map[string]ImThere
}

//---------

func openMUdpConn() *net.UDPConn {
	// resolve multicast udp address
	mcaddr, err := net.ResolveUDPAddr("udp4", multicastUdpAddr)
	if err != nil {
		log.Fatal("failed to resolve multicast udp address. ", err)
	}

	// open listen connection on default system interface
	mudpConn, err := net.ListenMulticastUDP("udp4", nil, mcaddr)
	if err != nil {
		log.Fatal("failed to open multicast udp listen connection. ", err)
	}

	return mudpConn
}

//---------------

func NewProvider() *Provider {
	mppq := new(Provider)
	mppq.Quit = make(chan bool)
	mppq.AddService = make(chan ImThere)
	mppq.DelService = make(chan ImThere)
	mppq.services = make(map[string]ImThere)

	return mppq
}

func (mppq *Provider) MarcoPoloLoop() {
	udpConn := openMUdpConn()
	defer udpConn.Close()

	udpChan := make(chan *UDPPacket)
	go udpReadLoop(udpConn, udpChan)

	var chanReadOk = true
	for chanReadOk {
		select {
		case _ = <-mppq.Quit:
			chanReadOk = false

		case jsonAdd, chanReadOk := <-mppq.AddService:
			if chanReadOk {
				mppq.addService(jsonAdd)
			}

		case jsonDel, chanReadOk := <-mppq.DelService:
			if chanReadOk {
				mppq.delService(jsonDel)
			}

		case udpPacket, chanReadOk := <-udpChan:
			if chanReadOk {
				fmt.Printf("MarcoPoloLoop, read %d bytes '%s' from udp (remote %s)\n",
					len(udpPacket.data), udpPacket.data, udpPacket.remoteAddr.String())
			}
		}
	}
}

//------------

func (mppq *Provider) addService(service ImThere) {
	mppq.services[service.ServiceName] = service
}

func (mppq *Provider) delService(service ImThere) {
	delete(mppq.services, service.ServiceName)
}
