#!/bin/bash -x

# gpg -d create.env.asc for PGPASSWORD
. ./create.env

# TODO this shares postgres (admin) account password with role user.  Should be
# different in PROD.

PGPASSWORD="$PGPASSWORD" psql -h "$PGHOST" -U postgres -p "$PGPORT" -c "CREATE USER insights WITH PASSWORD '$PGPASSWORD'"
PGPASSWORD="$PGPASSWORD" psql -h "$PGHOST" -U postgres -p "$PGPORT" -c "ALTER USER insights WITH CREATEDB"
