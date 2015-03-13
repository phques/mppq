#!/bin/sh

set -eux

go install -v google.code/p/marcoPoloGo..."$@"

#~ go install -v google.code/p/marcoPoloGo/marcopolo"$@"
#~ go build -v google.code/p/marcoPoloGo/clientQryTest "$@"
#~ go build -v google.code/p/marcoPoloGo/clientRegTest "$@"
#~ go build -v google.code/p/marcoPoloGo/marcopoloSrv "$@"
