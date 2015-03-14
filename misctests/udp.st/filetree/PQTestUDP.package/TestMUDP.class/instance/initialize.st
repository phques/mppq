initialization
initialize
	super initialize.
	broadcastHost := NetNameResolver addressFromString: '224.0.1.60'.
	broadcastPort := 1888.
	broadcastHost := NetNameResolver addressFromString: '239.255.0.13'.
	broadcastHost := #[239 255 0 13] asSocketAddress .
	broadcastPort := 1440.
