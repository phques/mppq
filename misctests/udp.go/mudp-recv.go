package main

import (
	"fmt"
	"log"
	"net"
	//"time"
)

func main() {
	// --- listen for udp msg ---
	mcaddr, err := net.ResolveUDPAddr("udp4", "239.255.0.13:1440")
	//mcaddr, err := net.ResolveUDPAddr("udp4", "239.255.43.99:1888")
	//mcaddr, err := net.ResolveUDPAddr("udp4", "224.0.1.60:1888")
	if err != nil {
		log.Fatal("error ResolveUDPAddr ", err)
	}

	//~ ethname := "{5BF6D791-D59A-40A0-BDD0-FADD0A065A8E}"
	//~ interf, err := net.InterfaceByName(ethname)
	//~ fmt.Printf("%v, %v\n", err, interf)

	//~ interf, err := net.InterfaceByIndex(4)
	//~ fmt.Printf("%v, %v\n", err, interf)
	//~ conn, err := net.ListenMulticastUDP("udp4", interf, mcaddr)

	conn, err := net.ListenMulticastUDP("udp4", nil, mcaddr)
	if err != nil {
		log.Fatal("error ListenMulticastUDP ", err)
	}
	defer conn.Close()

	// wait for msg
	data := make([]byte, 1024)
	for {
		fmt.Println("wait for msg")
		nbRead, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error reading udp socket ", err)
			return
		}
		udpData := data[:nbRead]

		fmt.Printf("read %d bytes '%s' from udp (remote %s)\n", nbRead, udpData, remoteAddr.String())
        
        // send back some response
        _, err = conn.WriteToUDP([]byte("some response"), remoteAddr)
        if err != nil {
            fmt.Println("error sending response ", err)
            return
        }

		if string(udpData) == "quit" {
			break
		}
	}
}
