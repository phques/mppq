private
processResponses:  packets
	"TODO"
	Transcript show: packets size; cr.
	
	packets do: [ :packet |
		Transcript show: packet printString; cr.
	].

	^ nil.