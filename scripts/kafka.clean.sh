#!/bin/bash
BASEDIR=$(dirname $0)
source $BASEDIR/kafka.conf
[[ -f $BASEDIR/kafka.local.conf ]] && source $BASEDIR/kafka.local.conf

if [[ $CLEAN_CHECKOUTS -eq 1 ]]; then
  rm -rf $BASEDIR/kafka/
fi

rm -rf /tmp/kraft-combined-logs/
