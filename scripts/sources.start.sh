#!/bin/bash
#
# See README.sources
#
source sources.conf
pushd sources-api-go
go build
./sources-api-go
popd
