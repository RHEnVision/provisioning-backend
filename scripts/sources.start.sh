#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
pushd $BASEDIR/sources-api-go
go build
./sources-api-go
popd
