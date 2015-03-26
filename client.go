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

// query holds the info required to execute a service query
type query struct {
	name         string
	waitFor      time.Duration
	useBroadcast bool
}

//----------

// QueryService sends a 'whosthere' for service name,
// returns array of found SerfviceDefs
func QueryService(name string, waitFor time.Duration, useBroadcast bool) ([]ServiceDef, error) {
	log.Printf("dbg QueryService useBroadcast = %v\n", useBroadcast)

	// create Client object
	q := &query{name: name, waitFor: waitFor, useBroadcast: useBroadcast}

	serviceDefs, err := q.doQuery()
	if err != nil {
		return nil, err
	}

	return serviceDefs, nil
}

//------------

// doQuery sends a 'whosthere' for service name,
// returns array of found SerfviceDefs
func (q *query) doQuery() ([]ServiceDef, error) {

	// open udp connection (any local address & port)
	var err error
	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)

	if err != nil {
		log.Println("failed to open local udp connection. ", err)
		return nil, err
	}
	defer udpConn.Close()

	// prep channel to recv messages from udp loop
	// & start udp read loop
	msgChan := make(chan *UDPPacket)
	//nb: will stop when udpConn is closed
	go udpReadLoop(udpConn, msgChan)

	// send query !
	query := whosthereStr + q.name
	if q.useBroadcast {
		//## for Window8/8.1, cant recv multicast, send broadcast
		udpConn.WriteToUDP([]byte(query), &broadcastUdpAddr)
	} else {
		// multicast
		udpConn.WriteToUDP([]byte(query), &multicastUdpAddr)
	}

	// loop until timeout, gathering recvd ServiceDef
	//##PQ TODO: should resend query every 1sec or something ! (?)
	//## but then we would recv multiple entries from same services !
	//## just let user call it multiple time to handle this !
	var serviceDefs []ServiceDef
	timer := time.NewTimer(q.waitFor)
	done := false
	for !done {
		select {
		case <-timer.C:
			// time is over, we're done
			done = true

		case udpPacket := <-msgChan:
			// received udp reponse packet, processs it
			serviceDef := q.processUdpPacket(udpPacket)
			if serviceDef != nil {
				serviceDefs = append(serviceDefs, *serviceDef)
			}
		}
	}

	// client.udpConn will close on return, so udpReadLoop() will stop
	return serviceDefs, nil
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
