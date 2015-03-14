private
sendQuery: query
	| addr |
	addr := (self useBroadcast) ifTrue: [ bcaddr  ] ifFalse: [ mcaddr ].
	clientConn sendUDPData: query toHost: addr port: port.
