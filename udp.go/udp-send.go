// main.go
package main

import (
	"fmt"
	"net"
)

func send(udpConn *net.UDPConn, remoteAddr net.UDPAddr, msg string) {
	fmt.Println("sending msg")

	nbBytes, err := udpConn.WriteToUDP([]byte(msg), &remoteAddr)
	if err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	fmt.Printf("wrote %d bytes\n", nbBytes)
}

func main() {

	// open local UDP port
	fmt.Println("net.DialUDP")

	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	//localUdpAddr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}

	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()

	fmt.Printf("%v", udpConn.LocalAddr())

	// destination broadcast UDP address
	//remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 4444}
	//remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(192, 168, 56, 255), Port: 4444}
	remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(192, 168, 1, 255), Port: 4444}
	//remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(192, 168, 1, 149), Port: 4444}

	// send 'msg's
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 1")
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 2")
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 3")
	send(udpConn, remoteBroadcastUdpAddr, "quit")
}
