private
processResponses:  packets
	"TODO"
	packets do: [ :packet | | jsonStr |
		"debug"
		Transcript show: packet; cr.
		
		(packet data beginsWith: imhereStr) ifTrue: [ 
			jsonStr := packet data allButFirst: imhereStr size.
			"debug"
			Transcript cr; show: jsonStr; cr.
		].
	].

	
	^ nil.