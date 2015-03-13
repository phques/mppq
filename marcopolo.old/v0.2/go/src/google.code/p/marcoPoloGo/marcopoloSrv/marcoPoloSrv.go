// MarcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// marcoPoloSrv.go
package main

import (
	"fmt"
	"google.code/p/marcoPoloGo/marcopolo"
	//"os"
	//"os/signal"
	//"syscall"

	"net"
)

//-----

// app holds the data for a registered app
type app struct {
	name    string
	appAddr *net.UDPAddr
}

// apps contains all the registered apps
var apps = make(map[string]app)

//-----------

//TODO: log
// regApp registers an app
func regApp(cmd marcopolo.CmdMsg, srvConn *marcopolo.ServerConn) {
	appName := cmd.AppNameParam.AppName

	fmt.Println("register app ", appName)

	// lookup any prev regd app w. this name
	prevApp, found := apps[appName]
	if found {
		fmt.Printf("regapp, removing prev old app @ %s\n", prevApp.appAddr)

		// already an app regd w. this name
		// remove old entry,
		delete(apps, appName)

		//TODO: ? send notif to prev app 'marco.polo.kickedout'
		_ = prevApp
	}

	// register new app
	app := app{name: appName, appAddr: cmd.UdpPacket.RemoteAddr}
	apps[appName] = app

	// send OK to caller
	srvConn.SendRespRegApp(cmd.UdpPacket.RemoteAddr)
}

//TODO: log
// unregApp unregisters an app
func unregApp(cmd marcopolo.CmdMsg, srvConn *marcopolo.ServerConn) {
	appName := cmd.AppNameParam.AppName
	fmt.Println("unregister app ", appName)

	// lookup regd app w. this name
	prevApp, found := apps[appName]
	_ = prevApp
	if found {
		fmt.Printf("unregapp, removing @ %s\n", prevApp.appAddr)

		// remove app entry
		delete(apps, appName)

		// send OK to caller
		srvConn.SendRespUnregApp(cmd.UdpPacket.RemoteAddr)
	} else {
		err := fmt.Errorf("unregapp, not found : '%s'", appName)
		fmt.Println(err)

		// send error back to app
		srvConn.SendRespUnregAppErr(err, cmd.UdpPacket.RemoteAddr)
	}
}

//TODO: log
// qryApp looks up an app and sends back the answer to the caller
func qryApp(cmd marcopolo.CmdMsg, srvConn *marcopolo.ServerConn) {
	appName := cmd.AppNameParam.AppName
	fmt.Println("qry app ", appName)

	// look for a regd app w. this name
	foundApp, found := apps[appName]
	if found {
		//TODO: ping app to confirm it is still there (&unreg if no answer)!
		//  or do we leave this to caller ?
		//  the problem I see w. pign here is that since we use UDP,
		//  the caller might do a sendQry/tryToGetResp loop .. as we would do here !
		//  Which would mean we would get multiple/dup queries from same caller,
		//  which we would need to know how to hanle (recognize & skip,
		//  or keep list of callers for an app name & send all answers back: not as clean but simpler)
		fmt.Printf("qryapp, found @ %s\n", foundApp.appAddr)

		// send back answer to caller
		srvConn.SendRespQryApp(foundApp.appAddr, cmd.UdpPacket.RemoteAddr)
	} else {
		// send back error to caller
		fmt.Println("qryapp, not found")

		errMsg := fmt.Errorf("qryapp, '%s' not found ", appName)
		srvConn.SendRespQryAppErr(errMsg, cmd.UdpPacket.RemoteAddr)
	}
}

//TODO: register for sigint/sigterm & close cleanly

func main() {

	// open local UDP marcolo server connection
	fmt.Println("open udp connection")

	var srvConn = new(marcopolo.ServerConn)
	err := srvConn.Open()
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer srvConn.Close()

	for {
		//TODO: ?? periodic pings to detect disconnected regd apps
		//TODO: once we start having multi goroutines, handle ownership of app objects !

		// wait for marco.polo cmdMsg
		cmdMsg, err := srvConn.RecvCmdMsg()
		if err != nil {
			// show error
			fmt.Printf("error recving marco.polo cmdMsg: %q\n", err)
			//return
		} else {
			// show results
			fmt.Println("---recvd:--")
			fmt.Println(" cmd:", cmdMsg.Cmd)
			fmt.Println(" appname:", cmdMsg.AppNameParam.AppName)
			fmt.Println(" ver:", cmdMsg.AppNameParam.Version)
			fmt.Println(" from:", cmdMsg.UdpPacket.RemoteAddr)

			// process the command
			switch cmdMsg.Cmd {
			case marcopolo.CmdRegApp:
				regApp(cmdMsg, srvConn)
			case marcopolo.CmdUnregApp:
				unregApp(cmdMsg, srvConn)
			case marcopolo.CmdQryApp:
				qryApp(cmdMsg, srvConn)
			}
		}
	}
}
