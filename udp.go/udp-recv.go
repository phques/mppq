package main

import (
	"fmt"
	"net"
	//"time"
)

func main() {
	// --- listen for udp msg ---
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 4444,
	})
	if err != nil {
		fmt.Println("error open udp socket ", err)
		return
	}
	defer udpConn.Close()

	// wait for msg
	for {
		data := make([]byte, 1024)
		fmt.Println("wait for msg")
		nbRead, remoteAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error reading udp socket ", err)
			return
		}
		data = data[:nbRead]

		fmt.Printf("read %d bytes '%s' from udp (remote %s)\n", nbRead, data, remoteAddr.String())
		if string(data) == "quit" {
			break
		}

		//fmt.Println("wait a sec...")
		//time.Sleep(1 * time.Second)
	}
}
