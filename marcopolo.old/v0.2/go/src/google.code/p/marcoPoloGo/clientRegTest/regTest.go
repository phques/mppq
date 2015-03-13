// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// marcoPoloClientTest
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

const myMsg = "kwez.org/androidPush.myMsg"

var handlers = make(map[string]chan bool)

//----------

func regHandler(conn *marcopolo.ClientConn, msg string, handler marcopolo.MsgHandler) {
	zechan := conn.RegisterHandler(msg, handler)
	handlers[msg] = zechan
}

func regFuncs(conn *marcopolo.ClientConn) {
	regHandler(conn, marcopolo.RespMsgRegApp, handleRegAppResp)
	regHandler(conn, marcopolo.RespMsgUnregApp, handleRegAppResp) // same object as RespMsgRegApp

	regHandler(conn, myMsg, handleMyMsg)
}

// waitResp waits for an handled response by msg name
func waitResp(msg string) {
	// find our handler's chanell by msg name
	handlerChan, found := handlers[msg]
	if found {
		// wait for msgHandler to complete
		select {
		case <-handlerChan:
		case <-time.After(time.Second):
			fmt.Println("waitResp, no response after 1sec")
			return
		}
	} else {
		fmt.Printf("waitResp, no func for %q\n", msg)
	}
}

// handleRegAppResp is called to handle marcopolo.RespMsgRegApp / marcopolo.RespMsgRegApp
func handleRegAppResp(conn *marcopolo.ClientConn, msgObj *marcopolo.ClientMsgObj) bool {
	fmt.Println("handle ", msgObj.Msg)

	// cast msgObj.RespParam to expected marcopolo.Resp
	if param, convOk := msgObj.MsgParam.(*marcopolo.Resp); convOk {
		// response says 'ok' ?
		if param.Ok {
			fmt.Println(" resp ok")
		} else {
			// oopsm recvd an error
			fmt.Println(" err: ", param.Error)
		}
		return true
	} else {
		fmt.Println("handleRegAppResp, failed to convert msgObj.RespParam")
		return false
	}
}

//
func handleMyMsg(conn *marcopolo.ClientConn, msgObj *marcopolo.ClientMsgObj) bool {
	fmt.Println("handleMyMsg")

	// cast msgObj.RespParam to expected
	if param, convOk := msgObj.MsgParam.(*MyMsgParam); convOk {
		fmt.Printf(" msgName:%s, msgParam:%d\n", param.MyMsgName, param.MsgParam)
		return true
	} else {
		fmt.Println("handleMyMsg, failed to convert msgObj.RespParam")
		return false
	}
}

func main() {

	// open local UDP client connection
	fmt.Println("---open client connection")
	var clientConn marcopolo.ClientConn
	if err := clientConn.Open(); err != nil {
		fmt.Println("error open client connection : ", err)
		return
	}
	defer clientConn.Close()
	fmt.Printf("connected as %s\n", clientConn.Conn.String())

	// register our clientConn msg & handlers
	clientConn.RegisterMsg(myMsg, MyMsgParam{})
	regFuncs(&clientConn)

	// register app with marcopolo server
	fmt.Println("---register app")

	if err := clientConn.SendRegAppName("kwez.org/androidPush"); err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	// wait for response from RecvResps
	waitResp(marcopolo.RespMsgRegApp)

	// unregister incorrect app name
	fmt.Println("---unregister incorrect app")

	if err := clientConn.SendUnregAppName("kwez.org/androidPouche"); err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	// wait for response from RecvResps (expected error)
	waitResp(marcopolo.RespMsgUnregApp)

	// wait a bit
	fmt.Println("---sleep 10secs")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		fmt.Printf("\r%d", i)
	}

	// unregister app with marcopolo server
	fmt.Println("---unregister app")

	if err := clientConn.SendUnregAppName("kwez.org/androidPush"); err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	waitResp(marcopolo.RespMsgUnregApp)

}
