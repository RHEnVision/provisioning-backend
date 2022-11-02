#!/bin/bash
BASEDIR=$(dirname $0)
source $BASEDIR/kafka.conf
[[ -f $BASEDIR/kafka.local.conf ]] && source $BASEDIR/kafka.local.conf

test -d $BASEDIR/kafka || mkdir $BASEDIR/kafka

pushd $BASEDIR/kafka
  if [[ -f kafka.tgz ]]; then
    echo "Kafka already downloaded"
  else
    echo "Downloading Kafka..."
    curl -L -o kafka.tgz "https://www.apache.org/dist/kafka/$VERSION/kafka_$SVERSION-$VERSION.tgz"

    echo "Extracting Kafka..."
    tar -xzf kafka.tgz --strip 1
  fi
popd

echo "All done, run ./kafka.start.sh"
