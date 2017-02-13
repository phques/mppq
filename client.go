// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"time"
)

// Query is used to query for services
type Query interface {
	Stop()
	Done() bool
	ServiceCh() chan *ServiceDef
}

// query holds the info required to execute a service query
type query struct {
	name         string
	useBroadcast bool
	// for goroutine loop
	done      chan struct{} // close to stop / closed when done
	serviceCh chan *ServiceDef
}

//----------

// QueryService sends 'whosthere' queries for service name,
// returns array of *ALL* received SerfviceDefs (incl. duplicates)
func QueryService(name string, waitFor time.Duration, useBroadcast bool) ([]ServiceDef, error) {
	log.Printf("dbg QueryService useBroadcast = %v\n", useBroadcast)

	// create Query object
	q, err := NewQuery(name, useBroadcast)
	if err != nil {
		return nil, err
	}

	// loop util timeout (waitFor), adding found services to list
	done := false
	timeout := time.NewTimer(waitFor)
	var services []ServiceDef
	for !done {
		select {
		case <-timeout.C:
			q.Stop() // stop Query loop
			done = true
		case s := <-q.ServiceCh():
			services = append(services, *s)
		}
	}

	return services, nil
}

//------------

// NewQuery creates and starts a Query
func NewQuery(name string, useBroadcast bool) (Query, error) {
	// create Query object
	q := &query{name: name, useBroadcast: useBroadcast}
	q.done = make(chan struct{})
	q.serviceCh = make(chan *ServiceDef)

	err := q.start()
	if err != nil {
		return nil, err
	}

	return q, nil
}

// ServiceCh returns the ServiceDef channel
func (q *query) ServiceCh() chan *ServiceDef {
	return q.serviceCh
}

// Stop stops the Query
func (q *query) Stop() {
	if !q.Done() { // not completely safe !& (race?)
		close(q.done)
	}
}

// Done returns true if Query was stopped
func (q *query) Done() bool {
	select {
	case <-q.done:
		return true
	default:
		return false
	}
}

// -------------

// Start launches the doQueryLoop goroutine with an UDPConn
func (q *query) start() error {
	// open udp connection (any local address & port)
	var err error
	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)

	if err != nil {
		log.Println("failed to open local udp connection. ", err)
		return err
	}

	// start loop
	go q.doQueryLoop(udpConn)
	return nil
}

// sendQuery sends the mppq whosThere query for a service
func (q *query) sendQuery(udpConn *net.UDPConn) {
	queryStr := whosthereStr + q.name
	log.Printf("sendQuery <%v>\n", queryStr)

	if q.useBroadcast {
		//## for Window8/8.1, cant recv multicast, send broadcast
		udpConn.WriteToUDP([]byte(queryStr), &broadcastUdpAddr)
	} else {
		// multicast
		udpConn.WriteToUDP([]byte(queryStr), &multicastUdpAddr)
	}
}

// doQueryLoop loops until timeout, sending recvd ServiceDef on channel
// normally runs in a goroutine. udpConn closed on exit
func (q *query) doQueryLoop(udpConn *net.UDPConn) {

	defer udpConn.Close()

	// prep channel to recv messages from udp loop
	// & start udp read loop
	udpReadMsgChan := make(chan *UDPPacket)
	udpReadQuitChan := make(chan struct{})
	//nb: will stop when udpConn is closed
	go udpReadLoop(udpConn, udpReadMsgChan, udpReadQuitChan)

	// send initial query !
	q.sendQuery(udpConn)
	sendRepeatDelay := time.Second
	sendTimeout := time.NewTimer(sendRepeatDelay)

	// loop until timeout, sending recvd ServiceDef on channel
	done := false
	for !done {
		select {
		case <-q.done:
			// we're done, stop loop
			log.Println("Query.doQueryLoop: recv quit")
			done = true

		case <-sendTimeout.C:
			// time to send a query again
			q.sendQuery(udpConn)
			sendTimeout.Reset(sendRepeatDelay)

		case udpPacket := <-udpReadMsgChan:
			// received udp reponse packet, processs it
			serviceDef := q.processUdpPacket(udpPacket)
			if serviceDef != nil {
				// send serviceDef,
				// note that if client does not read fast enough,
				// we will timeout above !
				q.serviceCh <- serviceDef
			}
		}
	}

	// client.udpConn will close on return, so udpReadLoop() will stop
	close(udpReadQuitChan) // signal that we have closed conn / stopping
	close(q.serviceCh)
}

// processUdpPacket
func (q *query) processUdpPacket(udpPacket *UDPPacket) *ServiceDef {
	log.Printf("recvd response [%s]", udpPacket.data)

	// did we receive a whosthere mppq query ?
	if !bytes.HasPrefix(udpPacket.data, []byte(ImhereStr)) {
		// debug (hmm, not a good idea too output received unknown data ! ;-p)
		log.Printf("received unknown response")
		return nil
	}

	// get serviceDef parameter of Imhere! "Imhere!serviceDefJson"
	serviceDefJson := udpPacket.data[len(ImhereStr):]

	// decode json ServiceDef response
	var serviceDef ServiceDef
	if err := json.Unmarshal(serviceDefJson, &serviceDef); err != nil {
		log.Println("error decoding json ServiceDef response. ", err)
		return nil
	}

	// add the remote udp address, taken from udpPacket
	remIP := udpPacket.remoteAddr.IP
	serviceDef.RemoteIP = &remIP

	return &serviceDef
}
