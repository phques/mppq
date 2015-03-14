as yet unclassified
readResp
	| buff resp |
	"clientConn setPeer: broadcastHost port: broadcastPort."
	buff := String new: 4096.
	3 timesRepeat:  [ 
		resp := clientConn receiveUDPDataInto: buff.
		Transcript show: resp; cr.
		200 milliSeconds asDelay wait.
	].
	