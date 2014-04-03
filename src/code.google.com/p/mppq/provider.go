// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"log"
	"net"
)

//---------

// a marcopolopq provider,
// listens on udp for 'whosthere' msgs,
// answering with info about service if we match the queried service
type Provider struct {
	run        bool
	udpConn    *net.UDPConn
	Quit       chan bool
	AddService chan ServiceDef
	DelService chan ServiceDef
	services   map[string]ServiceDef
}

//---------

func openUdpConn() *net.UDPConn {
	// open listen connection on default system interface
	// NB: on Win8/8.1, we *can* use multicast to listen,
	//     it will work if we *send broadcast* !
	udpConn, err := net.ListenMulticastUDP("udp4", nil, &multicastUdpAddr)
	if err != nil {
		log.Fatal("failed to open multicast udp listen connection. ", err)
	}

	return udpConn
}

//---------------

func NewProvider() *Provider {
	prov := new(Provider)
	prov.run = false
	prov.Quit = make(chan bool)
	prov.AddService = make(chan ServiceDef)
	prov.DelService = make(chan ServiceDef)
	prov.services = make(map[string]ServiceDef)

	return prov
}

func (prov *Provider) MarcoPoloLoop() {
	// open udp connection
	prov.udpConn = openUdpConn()
	defer prov.udpConn.Close()

	udpChan := make(chan *UDPPacket)
	go udpReadLoop(prov.udpConn, udpChan)

	prov.run = true

	var chanReadOk = true
	for chanReadOk && prov.run {
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
				prov.processUdpPacket(udpPacket)
			}
		}
	}
}

func (prov *Provider) Stop() {
	//if (!prov.stopped) {
	prov.run = false
	prov.Quit <- true
	//close(prov.Quit)
	//}
}

func (prov *Provider) Close() {
	prov.Stop()

	// prov.udpconn closed by defer in MarcoPoloLoop
	/*
		stopped    bool
		udpConn    *net.UDPConn
		Quit       chan bool
		AddService chan ServiceDef
		DelService chan ServiceDef
		services   map[string]ServiceDef
	*/
}

//------------

func (prov *Provider) addService(service ServiceDef) {
	prov.services[service.ServiceName] = service
}

func (prov *Provider) delService(service ServiceDef) {
	delete(prov.services, service.ServiceName)
}

func (prov *Provider) processUdpPacket(packet *UDPPacket) {

	// did we receive a whosthere mppq query ?
	if !bytes.HasPrefix(packet.data, []byte(whosthereStr)) {
		// debug (hmm, not a good idea too output received unknown data ! ;-p)
		log.Printf("received unknown msg [%s]", packet.data)
		return
	}

	// get serviceName parameter of whosthere "whosthere?serviceName"
	serviceName := string(packet.data[len(whosthereStr):])
	// debug
	log.Printf("received whosthere for [%s]", serviceName)

	// lookup serviceName in our registered services
	serviceDef, ok := prov.services[serviceName]
	if !ok {
		log.Printf(" .. not registered")
		return
	}

	// got it ! create json response
	jsonmsg, err := json.Marshal(serviceDef)
	if err != nil {
		// ooops ?
		log.Fatal("error json marshaling serviceDef. ", err)
	}

	// send back json response to udp sender
	if _, err := prov.udpConn.WriteToUDP(jsonmsg, packet.remoteAddr); err != nil {
		log.Printf("error sending back udp response. ", err)
	}

}
