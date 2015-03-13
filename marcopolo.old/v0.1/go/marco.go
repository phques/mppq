// main.go
package main

import (
	"fmt"
	"net"
	"time"
)

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func main() {

	// open local UDP port
	fmt.Println("net.DialUDP")

	localUdpAddr := net.UDPAddr{ IP: net.IPv4(0, 0, 0, 0), Port: 0 }
	
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}


	// destination 'marco' broadcast UDP address
    remoteBroadcastUdpAddr := net.UDPAddr{ IP: net.IPv4(255, 255, 255, 255), Port: 4444 }

	answer := make([]byte, 1024)

	var success bool = false
	for !success {
		// send 'marco'
		fmt.Println("sending marco")
		
		marcoMsg := "marco|testMarcoPolo"
		nbBytes, err := udpConn.WriteToUDP([]byte(marcoMsg), &remoteBroadcastUdpAddr)
		if err != nil {
			fmt.Println("error sending marco ", err)
			return
		}
		_ = nbBytes
		//fmt.Printf("wrote %d bytes\n", nbBytes)

		// read back answer, try 10x 100ms (total 1s)
		fmt.Println("get polo answer")

		for i := 0; !success && i < 10; i++ {
			// set read timeout on connection to 100ms
			deadline := time.Now().Add(time.Millisecond * 100)
			udpConn.SetReadDeadline(deadline)

			// Try to read answer
			//nbBytes, udpAddr, err := udpConnRead.ReadFromUDP(answer)
			//nbBytes, addr, err := udpConnRead.ReadFrom(answer)
			nbBytes, err := udpConn.Read(answer)
			if err != nil {
				if !isTimeout(err) {
					fmt.Println("error reading back answer ", err)
					return
				}
				//fmt.Println("timeout reading back answer")
			} else {
				//fmt.Printf("read %d bytes: '%s', from %s\n", nbBytes, answer[:nbBytes], udpAddr.String())
				fmt.Printf("read %d bytes: '%s'\n", nbBytes, answer[:nbBytes])
				success = true
			}
		}
	}
}
