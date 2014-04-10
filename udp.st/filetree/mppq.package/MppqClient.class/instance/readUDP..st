private
readUDP: duration
	"reads the UDP response(s) (as string) into a collection"

	| buffer resultArray udpPacket nbBytes packets |
	buffer := String new: bufferSize.
	packets := OrderedCollection new.
	
	(duration asSeconds * 5) timesRepeat: [
		"read from UDP socket, non-blocking"
		resultArray := clientConn receiveUDPDataInto: buffer.
		"debug"
		Transcript show: resultArray; cr.

		"get results from returned array"
		nbBytes := resultArray at: 1.
		nbBytes > 0 ifTrue: [ 
			"create a new udp packet and add it to the collection"
			udpPacket := UDPPacket new.
			udpPacket data: (buffer copyFrom: 1 to: nbBytes).
			udpPacket remoteAddr: (resultArray at: 2).
			
			packets add: udpPacket.
		].
		200 milliSeconds wait.
	].

	^ packets.