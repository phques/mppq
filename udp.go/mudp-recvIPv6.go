package main

import (
	"fmt"
	"log"
	"net"
	//"time"
)

func main() {
	// --- listen for udp msg ---
	mcaddr, err := net.ResolveUDPAddr("udp6", "[FF01::1]:1440")
	if err != nil {
		log.Fatal("error ResolveUDPAddr ", err)
	}

	conn, err := net.ListenMulticastUDP("udp6", nil, mcaddr)
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
		if string(udpData) == "quit" {
			break
		}

		//fmt.Println("wait a sec...")
		//time.Sleep(1 * time.Second)
	}
}
