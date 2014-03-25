examples
testArith
	"test jsonrpc call Arith.Multiply & Divide"

	| arg arithReq response |
	
	"create ArithArg parameter for Arith.xx"
	arg := ArithArg new.
	arg a: 14.
	arg b: 3.

	"call Arith.Multiply : args"
	arithReq := 	ArithReq multiply: arg.
	response := self callArith: arithReq.

	"call Arith.Divide : args"
	arithReq := 	ArithReq divide: arg.
	response := self callArith: arithReq.
