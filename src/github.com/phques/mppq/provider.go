// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	"bytes"
	"encoding/json"
	"errors"
	//"fmt"
	"log"
	"net"
)

//---------

// a marcopolopq provider,
// listens on udp for 'whosthere' msgs,
// answering with info about service if we match the queried service
type Provider struct {
	run      bool
	udpConn  *net.UDPConn
	quit     chan bool
	addSrvCh chan ServiceDef
	delSrvCh chan ServiceDef
	services map[string]ServiceDef
}

//---------

func openUdpConn() (*net.UDPConn, error) {
	// open listen connection on default system interface
	// NB: on Win8/8.1, we *can* use multicast to listen,
	//     it will work if we *send broadcast* !
	udpConn, err := net.ListenMulticastUDP("udp4", nil, &multicastUdpAddr)
	if err != nil {
		log.Println("failed to open multicast udp listen connection. ", err)
	}

	return udpConn, err
}

//---------------

func NewProvider() *Provider {
	prov := new(Provider)
	prov.run = false
	prov.quit = make(chan bool)
	prov.addSrvCh = make(chan ServiceDef)
	prov.delSrvCh = make(chan ServiceDef)
	prov.services = make(map[string]ServiceDef)

	return prov
}

func (prov *Provider) Start() error {
	// open udp connection
	conn, err := openUdpConn()
	if err != nil {
		return err
	}

	// launch goroutine marcoPoloLoop, will close prov.udpConn
	started := make(chan bool)
	prov.udpConn = conn
	go prov.marcoPoloLoop(conn, started)

	// wait for it to be ready before returning
	<-started
	return nil
}

func (prov *Provider) Stop() {
	//if (!prov.stopped) {
	prov.run = false
	prov.quit <- true
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

func (prov *Provider) AddService(service ServiceDef) error {
	if !prov.run || prov.addSrvCh == nil {
		return errors.New("AddService, Provider is not running / not initialized")
	}
	prov.addSrvCh <- service
	return nil
}

func (prov *Provider) DelService(service ServiceDef) error {
	if !prov.run || prov.delSrvCh == nil {
		return errors.New("DelService, Provider is not running / not initialized")
	}
	prov.delSrvCh <- service
	return nil
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
	response := []byte(ImhereStr + string(jsonmsg))
	if _, err := prov.udpConn.WriteToUDP(response, packet.remoteAddr); err != nil {
		log.Printf("error sending back udp response. ", err)
	}

}

// main loop (goroutine)
// will close conn
func (prov *Provider) marcoPoloLoop(conn *net.UDPConn, started chan<- bool) {
	defer prov.udpConn.Close()

	udpChan := make(chan *UDPPacket)
	go udpReadLoop(prov.udpConn, udpChan)

	prov.run = true

	// signal that we are ready
	started <- true

	var chanReadOk = true
	for chanReadOk && prov.run {
		select {
		case _ = <-prov.quit:
			chanReadOk = false

		case serviceAdd, chanReadOk := <-prov.addSrvCh:
			if chanReadOk {
				prov.addService(serviceAdd)
			}

		case serviceDel, chanReadOk := <-prov.delSrvCh:
			if chanReadOk {
				prov.delService(serviceDel)
			}

		case udpPacket, chanReadOk := <-udpChan:
			if chanReadOk {
				prov.processUdpPacket(udpPacket)
			}
		}
	}
}
