#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
[[ -f $BASEDIR/sources.local.conf ]] && source $BASEDIR/sources.local.conf

if [[ $CLEAN_CHECKOUTS -eq 1 ]]; then
	rm -rf $BASEDIR/sources-api-go/ $BASEDIR/sources-database-populator/
fi
sudo su - postgres -c "dropdb $DATABASE_NAME"
sudo su - postgres -c "dropuser $DATABASE_USER"
