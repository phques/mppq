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

type Client struct {
	name         string
	waitFor      time.Duration
	useBroadcast bool

	udpConn     *net.UDPConn // created/opened only during query func
	serviceDefs []ServiceDef
}

//----------

func QueryService(name string, waitFor time.Duration, useBroadcast bool) ([]ServiceDef, error) {
	log.Printf("dbg QueryService useBroadcast = %v\n", useBroadcast)

	// create Client object
	client := &Client{name: name, waitFor: waitFor, useBroadcast: useBroadcast}

	if err := client.doQuery(); err != nil {
		return nil, err
	}

	return client.serviceDefs, nil
}

//------------

func (client *Client) doQuery() error {

	// open udp connection (any local address & port)
	var err error
	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}

	client.udpConn, err = net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		log.Println("failed to open local udp connection. ", err)
		return err
	}
	defer client.udpConn.Close()

	// prep channel to recv messages from udp loop
	// & start udp read loop
	msgChan := make(chan *UDPPacket)
	go udpReadLoop(client.udpConn, msgChan)

	// send query !
	query := whosthereStr + client.name
	if client.useBroadcast {
		//## for Window8/8.1, cant recv multicast, send broadcast
		client.udpConn.WriteToUDP([]byte(query), &broadcastUdpAddr)
	} else {
		// multicast
		client.udpConn.WriteToUDP([]byte(query), &multicastUdpAddr)
	}

	// loop until timeout, gathering recvd ServiceDef
	//##PQ TODO: should resend query every 1sec or something ! (?)
	//## but then we would recv multiple entries from same services !
	//## just let user call it multiple time to handle this !
	timer := time.NewTimer(client.waitFor)
	done := false
	for !done {
		select {
		case <-timer.C:
			// time is over, we're done
			done = true

		case udpPacket := <-msgChan:
			// received udp reponse packet, processs it
			client.processUdpPacket(udpPacket)
		}
	}

	// client.udpConn will close on return, so udpReadLoop() will stop
	return nil
}

func (client *Client) processUdpPacket(udpPacket *UDPPacket) {
	log.Printf("recvd response [%s]", udpPacket.data)

	// did we receive a whosthere mppq query ?
	if !bytes.HasPrefix(udpPacket.data, []byte(ImhereStr)) {
		// debug (hmm, not a good idea too output received unknown data ! ;-p)
		log.Printf("received unknown response")
		return
	}

	// get serviceDef parameter of Imhere! "Imhere!serviceDefJson"
	serviceDefJson := udpPacket.data[len(ImhereStr):]

	// decode json ServiceDef response
	var serviceDef ServiceDef
	if err := json.Unmarshal(serviceDefJson, &serviceDef); err != nil {
		log.Printf("error decoding json ServiceDef response. ", err)
		return
	}

	// add the remote udp address, taken from udpPacket
	remIP := udpPacket.remoteAddr.IP
	serviceDef.RemoteIP = &remIP

	client.serviceDefs = append(client.serviceDefs, serviceDef)
}
