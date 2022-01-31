#! /bin/sh

exec docker run \
    -it \
    --env PGHOST="${PGHOST_INTERNAL:-$PGHOST}" \
    --env PGPORT="${PGPORT_INTERNAL:-$PGPORT}" \
    --env PGDATABASE \
    --env PGUSER \
    --env PGPASSWORD \
    --env POSTGRES_PASSWORD \
    -v "$PWD:/tmp" \
    --workdir /tmp \
    --rm \
    postgis/postgis \
    psql "$@"
