initialization
initialize
	super initialize.
	broadcastHost := NetNameResolver addressFromString: '224.0.1.60'.
	broadcastPort := 1888.
	