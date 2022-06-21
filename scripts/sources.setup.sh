#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
git clone https://github.com/RedHatInsights/sources-api-go $BASEDIR/sources-api-go
sudo dnf install redis postgresql postgresql-server
export PGSETUP_INITDB_OPTIONS="--auth=trust"
sudo postgresql-setup --initdb --unit postgresql
sudo systemctl enable --now redis postgresql
sudo su - postgres -c "createuser $DATABASE_USER"
sudo su - postgres -c "createdb $DATABASE_NAME --owner $DATABASE_USER"

echo "Local users are configured as trusted in postgresql. If this is not"
echo "what you want, review pg_hba.conf file."
echo "All done, run ./sources.start.sh"
