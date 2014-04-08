as yet unclassified
queryService: serviceName waitFor: waitSecs
	"does a query for service, returns found serviceDefs"

	| query responses |

	"create query"
	query := self createQuery: serviceName .
	
	"send query"
	self connect.
	self sendQuery: query.
	
	"get response(s)"
	responses := self getResponses.
	