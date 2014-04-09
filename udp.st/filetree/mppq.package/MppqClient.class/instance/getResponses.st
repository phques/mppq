as yet unclassified
getResponses
	"reads the response(s) after sending the query"

	| buffer process responses |
	buffer := String new: bufferSize.
	
	"launch a thread to read udp"
	responses := OrderedCollection new.
	
	"debug"
	self readUDP: responses.
	
	process := [ self readUDP: responses ] fork.
	
	"let it run X seconds"
	responseWaitTime wait.
	process terminate.
	
	"process reponses"
	self processResponses:  responses.