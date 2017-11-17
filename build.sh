#!/usr/bin/env bash
# TODO We can deploy this from a dockerfile, as a binary, once we have everything we want in it.

export CANARY=$GOPATH/src/github.com/blackducksoftware/canary/

if [ ! -d $CANARY ]; then
	echo "Exiting the build.  Looks like your gopath isnt set up to have $CANARY !"
	exit 1
fi 	

set -x

rm main

# This will put the 'sidecar' binary into your GOPATH.
go build ./cmd/sidecar/service_scanner.go
$GOPATH/bin/sidecar
