// MarcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package mppq

import (
	"log"
	"net"
	"time"
)

type Client struct {
	name    string
	waitFor time.Duration

	udpConn     *net.UDPConn // created/opened only during query func
	serviceDefs []ServiceDef
}

//----------

func QueryService(name string, waitFor time.Duration) ([]ServiceDef, error) {
	client := &Client{name: name, waitFor: waitFor}

	err := client.doQuery()
	if err != nil {
		return nil, err
	}

	return client.serviceDefs, nil
}

//------------

func (client *Client) doQuery() error {

	// open udp connection (any local address & port)
	localUdpAddr, err := net.ResolveUDPAddr("udp4", ":0")
	if err != nil {
		log.Println("error ResolveUDPAddrt local", err)
		return err
	}

	client.udpConn, err = net.ListenUDP("udp4", localUdpAddr)
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
	client.udpConn.WriteToUDP([]byte(query), multicastUdpAddr)

	// loop until timeout, gathering recvd ServiceDef
	//##PQ TODO: should resend query every 1sec or something ! (?)
	//## but then we would recv multiple entries from same services !
	timer := time.NewTimer(client.waitFor)
	done := false
	for !done {
		select {
		case <-timer.C:
			done = true

		case udpPacket := <-msgChan:
			client.processUdpPacket(udpPacket)
		}
	}

	// client.udpConn will close on return, so udpReadLoop() will stop
	return nil
}

func (client *Client) processUdpPacket(udpPacket *UDPPacket) {
	log.Printf("recvd response [%s]", udpPacket.data)
	//## PQ TODO: parse json response & save
	//client.serviceDefs= append(client.serviceDefs, serviceDef)
}
