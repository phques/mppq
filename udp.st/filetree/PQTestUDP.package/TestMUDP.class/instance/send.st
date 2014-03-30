others
send
	clientConn sendUDPData: 'mess 1' toHost: broadcastHost port: broadcastPort.
	clientConn sendUDPData: 'mess 2' toHost: broadcastHost port: broadcastPort.
	clientConn sendUDPData: 'quit' toHost: broadcastHost port: broadcastPort.
