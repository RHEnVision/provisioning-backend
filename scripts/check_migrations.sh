#!/bin/sh

CHANGED_COMMITED=$(git diff --name-status --oneline HEAD..HEAD^ | grep '^M' | grep 'internal/db/migrations/.*sql')
CHANGED=$(git diff --name-status --oneline | grep '^M' | grep 'internal/db/migrations/.*sql')

if [[ -n "$CHANGED$CHANGED_COMMITED" ]]; then
  echo "Last commit, staged or unstaged changes do appear to modify migration(s)."
  echo "This is not allowed, create a new migration instead either manually or:"
  echo
  echo "  make generate-migration MIGRATION_NAME=add_new_column"
  echo
  echo "If you still want to pass CI tests, create an additional commit to get this"
  echo "check passing as it only checks the latest one. Use with care!"
  echo
  echo "The following files were modified:"
  echo
  echo $CHANGED
  echo $CHANGED_COMMITED
  exit 1
fi
