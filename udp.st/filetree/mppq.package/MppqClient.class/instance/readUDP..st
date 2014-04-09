as yet unclassified
readUDP: aCollection
	"reads the UDP response(s) (as string) into the collection
	this should be forked, since it runs forever, blocking on Read	"

	| buffer resultArray udpPacket nbBytes |
	buffer := String new: bufferSize.
	
	[	"read from UDP socket"
		resultArray := clientConn receiveUDPDataInto: buffer.
		"debug"
		self halt.

		"get results from returned array"
		nbBytes := resultArray at: 1.
		nbBytes > 0 ifTrue: [ 
			"create a new udp packet and add it to the collection"
			udpPacket := UDPPacket new.
			udpPacket data: (buffer copyFrom: 1 to: nbBytes).
			udpPacket remoteAddr: (resultArray at: 2) asSocketAddress.
			
			aCollection add: udpPacket.
		].
	] repeat.
