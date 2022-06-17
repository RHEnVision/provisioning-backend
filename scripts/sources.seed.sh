#!/bin/bash
#
# See README.sources
#
source sources.conf
git clone https://github.com/MikelAlejoBR/sources-database-populator
pushd sources-database-populator
go run main.go
popd
