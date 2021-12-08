#!/bin/bash -x

# gpg -d create.env.asc for PGPASSWORD
. ./create.env

PGPASSWORD="$PGPASSWORD" psql -h "$PGHOST" -U postgres -p "$PGPORT" -c "CREATE USER insights WITH PASSWORD '$PGPASSWORD'"
PGPASSWORD="$PGPASSWORD" psql -h "$PGHOST" -U postgres -p "$PGPORT" -c "ALTER USER insights WITH CREATEDB"
