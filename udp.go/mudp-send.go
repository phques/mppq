// main.go
package main

import (
	"fmt"
	"log"
	"net"
)

func send(udpConn *net.UDPConn, remoteAddr *net.UDPAddr, msg string) {
	fmt.Println("sending msg")

	nbBytes, err := udpConn.WriteToUDP([]byte(msg), remoteAddr)
	if err != nil {
		fmt.Println("error sending marco ", err)
		return
	}

	fmt.Printf("wrote %d bytes\n", nbBytes)
    
    // recv response
    data := make([]byte, 1024)
    fmt.Println("wait for response")
    nbRead, remoteAddr, err := udpConn.ReadFromUDP(data)
    if err != nil {
        fmt.Println("error reading udp socket ", err)
        return
    }
    data = data[:nbRead]

    fmt.Printf("read %d bytes '%s' response from udp\n", nbRead, data)
}

func main() {

	// destination broadcast UDP address
	//mcaddr := &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 4444}
	//mcaddr := &net.UDPAddr{IP: net.IPv4(192, 168, 56, 255), Port: 4444}
	//mcaddr := &net.UDPAddr{IP: net.IPv4(192, 168, 1, 255), Port: 4444}
	//mcaddr := &net.UDPAddr{IP: net.IPv4(192, 168, 1, 149), Port: 4444}

	//mcaddr, err := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	//mcaddr, err := net.ResolveUDPAddr("udp4", "224.0.1.60:1888")
	mcaddr, err := net.ResolveUDPAddr("udp4", "239.255.0.13:1440")
	//mcaddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:1440")
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
	//~ udpConn, err := net.ListenMulticastUDP("udp4", interf, mcaddr)

	if err != nil {
		log.Fatal("error ListenUDP: ", err)
	}
	defer udpConn.Close()

	fmt.Printf("%v", udpConn.LocalAddr())

	// send 'msg's
	send(udpConn, mcaddr, "this is msg 1")
	send(udpConn, mcaddr, "this is msg 2")
	send(udpConn, mcaddr, "mppq.whosthere?androidPush")
	send(udpConn, mcaddr, "quit")
}
