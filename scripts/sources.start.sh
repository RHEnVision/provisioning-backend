#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
[[ -f $BASEDIR/sources.local.conf ]] && source $BASEDIR/sources.local.conf

pushd $BASEDIR/sources-api-go
go build
./sources-api-go
popd
