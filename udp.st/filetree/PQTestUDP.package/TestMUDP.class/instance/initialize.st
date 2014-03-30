initialization
initialize
	super initialize.
	broadcastHost := NetNameResolver addressFromString: '239.255.13.4'.
	broadcastPort := 1440.
	