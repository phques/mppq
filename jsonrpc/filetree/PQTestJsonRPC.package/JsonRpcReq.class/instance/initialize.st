initialization
initialize
	super initialize.
	"jsonrpc := '2.0'."	"note: golang jsonrpc is NOT 2.0 !"
	id := self class nextId.