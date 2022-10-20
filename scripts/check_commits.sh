#!/bin/bash

MESSAGE=$(git log --pretty=format:%s HEAD^..HEAD)
if [[ "$MESSAGE" =~ ^WIP ]]; then
  exit 0
fi

MAX_LENGTH=70
TYPES="build ci docs feat fix perf refactor test chore"
PATTERN="^[a-z\(\)]+\:\ (.*)$"
SCOPED_PATTERN="^([a-z]+)\([a-z\-\*]+\)\:\ (.*)$"

if [[ ${#MESSAGE} > $MAX_LENGTH ]]; then
  echo "ERROR: Commit message length too long (70): ${#MESSAGE}"
  exit 1
fi

if ! [[ "$MESSAGE" =~ $PATTERN ]]; then
  echo "ERROR: Commit message did not match 'type: subject': $MESSAGE"
  exit 1
fi

if [[ "$MESSAGE" =~ $SCOPED_PATTERN ]]; then
  TYPE=${BASH_REMATCH[1]}
  if ! [[ $TYPES =~ (^| )$TYPE($| ) ]]; then
    echo "ERROR: Commit message's type '$TYPE' must be one of '$TYPES'"
    exit 1
  fi
fi
