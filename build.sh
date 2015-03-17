#!/bin/sh

export GOPATH=`pwd`

if [ ! -d src/github.com/comstud/slopher ] ; then
    echo "Fetching slopher..."
    go get github.com/comstud/slopher
    echo "Done"
fi

go build src/bot.go
