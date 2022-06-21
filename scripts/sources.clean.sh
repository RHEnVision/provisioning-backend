#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
rm -rf $BASEDIR/sources-api-go/ $BASEDIR/sources-database-populator/
sudo su - postgres -c "dropdb $DATABASE_NAME"
sudo su - postgres -c "dropuser $DATABASE_USER"
