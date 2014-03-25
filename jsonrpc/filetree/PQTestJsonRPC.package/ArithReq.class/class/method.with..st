private
method: methodName with: arithArg
	| req |
	req := self new.
	req method: 'Arith.' , methodName.
	req params: { arithArg }. "golang jsonrpc is ver 1, needs to be array"
	^ req.
	