private
callArith: arithReq
	"call the jsonrpc Arith server, show the result json string"

	| socket resp |
	[ 
	"connect to server"
	socket := Socket newTCP.
	socket connectTo: NetNameResolver localHostAddress port: 8222.
	socket waitForConnectionFor: 1.	"send request & read response"
	socket sendData: arithReq toJSON.
	resp := socket receiveData ]
		ensure: [ socket closeAndDestroy ].	"show response"
	resp ifNotNil: [ self showResponse: resp FromReq: arithReq ]