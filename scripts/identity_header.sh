#!/bin/bash

ENC=`echo "{\"identity\": {\"type\": \"User\", \"account_number\":\"$1\", \"internal\":{\"org_id\":\"$2\"}}}" | base64 -w0`
echo "x-rh-identity: $ENC"
