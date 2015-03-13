// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package marcopolo

import (
	"fmt"
	"net"
)

// Server connection
type ServerConn struct {
	Conn
}

// ------

// ServerConn.Open opens an UDP Conn on local marco.polo port
func (conn *ServerConn) Open() (err error) {
	return conn.open(marcoPoloUdpPort)
}

// ServerConn.Close closes the UDP Conn
func (conn *ServerConn) Close() {
	conn.close()
}

//-----

// RecvCmdMsg receives a command string, validates & parses it then setups CmdMsg
func (conn *ServerConn) RecvCmdMsg() (cmdMsg CmdMsg, err error) {
	// get msg as strings
	var udpPacket UdpPacket
	if udpPacket, err = conn.recvMsgStr(); err != nil {
		return
	}

	// check for valid marco.polo command, gets data split into cmd & param
	var cmdMsgStrs [2][]byte
	if cmdMsgStrs, err = validateCmd(udpPacket.data); err != nil {
		return
	}

	// parse json AppName param
	var appNameParam AppNameParam
	if err = parseJsonMsgParam(cmdMsgStrs[1], &appNameParam); err != nil {
		return
	}

	// check version from param
	if appNameParam.Version != version {
		err = fmt.Errorf("got cmd version '%v', expecting '%v'", appNameParam.Version, version)
		return
	}

	// done, setup CmdMsg
	cmdMsg.Cmd = string(cmdMsgStrs[0])
	cmdMsg.AppNameParam = appNameParam
	cmdMsg.UdpPacket = udpPacket
	return
}

//---------

// SendRespRegApp sends a RespRegApp msg back to a client
func (conn *ServerConn) SendRespRegApp(clientDestAddr *net.UDPAddr) (err error) {
	resp := MakeRespRegApp()
	return conn.SendMsg(RespMsgRegApp, resp, clientDestAddr)
}

// SendRespUnregApp sends a RespUnregApp msg back to a client
func (conn *ServerConn) SendRespUnregApp(clientDestAddr *net.UDPAddr) (err error) {
	resp := MakeRespUnregApp()
	return conn.SendMsg(RespMsgUnregApp, resp, clientDestAddr)
}

// SendRespUnregAppErr sends an error RespUnregApp msg back to a client
func (conn *ServerConn) SendRespUnregAppErr(error error, clientDestAddr *net.UDPAddr) (err error) {
	resp := MakeRespUnregAppErr(error)
	return conn.SendMsg(RespMsgUnregApp, resp, clientDestAddr)
}

// SendRespQryApp sends a RespQryApp msg back to a client
func (conn *ServerConn) SendRespQryApp(foundAppAddr *net.UDPAddr, clientDestAddr *net.UDPAddr) (err error) {
	resp := MakeRespQryApp(foundAppAddr)
	return conn.SendMsg(RespMsgQryApp, resp, clientDestAddr)
}

// SendRespQryAppErr sends an error RespQryApp msg back to a client
func (conn *ServerConn) SendRespQryAppErr(error error, clientDestAddr *net.UDPAddr) (err error) {
	resp := MakeRespQryAppErr(error)
	return conn.SendMsg(RespMsgQryApp, resp, clientDestAddr)
}
