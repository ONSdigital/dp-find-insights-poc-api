#!/bin/bash
# We don't actually need these env vars but it's more explicit to use them
# directly

set -e

if [[ $1 == "" ]]; then
    echo "must pass db dump as arg"
    exit 1
fi

createdb "$PGDATABASE"
psql -U postgres -d "$PGDATABASE" -c "CREATE EXTENSION postgis"

psql <<EOT
\x
SET synchronous_commit TO off;
\i $1
EOT

cd .. && make update-schema
