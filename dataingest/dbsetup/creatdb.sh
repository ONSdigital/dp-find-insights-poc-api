#!/bin/bash
# Depends on the presence of PG_ env in the env

set -e

if [[ $1 == "" ]]; then
    echo "must pass db dump as arg"
    exit 1
fi

# generate SQL without exposing password on a command line
echo "If $PGUSER already exists, the next line will fail; that's OK"
cat <<EOF | PGPASSWORD="$POSTGRES_PASSWORD" psql -U postgres -d postgres -f -
CREATE USER $PGUSER WITH PASSWORD '$PGPASSWORD' CREATEDB
EOF

createdb "$PGDATABASE"
PGPASSWORD="$POSTGRES_PASSWORD" psql -U postgres -d "$PGDATABASE" -c "CREATE EXTENSION postgis"

psql <<EOT
\x
SET synchronous_commit TO off;
\i $1
EOT

cd ../.. && make update-schema
