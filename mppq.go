// MarcoPoloPQ project
// Copyright 2024 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//
// https://github.com/phques/mppq
// Simple service discovery using multicast UDP.
// 'marco' client query, 'polo' server response.
// Use MppqServer.Serve(ctx) to start the server.
// Use MppqQueryServiceDefs() to query the server.

package mppq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/phques/gomisc/ordone"
)

//-----------------

const (
	mppqQueryHeader    = "mppq?marco"
	mppqResponseHeader = "mppq!polo"
	udpPort            = 1440
)

var (
	// setup in init()
	multicastUdpAddr net.UDPAddr
)

func init() {
	//RFC 2365 - Administratively Scoped IP Multicast
	//  The IPv4 Local Scope -- 239.255.0.0/16
	multicastUdpAddr = net.UDPAddr{IP: net.IPv4(239, 255, 0, 13), Port: udpPort}
}

//-----------------

// ServiceDef describes a service the provider offers.
type ServiceDef struct {
	ServiceName  string // the name of the service provided
	ProviderName string // the name of the provider
	HostIP       string // IP address on which the service is available
	HostPort     int    // port on which the service is available
	Protocol     string // app arbitrary protocol name, ie "jsonrpcv1"
}

// Create a service definition
func NewServiceDef(serviceName, providerName string, hostIP string, hostPort int, protocol string) ServiceDef {
	return ServiceDef{
		ServiceName:  serviceName,
		ProviderName: providerName,
		HostIP:       hostIP,
		HostPort:     hostPort,
		Protocol:     protocol,
	}
}

// Create a service definition using local hostname and local IP
func NewLocalServiceDef(serviceName string, hostPort int, protocol string) ServiceDef {
	hostname, _ := os.Hostname() // get the local hostname
	hostIP := GetOutboundIP()    // get the local IP of this machine
	return NewServiceDef(serviceName, hostname, hostIP.String(), hostPort, protocol)
}

//-----------------

// MppqQuery defines the mppq query (cf NewStdMppqQuery)
type MppqQuery struct {
	Query string
}

// Create a MppqQuery, with the query's correct values
func NewStdMppqQuery() MppqQuery {
	return MppqQuery{
		Query: mppqQueryHeader,
	}
}

// marshal json into MppqQuery
func (q MppqQuery) toJSON() ([]byte, error) {
	return json.MarshalIndent(q, "", "  ")
}

// unmarshal MppqQuery from json
func (q *MppqQuery) fromJSON(data []byte) error {
	return json.Unmarshal(data, q)
}

//-----------------

// MppqResponse defines the mppq response
type MppqResponse struct {
	Header      string // holds the expected mppqResponseHeader
	ServiceDefs []ServiceDef
}

// Create a MppqResponse, with header set to correct value
func NewMppqResponse(serviceDefs []ServiceDef) MppqResponse {
	return MppqResponse{
		Header:      mppqResponseHeader,
		ServiceDefs: serviceDefs,
	}
}

// Validate MppqResponse, based on header / mppqResponseHeader
func IsValidMppqResponse(r MppqResponse) bool {
	return r.Header == mppqResponseHeader
}

// marshal json into MppqResponse
func (r MppqResponse) toJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// unmarshal MppqResponse from json
func (r *MppqResponse) fromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

// ------------
// holds data received from udp
type udpData struct {
	data []byte
	addr *net.UDPAddr
}

// create a new udpData, *making a copy* of the data
func newUDPData(data []byte, addr *net.UDPAddr) udpData {
	newData := make([]byte, len(data))
	copy(newData, data)

	return udpData{
		data: newData,
		addr: addr,
	}
}

// ------------

// MppqServer implements a mppq server
type MppqServer struct {
	serviceDefs []ServiceDef
}

// Create a new MppqServer
func NewMppqServer(serviceDefs []ServiceDef) *MppqServer {
	return &MppqServer{
		serviceDefs: serviceDefs,
	}
}

// Listen on udp and process / respond to received MppqQuery
func (s *MppqServer) Serve(ctx context.Context) {

	// create udp connection that listens on multicast udp
	udpConn, err := net.ListenMulticastUDP("udp4", nil, &multicastUdpAddr)
	if err != nil {
		log.Println("failed to open multicast udp listen connection. ", err)
		return
	}
	defer udpConn.Close()

	log.Println("listening on", udpConn.LocalAddr())

	// read from udp and send data to dataCh, in a goroutine loop until ctx.done()
	dataCh := readFromUDP(ctx, udpConn)

	// process / respond to received MppqQuery,
	// stop on ctx.done()
	for recvdData := range ordone.OrDone(ctx.Done(), dataCh) {
		s.processRecvdQuery(recvdData, udpConn)
	}

	//nb: udpConn is closed here, stopping readFromUDP

	log.Println("MppqServer stopped")
}

// process a received MppqQuery, unmarshalling the data and sending back a MppqResponse
func (s *MppqServer) processRecvdQuery(recvdData udpData, udpConn *net.UDPConn) {
	// unmarshal and validate MppqQuery
	var recvdMppQuery MppqQuery
	if err := recvdMppQuery.fromJSON(recvdData.data); err != nil {
		log.Println("failed to unmarshal MppqQuery from json. ", err)
		return
	}

	if recvdMppQuery != NewStdMppqQuery() {
		log.Println("UDP is not a valid MppqQuery query. ", recvdMppQuery)
		return
	}

	log.Printf("valid std MppqQuery received! %#v\n", recvdMppQuery)

	// create MppqResponse
	response := NewMppqResponse(s.serviceDefs)
	rData, err := response.toJSON()
	if err != nil {
		log.Println("failed to marshal MppqResponse to json. ", err)
		return
	}

	// write MppqResponse to udp
	log.Println("writing MppqResponse to udp")
	_, err = udpConn.WriteToUDP(rData, recvdData.addr)
	if err != nil {
		log.Println("failed to write MppqResponse to udp. ", err)
		return
	}
}

// read from udp and send data to dataCh, in an goroutine loop,
// stop and close dataCh if the udpConn or ctx.done() is closed
// note that udpConn.ReadFromUDP is blocking
func readFromUDP(ctx context.Context, udpConn *net.UDPConn) <-chan udpData {
	dataCh := make(chan udpData)

	go func() {
		defer close(dataCh)
		buf := make([]byte, 1024)
		for {
			// read from udp (this is blocking!!) (will stop if conn is closed)
			nbRead, addr, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				// Could not find a better way to check if the error is because the connection is closed
				if strings.Contains(err.Error(), "use of closed network connection") {
					log.Println("readFromUDP, stopping: closed network connection.")
				} else {
					log.Println("readFromUDP stopping, failed to read from udp connection:", err)
				}
				return
			}

			log.Printf("readFromUDP received %d bytes from %s\n", nbRead, addr)

			// send to dataCh
			select {
			case <-ctx.Done():
				return
			case dataCh <- newUDPData(buf[:nbRead], addr):
			}
		}
	}()

	return dataCh
}

// QueryServiceDefs sends a MppqQuery to get service definitions from all mppq servers
// nbTries: number of tries (1 second timeout each)
func QueryServiceDefs(nbTries int) (*MppqResponse, error) {
	// create udp connection
	var err error
	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		log.Println("failed to open udp connection. ", err)
		return nil, err
	}
	defer udpConn.Close()

	// create a mppqQuery to send
	mppqQuery := NewStdMppqQuery()
	queryData, err := mppqQuery.toJSON()
	if err != nil {
		log.Println("failed to marshal MppqQuery to json. ", err)
		return nil, err
	}

	// to read the response
	data := make([]byte, 16*1024)
	var mppqResponse MppqResponse

	// try nbTries times to get a valid MppqResponse
	for i := 0; i < nbTries; i++ {

		// send a MppqQuery to the multicast udp
		log.Println("sending MppqQuery to", multicastUdpAddr.String())
		_, err = udpConn.WriteToUDP(queryData, &multicastUdpAddr)
		if err != nil {
			log.Println("failed to write MppqQuery to udp. ", err)
			return nil, err
		}

		log.Printf(
			"(%d/%d) waiting (1sec) for MppqResponse from %s\n",
			i,
			nbTries,
			multicastUdpAddr.String(),
		)

		// read from udp, 1 sec timeout
		udpConn.SetReadDeadline(time.Now().Add(time.Second))
		nbRead, remoteAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			// don't fail on timeout
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}
			// failed to read from udp, stop
			log.Println("failed to read response from udp. ", err)
			return nil, err
		}

		// we received a response
		log.Printf("read %d bytes from %s\n", nbRead, remoteAddr)

		// unmarshal to MppqResponse
		if err = mppqResponse.fromJSON(data[:nbRead]); err != nil {
			log.Println("failed to unmarshal MppqResponse from json. ", err)
			continue
		}

		// validate mppqResponse
		if !IsValidMppqResponse(mppqResponse) {
			log.Println("invalid MppqResponse, header does not match. ", mppqResponse.Header)
			continue
		}

		// success!
		return &mppqResponse, nil
	}

	log.Println("failed to get valid MppqResponse after", nbTries, "tries")
	return nil, errors.New("failed to get valid MppqResponse")
}
