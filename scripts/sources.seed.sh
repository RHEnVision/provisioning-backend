#!/bin/bash
#
# See README.sources
#
BASEDIR=$(dirname $0)
source $BASEDIR/sources.conf
[[ -f $BASEDIR/sources.local.conf ]] && source $BASEDIR/sources.local.conf

if [[ -z "$ARN_ROLE" ]]; then
	echo "ARN_ROLE must be defined in sources.local.conf!"
	exit 1
fi

echo "Creating $ARN_ROLE with account_id $ACCOUNT_ID org_id $ORG_ID"
IDENTITY=$($BASEDIR/identity_header.sh $ACCOUNT_ID $ORG_ID)
# create source for provisioning type with account 13/00013
curl --location -g  --request POST "http://localhost:$PORT/api/sources/v3.1/bulk_create" \
--header "$IDENTITY" \
-d "$(cat <<EOF
{
	"sources": [
		{
			"name": "Amazon source",
			"source_type_name": "amazon",
			"app_creation_workflow": "manual_configuration"
		}
	],
	"applications": [
		{
			"source_name": "Amazon source",
			"application_type_name": "provisioning"
		}
	],
	"authentications": [
		{
			"resource_type": "Application",
			"resource_name": "provisioning",
			"username": "$ARN_ROLE",
			"authtype":"provisioning-arn"
		}
	]
}
EOF
)"

# the following is only useful when you want a lot of data in the db
#git clone https://github.com/MikelAlejoBR/sources-database-populator $BASEDIR/sources-database-populator
#pushd $BASEDIR/sources-database-populator
#go run main.go
#popd
