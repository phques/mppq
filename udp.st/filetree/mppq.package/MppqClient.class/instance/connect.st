as yet unclassified
connect
	"create the udp socket"

	clientConn := Socket newUDP.
	self useBroadcast ifTrue: [
		clientConn setOption: 'SO_BROADCAST' value: true. "win8 multicast not working local to local !"
	].
