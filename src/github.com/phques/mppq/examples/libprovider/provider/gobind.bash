#!/bin/bash

if [[ -z $1 ]]; then
    echo 'missing parameter: go | java'
    exit 1
fi

export GOPATH=/home/philippe/code/mppq:$GOPATH
gobind -lang=$1 -outdir=gen github.com/phques/mppq/examples/libprovider/provider

