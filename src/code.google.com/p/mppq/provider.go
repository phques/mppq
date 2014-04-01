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
	AddService chan ServiceDef
	DelService chan ServiceDef
	services   map[string]ServiceDef
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
	prov := new(Provider)
	prov.Quit = make(chan bool)
	prov.AddService = make(chan ServiceDef)
	prov.DelService = make(chan ServiceDef)
	prov.services = make(map[string]ServiceDef)

	return prov
}

func (prov *Provider) MarcoPoloLoop() {
	udpConn := openMUdpConn()
	defer udpConn.Close()

	udpChan := make(chan *UDPPacket)
	go udpReadLoop(udpConn, udpChan)

	var chanReadOk = true
	for chanReadOk {
		select {
		case _ = <-prov.Quit:
			chanReadOk = false

		case jsonAdd, chanReadOk := <-prov.AddService:
			if chanReadOk {
				prov.addService(jsonAdd)
			}

		case jsonDel, chanReadOk := <-prov.DelService:
			if chanReadOk {
				prov.delService(jsonDel)
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

func (prov *Provider) addService(service ServiceDef) {
	prov.services[service.ServiceName] = service
}

func (prov *Provider) delService(service ServiceDef) {
	delete(prov.services, service.ServiceName)
}
