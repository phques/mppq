#!/bin/bash

export GOPATH=/home/philippe/code/mppq:$GOPATH
gobind -lang=java -outdir=gen code.google.com/p/mppq/examples/libprovider/provider 


