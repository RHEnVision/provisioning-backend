#!/bin/bash
#
# See README.sources
#
source sources.conf
rm -rf ./sources-api-go ./sources-database-populator/
sudo su - postgres -c "dropdb $DATABASE_NAME"
sudo su - postgres -c "dropuser $DATABASE_USER"
