// main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

func f1() {
	fmt.Println("---f1---")
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	for _, a := range addrs {
		fmt.Println(a)
		ip := net.ParseIP(a)
		if ip != nil {
			if mask := ip.DefaultMask(); mask != nil {
				maskedIP := ip.Mask(mask)
				fmt.Println("  DefaultMask ", mask)
				fmt.Println("  masked  ", maskedIP)
				fmt.Println("  rev masked  ", revMask(ip, mask))
				//fmt.Println(" IsGlobalUnicast ", ip.IsGlobalUnicast())
				//fmt.Println(" IsLinkLocalMulticast ", ip.IsLinkLocalMulticast())
				//fmt.Println(" IsLinkLocalUnicast ", ip.IsLinkLocalUnicast())
				//fmt.Println(" IsLoopback ", ip.IsLoopback())
				//fmt.Println(" IsMulticast ", ip.IsMulticast())
			}
		}
	}

}

func f2() {
	fmt.Println("---f2---")
	interfs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, interf := range interfs {
		addrs, _ := interf.Addrs()
		fmt.Printf("%v %v %b %d\n", interf.Name, addrs, interf.Flags, interf.Index)
	}
}

// return  default local address IP
// ie- address that is used when connecting to a remove IP using 'default local IP' (no IP used)
func getDefaultLocalIP() net.IP {
	udpConn, err := net.DialUDP("udp4", nil, // no local IP
		&net.UDPAddr{
			//IP: net.IPv4(192, 168, 1, 1),
			IP:   net.IPv4(255, 255, 255, 255),
			Port: 0,
		})
	if err != nil {
		return nil
	}
	defer udpConn.Close()

	// unfortunatly, we can ONLY get the local address with locaAddr() ,
	// which is a string "ip:port", so we have to process the string to find the IP
	localAddr := udpConn.LocalAddr().String()
	fmt.Printf("udpconn local addr %v\n", localAddr)

	host, _, err := net.SplitHostPort(localAddr)
	return net.ParseIP(host)
}

func f3() {
	ip := getDefaultLocalIP()
	fmt.Printf("udpconn local addr IP %v\n", ip)
	ip = revMask(ip, ip.DefaultMask())
	fmt.Printf("broadcast %v\n", ip)
}

func f4() {
	// windows, hmm, superb name ;-p
	ethname := "{5BF6D791-D59A-40A0-BDD0-FADD0A065A8E}"
	interf, err := net.InterfaceByName(ethname)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", interf)
}

func main() {

	//f1()
	f2()
	//f3()

}
