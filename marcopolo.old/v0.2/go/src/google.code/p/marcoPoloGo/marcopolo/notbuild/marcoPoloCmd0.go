// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// marcopolo lets apps register for msgs they are interested in,
// other apps then requests for the registered app's address for a msg
package marcopolo

import (
	"errors"
	"fmt"
	"net"
	//"time"
	"encoding/json"
	"strings"
)

//--------

const (
	marcoPoloUdpPort int = 4444      // fixed / hard-coded port for now ;-p
	udpMsgBufferSize     = 1024 * 16 // max of 16k for marcoPolo msgs
	msgPrefix            = "marco.polo:"
)

//--------

// Version of marco.polo (msg)
type Version struct {
	Major int
	Minor int
}

// CmdMsg is a marco.polo command/msg, sent as a JSON string
type CmdMsg struct {
	Version    Version
	Action     string
	Name       string
	OptPayload string // optional broadcast msg payload
}

//---------

// Conn holds a udp client / server connection & address
type Conn struct {
	udpConn *net.UDPConn
	udpAddr net.UDPAddr
}

// Client connection
type ClientConn struct {
	Conn
}

// Server connection
type ServerConn struct {
	Conn
}

//---------

// Conn.open opens local udp connection, system assigned address, on 'port'
// use port=0 for sys assigned port
func (conn *Conn) open(port int) (err error) {
	localAnyAddr := net.IPv4(0, 0, 0, 0)
	conn.udpAddr = net.UDPAddr{IP: localAnyAddr, Port: port}
	conn.udpConn, err = net.ListenUDP("udp4", &conn.udpAddr)
	return err
}

// ------

// ServerConn.Open opens an UDP Conn on local marco.polo port
func (conn *ServerConn) Open() (err error) {
	return conn.open(marcoPoloUdpPort)
}

// ClientConn.Open opens an UDP Conn on a sys allocated local address&port
func (conn *ClientConn) Open() (err error) {
	return conn.open(0)
}

//-----

// Close closes the udp connection
func (conn *Conn) Close() {
	if conn.udpConn != nil {
		conn.udpConn.Close()
	}
}

//--------

// Conn.RecvMsg waits for & receives a marcopolo.CmdMsg
func (conn *ServerConn) RecvCmdMsg() (cmdMsg CmdMsg, remoteAddr *net.UDPAddr, err error) {
	// recv / wait for an UDP message, nbRead into buffer
	data := make([]byte, udpMsgBufferSize)
	nbRead, remoteAddr, err := conn.udpConn.ReadFromUDP(data)
	if err != nil {
		return
	}

	// get slice of proper length (nb bytes nbRead)
	data = data[:nbRead]

	//## debug
	fmt.Printf("nbRead %d bytes '%s'\n", nbRead, data)
	fmt.Printf("from udp (remote %s)\n", remoteAddr.String())

	// check for valid marco polo msg: does it start with "marco.polo:" ?
	if strings.HasPrefix(string(data), msgPrefix) {
		// ok valid prefix, get the JSON string (a marcopolo.CmdMsg)
		// "marco.polo:{JSON-marcopolo.CmdMsg}"
		marcoPoloMsgJson := data[len(msgPrefix):]

		//## debug
		fmt.Println("marcoPoloMsgJson:", string(marcoPoloMsgJson))

		// unmarshall json string to MarcoPoloMsg
		err = json.Unmarshal(marcoPoloMsgJson, &cmdMsg)
		if err != nil {
			// not a valid marcopolo.CmdMsg JSON object
			err = fmt.Errorf("not a valid marco.polo msg, JSON error: %q", err)
			return
		}
	} else {
		// not valid, does not start w "marco.polo:"
		err = errors.New("not a marco.polo msg")
		return
	}

	return
}
