// main.go
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	// open local UDP port
	fmt.Println("net.DialUDP")

	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()

	// destination 'marco' broadcast UDP address
	// will cause WriteToUDP error
	//~     remoteBroadcastUdpAddr := net.UDPAddr{ IP: net.IPv4(0,0,0,0), Port:0}
	remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 4444}

	// read
	readDoneChan := make(chan bool, 1)
	go func() {
		answer := make([]byte, 2000)
		nbBytes, err := udpConn.Read(answer)
		if err != nil {
			fmt.Println("error reading back answer: ", err)
		} else {
			fmt.Printf("recvd %d bytes: '%s'\n", nbBytes, answer[:nbBytes])
		}
		// test debug
		time.Sleep(time.Second)
		readDoneChan <- true
	}()

	var done bool = false
	for !done {
		// send 'marco'
		marcoMsg := "marco|testMarcoPolo"
		fmt.Println("sending ", marcoMsg)

		_, err = udpConn.WriteToUDP([]byte(marcoMsg), &remoteBroadcastUdpAddr)
		if err != nil {
			fmt.Println("error sending marco: ", err)

			// stop & wait for goroutine
			udpConn.Close() // will cause blocking Read above stop w. err
			<-readDoneChan
			done = true
		} else {
			select {
			// read is completed
			case <-readDoneChan:
				done = true

			// wait a bit before next send
			case <-time.After(time.Millisecond * 500):
			}
		}
	}

}
