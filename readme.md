// marcoPoloPQ project
// Copyright 2014 Philippe Quesnel
// Licensed under the Academic Free License version 3.0


repos to hold work in progress, 
yet another, very basic, UDP 'marco polo' and JSONRPC tests

Goal is to use multicast UDP,
even though currently (March 2014), my Windows platform (Win8.1) has probs with this !?
Probably will need to support broadcast UDP for this

Also considered 'borrowing' techniques from SLP .. ?


### Using addresses:
RFC 2365 - Administratively Scoped IP Multicast

The IPv4 Local Scope -- 239.255.0.0/16
(??The IPv4 Organization Local Scope -- 239.192.0.0/14)

Picked port 1440 (not listed in wikipedia well know ports)
(-> eicon-SLP ?)


