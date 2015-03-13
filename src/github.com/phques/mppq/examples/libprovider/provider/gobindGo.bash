#!/bin/bash

export GOPATH=/home/philippe/code/mppq:$GOPATH
gobind -lang=go -outdir=gen github.com/phques/mppq/examples/libprovider/provider

