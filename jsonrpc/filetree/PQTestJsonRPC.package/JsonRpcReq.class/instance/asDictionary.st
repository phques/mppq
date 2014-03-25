converting
asDictionary
	^ Dictionary new
		at: #id put: id;
		at: #method put: method;
		at: #params put: params;
		yourself