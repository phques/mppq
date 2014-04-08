as yet unclassified
getResponses
	"reads the response(s) after sending the query"

	| buffer |
	buffer := String new: bufferSize.
	
	"receiveUDPDataInto: is blocking !
	 how can we handle waiting for responses for only a specific time ?
	-> start a thread ? How to stop it ?...
	"