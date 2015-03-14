converting
toJSON
	"convert to JSON string"
	^ String
		streamContents: [ :stream | 
			(NeoJSONWriter on: stream)
				mapInstVarsFor: (self params at: 1) class;
				nextPut: self asDictionary ]