// main.go
package main

import (
	"fmt"
	"net"
    "log"
)

func send(udpConn *net.UDPConn, remoteAddr *net.UDPAddr, msg string) {
	fmt.Println("sending msg")

	nbBytes, err := udpConn.WriteToUDP([]byte(msg), remoteAddr)
	if err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	fmt.Printf("wrote %d bytes\n", nbBytes)
}

func main() {

	// destination broadcast UDP address
	//remoteBroadcastUdpAddr := &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 4444}
	//remoteBroadcastUdpAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 56, 255), Port: 4444}
	//remoteBroadcastUdpAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 1, 255), Port: 4444}
	//remoteBroadcastUdpAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 1, 149), Port: 4444}

//~     remoteBroadcastUdpAddr, err := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
    remoteBroadcastUdpAddr, err := net.ResolveUDPAddr("udp4", "224.0.1.60:1888")
	if err != nil {
		log.Fatal("error ResolveUDPAddrt ", err)
	}

	// open local UDP port
	fmt.Println("net.DialUDP")

	//localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	//localUdpAddr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
    localUdpAddr, err := net.ResolveUDPAddr("udp4", ":0")
	if err != nil {
		log.Fatal("error ResolveUDPAddrt ", err)
	}
    
	udpConn, err := net.ListenUDP("udp4", localUdpAddr)

    //~ ethname := "{5BF6D791-D59A-40A0-BDD0-FADD0A065A8E}"
    //~ interf, err := net.InterfaceByName(ethname)
    //~ fmt.Printf("%v, %v\n", err, interf)

    //~ udpConn, err := net.ListenMulticastUDP("udp4", nil, mcaddr)
    //~ udpConn, err := net.ListenMulticastUDP("udp4", interf, remoteBroadcastUdpAddr)

	if err != nil {
		log.Fatal("error ListenUDP: ", err)
	}
	defer udpConn.Close()

	fmt.Printf("%v", udpConn.LocalAddr())

	// send 'msg's
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 1")
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 2")
	send(udpConn, remoteBroadcastUdpAddr, "this is msg 3")
	send(udpConn, remoteBroadcastUdpAddr, "quit")
}
