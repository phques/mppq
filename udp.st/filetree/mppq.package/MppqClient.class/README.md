MppqCLient, client side of marcoPoloPQ.
Providers (servers) wait for service queries and respond with a serviceDefinition.

Query sent on mujlticast UDP, response sent back unicast udp to querier.
On Win8/8.1, multicast is not working properly, so use broadcast instead.