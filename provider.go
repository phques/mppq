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
	quit     chan struct{}
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

// NewProvider creates a new Provider
func NewProvider() *Provider {
	prov := new(Provider)
	prov.quit = make(chan struct{})
	prov.addSrvCh = make(chan ServiceDef)
	prov.delSrvCh = make(chan ServiceDef)
	prov.services = make(map[string]ServiceDef)

	return prov
}

// Start opens a UDP connection and starts the 'marcoPolo' loop listening for queries
func (prov *Provider) Start() error {
	// open udp connection
	conn, err := openUdpConn()
	if err != nil {
		return err
	}

	// launch goroutine marcoPoloLoop, will close conn
	started := make(chan bool)
	go prov.marcoPoloLoop(conn, started)

	// wait for it to be ready before returning
	<-started
	return nil
}

// IsRunning returns true if the provider was started
func (prov *Provider) IsRunning() bool {
	if prov.quit == nil {
		return false
	}

	select {
	case <-prov.quit:
		return false
	default:
		return true
	}
}

// Stop signals the provider to stop running / listening for queries
func (prov *Provider) Stop() error {
	if !prov.IsRunning() {
		return errors.New("Stop, Provider is not running")
	}

	close(prov.quit)
	return nil
}

// AddService adds a known service to the provider
//nb: provider must be Start)ed
func (prov *Provider) AddService(service ServiceDef) error {
	if !prov.IsRunning() {
		return errors.New("AddService, Provider is not running")
	}
	prov.addSrvCh <- service
	return nil
}

// DelService removes a known service from the provider
//nb: provider must be Start)ed
func (prov *Provider) DelService(service ServiceDef) error {
	if !prov.IsRunning() {
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

func (prov *Provider) processUdpPacket(conn *net.UDPConn, packet *UDPPacket) {

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
	if _, err := conn.WriteToUDP(response, packet.remoteAddr); err != nil {
		log.Println("error sending back udp response. ", err)
	}

}

// main loop (goroutine), will close conn on exit
// waits for add/del service commands and
// for data (mppq serfvice queries) from the UDP conn
func (prov *Provider) marcoPoloLoop(conn *net.UDPConn, started chan<- bool) {
	defer conn.Close()

	udpChan := make(chan *UDPPacket)
	quitChan := make(chan struct{})
	go udpReadLoop(conn, udpChan, quitChan)

	// signal that we are ready
	started <- true

	var chanReadOk = true
	for chanReadOk {
		select {
		case _ = <-prov.quit:
			log.Println("marcoPoloLoop recvd quit")
			chanReadOk = false
			close(quitChan)

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
				prov.processUdpPacket(conn, udpPacket)
			}
		}
	}
}
