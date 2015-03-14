private
getResponses: duration
	"reads the response(s) after sending the query"

	| packets  |

	"read packets responses"	
	packets := self readUDP: duration.

	"process reponses"
	^ self processResponses:  packets.