#!/bin/bash

export GOPATH=/home/philippe/code/mppq:$GOPATH
gobind -lang=go -outdir=gen code.google.com/p/mppq/examples/libprovider/provider

