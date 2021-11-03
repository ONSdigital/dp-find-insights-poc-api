#!/bin/bash -x

# gpg -d create.env.asc for PGPASSWORD
. ./create.env

HOST="fi-database-1.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"
PORT=54322

if [ "$1" == "-create-user" ]; then
    PGPASSWORD="$PGPASSWORD" psql -h "$HOST" -U postgres -p "$PORT" -c "CREATE USER insights WITH PASSWORD '$PGPASSWORD'"
    PGPASSWORD="$PGPASSWORD" psql -h "$HOST" -U postgres -p "$PORT" -c "ALTER USER insights WITH CREATEDB"
fi

# connect to db which exists!
PGPASSWORD="$PGPASSWORD" psql -h "$HOST" -U insights -p "$PORT" -d postgres -f insights.sql
