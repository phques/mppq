as yet unclassified
processResponses:  responses
	"TODO"
	Transcript show: responses size; cr.
	
	responses do: [ :response|
		Transcript show: response printString; cr.
	]