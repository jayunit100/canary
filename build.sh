#!/usr/bin/env bash
# TODO We can deploy this from a dockerfile, as a binary, once we have everything we want in it.

export GOPATH=/Users/`whoami`/go

echo "assuming that you've cloned this to $GOPATH/src/bitbucket.org/bdsengineering/hub-sidecar"

if [ ! -f $GOPATH ]; then
	GOPATH=/Users/`whoami`/work/go
	echo "tried a different path : $GOPATH"
fi
if [ ! -d $GOPATH/src/bitbucket.org/bdsengineering/hub-sidecar/ ]; then
	echo "Exiting the build.  Looks like your gopath isnt set up!"
	exit 1
fi 	

set -x
rm main 

go build cmd/main.go && ./main


#go build ./src/github.com/jayunit100/hub-sidecar/main/main.go
