initialization
initialize
	super initialize.
	"broadcastHost := NetNameResolver addressFromString: '255.255.255.255'."
	broadcastHost := #[255 255 255 255] asSocketAddress .
	broadcastPort := 1440.
	