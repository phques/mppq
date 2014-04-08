initialization
initialize
	super initialize.
	port := 1440.
	mcaddr := #[239 255 0 13] asSocketAddress . "multicast address"
	bcaddr := #[255 255 255 255] asSocketAddress . "broadcast address"
	clientConn := nil.
	useBroadcast := false.
	
	whosthereStr := 'mppq.whosthere?'. "query sent to find services"
	imhereStr := 'mppq.imHere!'.    "response with serviceDef"
