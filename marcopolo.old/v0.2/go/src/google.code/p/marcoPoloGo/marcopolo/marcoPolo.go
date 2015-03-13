// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package marcopolo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

//--------

const (
	marcoPoloUdpPort int = 4444     // marcopoloServer port, fixed / hard-coded for now ;-p
	udpMsgBufferSize     = 1024 * 4 // max of 4k for marcoPolo msgs
	cmdMsgSep            = '|'      // cmdMsg separator: cmdMessage + cmdMsgSep + jsonParam

	// marco polo command messages
	CmdRegApp   = "marco.polo.regapp"
	CmdUnregApp = "marco.polo.unregapp"
	CmdQryApp   = "marco.polo.qryapp" // resp = RespQryApp

	// client app messages (responses)
	RespMsgRegApp   = "marco.polo.resp.regapp"   // resp = Resp
	RespMsgUnregApp = "marco.polo.resp.unregapp" // resp = Resp
	RespMsgQryApp   = "marco.polo.resp.qryapp"   // resp = RespQryApp
)

var (
	// marco polo command messages & response prefixes:
	// as []byte : (cmdMsg + cmdMsgSep) ie []byte("marco.polo.regapp|")
	cmdMsgPrefixes  [][]byte
	respMsgPrefixes [][]byte

	// current API version
	version Version = Version{Major: 0, Minor: 2}

	// destination 'marcoPolo' broadcast UDP address on marcoPoloUdpPort
	broadcastAddr    = net.IPv4(255, 255, 255, 255)
	marcoPoloUdpAddr = net.UDPAddr{IP: broadcastAddr, Port: marcoPoloUdpPort}
)

//---------

// Conn holds a udp client / server connection & address
type Conn struct {
	udpConn *net.UDPConn
	//udpAddr net.UDPAddr
}

// --- messages data ---

type Version struct {
	Major int
	Minor int
}

type StdHeader struct {
	Version
}

// Resp is the common/base object for response msgs
type Resp struct {
	StdHeader
	Ok    bool
	Error string // when Ok==false
}

// AppAddr represents the address of an app
type AppAddr struct {
	IP   []byte
	Port int
}

// RespQryApp is the response msg for QryApp
type RespQryApp struct {
	Resp
	AppAddr AppAddr
}

// AppNameParam is the common parameter for marcopolo server msgs
type AppNameParam struct {
	StdHeader
	AppName string
}

// UdpPacket holds the recvd datagram (data & remot addr)
type UdpPacket struct {
	data       []byte       // recvd data
	RemoteAddr *net.UDPAddr // udp adr of caller
}

//----------

// CmdMsg holds a marco.polo server command message.
// right now, all cmds/msgs have a single AppName parameter
type CmdMsg struct {
	Cmd          string       // "marco.polo.regapp"
	AppNameParam AppNameParam // the app name parameter to the cmd
	UdpPacket    UdpPacket    // the data recvd & the udp addr of caller
}

//----------

func init() {
	// init cmdMsgPrefixes from known cmd messages
	msgs := []string{CmdRegApp, CmdUnregApp, CmdQryApp}
	addMsgPrefixes(msgs, &cmdMsgPrefixes)

	// init respMsgPrefixes from known cmd messages
	msgs = []string{RespMsgRegApp, RespMsgUnregApp, RespMsgQryApp}
	addMsgPrefixes(msgs, &respMsgPrefixes)

}

// addMsgPrefixes appends values from msgs to msgPrefixes as []byte w. '|' at end
func addMsgPrefixes(msgs []string, msgPrefixes *[][]byte) {
	for _, msg := range msgs {
		// append separator to msg
		prefix := append([]byte(msg), '|')

		// append new prefix to msgPrefixes
		*msgPrefixes = append(*msgPrefixes, prefix)
	}
}

//---------

func (resp RespQryApp) UDPAddr() net.UDPAddr {
	return net.UDPAddr{IP: resp.AppAddr.IP, Port: resp.AppAddr.Port}
}

//---------

// Compare compares two version, rets -1, 0, 1
func (v1 Version) Compare(v2 Version) int {
	if v1.Major < v2.Major {
		return -1
	}
	if v1.Major > v2.Major {
		return 1
	}
	if v1.Minor < v2.Minor {
		return -1
	}
	if v1.Minor > v2.Minor {
		return 1
	}
	return 0
}

func (v1 Version) Equals(v2 Version) bool {
	return v1.Major == v2.Major && v1.Minor == v2.Minor
}

func (v1 Version) GreatherThan(v2 Version) bool {
	return v1.Compare(v2) == 1
}

func (v1 Version) SmallerThan(v2 Version) bool {
	return v1.Compare(v2) == -1
}

// if I use this, fmt.Printf("%+v", marcopolo.RespQryApp) only displays the version !
//func (v Version) String() (str string) {
//	return fmt.Sprintf("v%d.%d", v.Major, v.Minor)
//}

//---------

func (conn *Conn) String() (str string) {
	return conn.udpConn.LocalAddr().String()
}

//---

// SendMsg sends a msg + (json) parameter
func (conn *Conn) SendMsg(msg string, param interface{}, remoteAddr *net.UDPAddr) (err error) {
	// convert reponse obj to json
	jsonStr, err := json.Marshal(param)
	if err != nil {
		return err
	}

	// full msg= msg + "|" + jsonStr
	marcoPoloMsg := fmt.Sprintf("%s%c%s", msg, cmdMsgSep, jsonStr)

	// send it
	_, err = conn.udpConn.WriteToUDP([]byte(marcoPoloMsg), remoteAddr)
	return err
}

// SendMsgToMarcoPolo sends a msg + (json) parameter to the marcopolo server
func (conn *Conn) SendMsgToMarcoPolo(msg string, param interface{}) (err error) {
	return conn.SendMsg(msg, param, &marcoPoloUdpAddr)
}

//-------

// Conn.open opens local udp connection, system assigned address, on 'port'
// use port=0 for sys assigned port
func (conn *Conn) open(port int) (err error) {
	localAnyAddr := net.IPv4(0, 0, 0, 0)
	//conn.udpAddr = net.UDPAddr{IP: localAnyAddr, Port: port}
	udpAddr := net.UDPAddr{IP: localAnyAddr, Port: port}
	conn.udpConn, err = net.ListenUDP("udp4", &udpAddr)
	return err
}

// Close closes the udp connection
func (conn *Conn) close() {
	if conn.udpConn != nil {
		conn.udpConn.Close()
	}
}

// recvMsgStr waits for & receives a marcopolo command msg string,
// it sets CmdMsg.data, CmdMsg.RemoteAddr
func (conn *Conn) recvMsgStr() (udpPacket UdpPacket, err error) {
	// recv / wait for an UDP message, read into buffer
	data := make([]byte, udpMsgBufferSize)
	nbRead, remoteAddr, err := conn.udpConn.ReadFromUDP(data)
	if err != nil {
		return
	}

	//## debug
	//fmt.Printf("nbRead %d bytes %s\n", nbRead, data[:nbRead])
	//fmt.Printf("from udp (remote %s)\n", remoteAddr.String())

	// fill 2 members of cmdMsg
	// create new array/slice of proper length = nb bytes nbRead
	newData := make([]byte, nbRead)
	copy(newData, data[:nbRead])

	udpPacket.data = newData
	udpPacket.RemoteAddr = remoteAddr
	return
}

//-----

// validateMsgStr checks data for a valid marco.polo.xxx command prefix.
// It returns two slices into data:
//  [0] = command name
//  [1] = command param JSON string
func validateMsgStr(msg []byte, prefixes [][]byte) (cmdMsgStr [2][]byte, err error) {

	// check for valid marco polo msg: does it start with "marco.polo.Xyz|" ?
	for _, prefix := range prefixes {

		if bytes.HasPrefix(msg, prefix) {
			// ok valid cmd msg prefix, get the JSON string param
			// "marco.polo.xxx|JSON-param"
			prefixLen := len(prefix)
			cmdMsgStr[0] = msg[:prefixLen-1] // strip separator
			cmdMsgStr[1] = msg[prefixLen:]
			return
		}
	}

	// not valid, does not start w "marco.polo.xx|"
	err = errors.New("not a marco.polo msg")
	return
}

// validateCmd validates a marcopolo server msg
func validateCmd(msg []byte) (cmdMsgStr [2][]byte, err error) {
	return validateMsgStr(msg, cmdMsgPrefixes)
}

// validateCmd validates a marcopolo client msg/response
func validateResp(msg []byte) (cmdMsgStr [2][]byte, err error) {
	return validateMsgStr(msg, respMsgPrefixes)
}

//------

// parseJsonMsgParam converts a message input param json string to an object
// eg: var respQryApp RespQryApp
//     parseJsonMsgParam(jsonStr, &respQryApp)
func parseJsonMsgParam(jsonStr []byte, object interface{}) (err error) {

	// unmarshall json string to AppNameParam (only param type avail!)
	if err = json.Unmarshal(jsonStr, &object); err != nil {
		// not a valid AppNameParam JSON object
		return fmt.Errorf("not a valid marco.polo msg, JSON error: %q", err)
	}

	return nil
}

//----------

// MakeRespQryApp makes a RespQryApp w. an UDPAddr response
func MakeRespQryApp(address *net.UDPAddr) (resp *RespQryApp) {
	resp = new(RespQryApp)
	resp.Version = version
	resp.Ok = true
	resp.AppAddr.Port = address.Port
	resp.AppAddr.IP = address.IP
	return resp
}

// MakeRespQryAppErr makes a RespQryApp w. an error string
func MakeRespQryAppErr(error error) (resp *RespQryApp) {
	resp = new(RespQryApp)
	resp.Version = version
	resp.Ok = false
	resp.Error = error.Error()
	return resp
}

// MakeRespRegApp makes a Resp
func MakeRespRegApp() (resp *Resp) {
	return makeResp()
}

// MakeRespUnregApp makes a Resp
func MakeRespUnregApp() (resp *Resp) {
	return makeResp()
}

// MakeRespUnregAppErr makes an error Resp
func MakeRespUnregAppErr(err error) (resp *Resp) {
	return makeRespErr(err)
}

//---

// makeResp makes a Resp
func makeResp() (resp *Resp) {
	resp = new(Resp)
	resp.Version = version
	resp.Ok = true
	return resp
}

// makeRespErr makes a Resp
func makeRespErr(err error) (resp *Resp) {
	resp = new(Resp)
	resp.Version = version
	resp.Ok = false
	resp.Error = err.Error()
	return resp
}
