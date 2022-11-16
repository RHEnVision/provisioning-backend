#!/bin/bash
BASEDIR=$(dirname $0)
source $BASEDIR/kafka.conf
[[ -f $BASEDIR/kafka.local.conf ]] && source $BASEDIR/kafka.local.conf

echo "Formatting journal..."
$BASEDIR/kafka/bin/kafka-storage.sh format -t $UUID -c $BASEDIR/kafka/config/kraft/server.properties

echo "Starting Kafka..."
$BASEDIR/kafka/bin/kafka-server-start.sh $BASEDIR/kafka/config/kraft/server.properties

echo "Creating topics..."
$BASEDIR/kafka/bin/kafka-topics.sh --create --topic platform.provisioning.internal.availability-check --bootstrap-server $KAFKA_HOST
$BASEDIR/kafka/bin/kafka-topics.sh --create --topic platform.sources.event-stream --bootstrap-server $KAFKA_HOST
