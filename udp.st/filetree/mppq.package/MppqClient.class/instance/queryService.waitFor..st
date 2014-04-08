as yet unclassified
queryService: serviceName waitFor: waitSecs
	"does a query for service, returns found serviceDefs"

	| query  |

	"create query"
	query := self createQuery: serviceName .
	
	"send query"
	self connect.
	self sendQuery: query.
	
	"get response(s)"