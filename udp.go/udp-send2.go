// main.go
package main

import (
	"fmt"
	"net"
)

func revMask(ip net.IP, mask net.IPMask) net.IP {
	ip = ip.To4()
	if ip != nil {
		for i, _ := range ip {
			ip[i] = (ip[i] & mask[i]) | (^mask[i])
		}
	}
	return ip
}

func makeBroadcast(addrStr string) net.IP {
	ip := net.ParseIP(addrStr)
	if ip == nil {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	mask := ip.DefaultMask()
	if mask == nil {
		return nil
	}
	return nil
}

func main() {

	// open local UDP port

	// destination broadcast UDP address
	remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 4444}
	//remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(192, 168, 1, 255), Port: 4444}
	//remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(192, 168, 1, 149), Port: 4444}

	udpConn, err := net.DialUDP("udp4", nil, &remoteBroadcastUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()

	fmt.Printf("%v", udpConn.LocalAddr())

	// send 'msg's
	udpConn.Write([]byte("this is msg 1"))
	udpConn.Write([]byte("this is msg 2"))
	udpConn.Write([]byte("this is msg 3"))
	udpConn.Write([]byte("quit"))
}
