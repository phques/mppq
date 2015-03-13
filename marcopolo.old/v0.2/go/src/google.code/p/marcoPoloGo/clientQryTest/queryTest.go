// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// clientQryTest, queries marcopolo for an app address
package main

import (
	"fmt"
	"google.code/p/marcoPoloGo/marcopolo"
	"time"
)

type MyMsgParam struct {
	MyMsgName string
	MsgParam  int
}

const MyMsg = "kwez.org/androidPush.myMsg"

//----------

func sendMyAppMsg(clientConn *marcopolo.ClientConn, resp *marcopolo.RespQryApp) {
	destUdpAddr := resp.UDPAddr()
	param := MyMsgParam{"msg1", 22}
	clientConn.SendMsg(MyMsg, param, &destUdpAddr)

	fmt.Printf(" sent %q %v\n", MyMsg, param)
}

// handleQryAppResp is called to handle marcopolo.RespMsgQryApp
func handleQryAppResp(conn *marcopolo.ClientConn, msgObj *marcopolo.ClientMsgObj) bool {
	fmt.Println("handleQryAppResp")

	// cast msgObj.RespParam to expected marcopolo.Resp
	if param, convOk := msgObj.MsgParam.(*marcopolo.RespQryApp); convOk {
		// response says 'ok' ?
		if param.Ok {
			fmt.Printf(" ok, found app: %v\n", param.AppAddr)
			// send myMsg
			sendMyAppMsg(conn, param)
		} else {
			// oopsm recvd an error
			fmt.Println(" qryApp err: ", param.Error)
		}
		return true
	} else {
		fmt.Println("handleQryAppResp, failed to convert msgObj.RespQryApp")
		return false
	}
}

func main() {
	// open local UDP client connection
	fmt.Println("open client connection")

	var clientConn marcopolo.ClientConn
	if err := clientConn.Open(); err != nil {
		fmt.Println("error open client connection : ", err)
		return
	}
	defer clientConn.Close()

	// register QryApp handler
	resChan := clientConn.RegisterHandler(marcopolo.RespMsgQryApp, handleQryAppResp)

	// qry marcopolo server for app
	fmt.Println("query for app")

	if err := clientConn.SendQryAppName("kwez.org/androidPush"); err != nil {
		fmt.Println("error in SendQryAppName() ", err)
		return
	}

	//
	fmt.Println("wait for response")

	select {
	case <-resChan:
	case <-time.After(time.Second * 2):
		fmt.Println("no response for qryApp after 2secs")
	}

	clientConn.UnregisterHandler(marcopolo.RespMsgQryApp)
}
