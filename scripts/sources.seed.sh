#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
git clone https://github.com/MikelAlejoBR/sources-database-populator $BASEDIR/sources-database-populator
pushd $BASEDIR/sources-database-populator
go run main.go
popd
