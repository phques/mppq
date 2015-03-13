// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package marcopolo

import (
	"bytes"
	"fmt"
	"reflect"
)

// Client connection
type ClientConn struct {
	Conn

	// recvMsgs
	//RespChan        chan ClientMsgObj
	msgs            map[string]*registeredMsg // messages to support on recv
	packetChan      chan UdpPacket            // recvd udp packets channel
	registerMsgChan chan registerMsg          // channel to register registeredMsg's
}

// ClientMsgObj is holds the msg + param obj from a recv message in recvMsgs
type ClientMsgObj struct {
	Msg      string      // the msg, ie "marco.polo.resp.regapp"
	MsgParam interface{} // the param object for the msg
}

// MsgHandler is the type for callback funcs called for a recvd client msg (cf registeredMsg)
type MsgHandler func(*ClientConn, *ClientMsgObj) bool

// registeredMsg holds info on a message that we want to process
type registeredMsg struct {
	msg          string       // the message to expect
	msgPrefix    []byte       // includes the separator '|'
	paramObjType reflect.Type // TypeOf the msg parameter object

	msgHandler MsgHandler
	resultChan chan bool // msgHandler sends on this when completed
}

// actions for registerMsg
const (
	regMsgAddMsg = iota
	regMsgRemMsg
	regMsgAddHandler
	regMsgRemHandler
)

// registerMsg is used to send a msg to register on ClentConn.registerMsgChan
type registerMsg struct {
	registeredMsg registeredMsg
	action        int // (regMsgAddMsg) used when sent to recvMsgs, indicates add/remv
}

//--------

// Open opens an UDP Conn on a sys allocated local address&port
// it also initializes secondary members
func (conn *ClientConn) Open() (err error) {
	//conn.RespChan = make(chan ClientMsgObj)
	conn.msgs = make(map[string]*registeredMsg)
	conn.packetChan = make(chan UdpPacket)
	conn.registerMsgChan = make(chan registerMsg)

	// open udp connection
	err = conn.open(0)
	if err != nil {
		return
	}

	// start goroutine that reads messages
	go conn.recvMsgs()

	// register the standard response msgs / types (recvMsgs must be started)
	conn.RegisterMsg(RespMsgRegApp, Resp{})
	conn.RegisterMsg(RespMsgUnregApp, Resp{})
	conn.RegisterMsg(RespMsgQryApp, RespQryApp{})

	return
}

// Close closes the connection
// and cleans-up secondary members
func (conn *ClientConn) Close() {
	close(conn.packetChan)
	close(conn.registerMsgChan)
	//close(conn.RespChan)
	//conn.RespChan = nil
	conn.msgs = nil
	conn.packetChan = nil
	conn.registerMsgChan = nil
	conn.Conn.close()
}

//-----

// RegisterMsg registers a message with it's corresponding param type
func (conn *ClientConn) RegisterMsg(msg string, msgObj interface{}) {
	// create a registeredMsg{} to register this msg w. it's param object type
	regMsg := registeredMsg{
		msg:          msg,
		msgPrefix:    append([]byte(msg), cmdMsgSep),
		paramObjType: reflect.TypeOf(msgObj),
		msgHandler:   nil,
		resultChan:   nil,
	}

	// register through a channel send !
	conn.registerMsgChan <- registerMsg{regMsg, regMsgAddMsg}
}

// RegisterHandler registers a msg handler, returns the handler's result channel
func (conn *ClientConn) RegisterHandler(msg string, handler MsgHandler) chan bool {

	// create a registeredMsg{}
	regMsg := registeredMsg{}
	regMsg.msg = msg
	regMsg.msgHandler = handler
	regMsg.resultChan = make(chan bool)

	// register through a channel send !
	conn.registerMsgChan <- registerMsg{regMsg, regMsgAddHandler}

	return regMsg.resultChan
}

// UnregisterHandler unregisters a msg handler (handler's channel is closed)
func (conn *ClientConn) UnregisterHandler(msg string) {

	// create a registeredMsg{}
	regMsg := registeredMsg{}
	regMsg.msg = msg

	// unregister through a channel send !
	conn.registerMsgChan <- registerMsg{regMsg, regMsgRemHandler}
}

//-----

// recvMsgs (goroutine) recvs and processes msgs,
// when a valid/registered msg is parsed, if it has a handler, then it is called
func (conn *ClientConn) recvMsgs() {

	// start goroutine that reads udp packets & sends them on channel
	go conn.readUdpPacket()

	// stop when channels are closed
	var chanReadOk = true
	for chanReadOk {
		var udpPacket UdpPacket
		var registerMsg registerMsg

		select {
		// recv & process udp packets
		case udpPacket, chanReadOk = <-conn.packetChan:
			if chanReadOk {
				conn.processMsgPacket(udpPacket)
			}

		// we have a registeredMsg to process (add/rem the registeredMsg entry)
		case registerMsg, chanReadOk = <-conn.registerMsgChan:
			if chanReadOk {
				conn.doRegisterMsg(registerMsg)
			}
		}
	}
}

func (conn *ClientConn) doRegisterMsg(registerMsg registerMsg) {
	//fmt.Printf("doRegisterMsg, %+v\n", registerMsg)

	regMsg := registerMsg.registeredMsg

	switch registerMsg.action {
	case regMsgAddMsg:
		// add new msg to process into conn.msgs
		conn.msgs[regMsg.msg] = &regMsg

	case regMsgRemMsg:
		// remove msg to be processed from conn.msgs
		delete(conn.msgs, regMsg.msg)

	case regMsgAddHandler:
		if regdMsg, ok := conn.msgs[regMsg.msg]; ok {
			regdMsg.msgHandler = regMsg.msgHandler
			regdMsg.resultChan = regMsg.resultChan
		} else {
			//##!! getting complicated to return error to caller here !
			fmt.Printf("ClientConn.doRegisterMsg::regMsgAddHandler, unregistered msg %q\n", regMsg.msg)
		}

	case regMsgRemHandler:
		if regdMsg, ok := conn.msgs[regMsg.msg]; ok {
			close(regdMsg.resultChan)
			regdMsg.msgHandler = nil
			regdMsg.resultChan = nil
		} else {
			//##!! getting complicated to return error to caller here !
			fmt.Printf("ClientConn.doRegisterMsg::regMsgRemHandler, unregistered msg %q\n", regMsg.msg)
		}
	}
}

// readUdpPacket (goroutine) continually recvs udp packets
// and then sends them on packetChan
func (conn *ClientConn) readUdpPacket() {
	for {
		// get msg as strings, (blocking) recv udpPacket
		if udpPacket, err := conn.recvMsgStr(); err == nil {
			// send udpPacket on channel
			conn.packetChan <- udpPacket
		} else {
			// if channels are closed, don't display read error
			var chanOk bool
			select {
			case _, chanOk = <-conn.packetChan:
			default:
			}
			// read error w/o closed channels/conn
			if chanOk {
				fmt.Println("Error receiving client msg :", err)
			}
			// exit when error
			break
		}
	}

	fmt.Println("ClientConn.readUdpPacket stop")
}

// processMsgPacket verfies / parses a msg udp packet
func (conn *ClientConn) processMsgPacket(udpPacket UdpPacket) {

	// check for valid & registered msg, splits into cmd & param strings
	var registeredMsg *registeredMsg
	var msgName string
	var msgParam []byte
	foundMsg := false
	for _, registeredMsg = range conn.msgs {

		if bytes.HasPrefix(udpPacket.data, registeredMsg.msgPrefix) {
			// ok valid cmd msg prefix, split in 2: msg & JSON string param
			// "marco.polo.xxx|JSON-param"
			prefixLen := len(registeredMsg.msgPrefix)
			msgName = string(udpPacket.data[:prefixLen-1]) // stripped separator
			msgParam = udpPacket.data[prefixLen:]
			foundMsg = true
			break
		}
	}

	if !foundMsg {
		len := len(udpPacket.data)
		if len > 32 {
			len = 32
		}
		fmt.Printf("Error, ClientConn.recvMsgs, unexpected msg %q...\n", udpPacket.data[:len])
		return
	}

	// create new msg param object according to msg type
	msgParamObj := reflect.New(registeredMsg.paramObjType).Interface()

	//## debug
	//fmt.Printf("recv msg: %s | %s\n", msgName, msgParam)

	// parse / populate msg param object with JSON
	if err := parseJsonMsgParam(msgParam, msgParamObj); err != nil {
		fmt.Println("Error parsing JSON msg param :", err)
		return
	}

	// check version from param
	//if msgParamObj appNameParam.Version != version {
	//	err = fmt.Errorf("got cmd version '%v', expecting '%v'", appNameParam.Version, version)
	//	return
	//}

	// if we have a func to handle this, call it
	if funcInfo, ok := conn.msgs[msgName]; ok {
		if funcInfo.msgHandler == nil {
			fmt.Println("no handler for msg ", msgName)
			return
		}

		// Create a ClientMsgObj
		clientMsgObj := ClientMsgObj{
			Msg:      msgName,
			MsgParam: msgParamObj,
		}

		// call func
		funcRes := funcInfo.msgHandler(conn, &clientMsgObj)

		// send func ret val on result channel (non blocking)
		select {
		case funcInfo.resultChan <- funcRes:
		default:
		}
	} else {
		//## debug only ?
		fmt.Println("recvd unhandled msg: ", msgName)
	}

}

//------------

// SendRegAppName sends a marcopolo msg to register an app
func (conn *ClientConn) SendRegAppName(appNameStr string) (err error) {
	return conn.sendAppNameMsg(CmdRegApp, appNameStr)
}

// SendQryAppName sends a marcopolo msg to query for an app
func (conn *ClientConn) SendQryAppName(appNameStr string) (err error) {
	return conn.sendAppNameMsg(CmdQryApp, appNameStr)
}

// SendUnregAppName sends a marcopolo msg to unregister an app
func (conn *ClientConn) SendUnregAppName(appNameStr string) (err error) {
	return conn.sendAppNameMsg(CmdUnregApp, appNameStr)
}

//---------

// sendAppNameMsg sends a marco.polo cmd msg with appNameStr as parameter
func (conn *ClientConn) sendAppNameMsg(cmdMsg, appNameStr string) (err error) {
	// setup an AppNameParam
	var appName = AppNameParam{StdHeader{version}, appNameStr}
	// send the message
	return conn.SendMsgToMarcoPolo(cmdMsg, appName)
}
