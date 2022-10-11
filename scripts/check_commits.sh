#!/bin/bash

MESSAGE=$(git log --pretty=format:%s HEAD^..HEAD)

MAX_LENGTH=70
TYPES="chore demo docs feat fix refactor revert style test"
PATTERN="^([a-z]+)\([a-z\-\*]+\)\:\ (.*)$"

if [[ ${#MESSAGE} > $MAX_LENGTH ]]; then
     echo "ERROR: Commit message was ${#MESSAGE} characters long, but should be at most $MAX_LENGTH characters"
     exit 1
fi

if [[ "$MESSAGE" =~ ^WIP ]]; then
    exit 0
fi

if ! [[ "$MESSAGE" =~ $PATTERN ]]; then
    echo "ERROR: Commit message did not match 'type(scope): subject': $MESSAGE"
    exit 1
fi

TYPE=${BASH_REMATCH[1]}
if ! [[ $TYPES =~ (^| )$TYPE($| ) ]]; then
    echo "ERROR: Commit message's type '$TYPE' must be one of '$TYPES'"
    exit 1
fi
