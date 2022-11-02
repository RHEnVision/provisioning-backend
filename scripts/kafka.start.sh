#!/bin/bash
BASEDIR=$(dirname $0)
source $BASEDIR/kafka.conf
[[ -f $BASEDIR/kafka.local.conf ]] && source $BASEDIR/kafka.local.conf

echo "Formatting journal..."
$BASEDIR/kafka/bin/kafka-storage.sh format -t $UUID -c $BASEDIR/kafka/config/kraft/server.properties

echo "Starting Kafka..."
$BASEDIR/kafka/bin/kafka-server-start.sh $BASEDIR/kafka/config/kraft/server.properties
