#!/bin/bash

set -e

if [[ $1 == linux-localhost ]]; then
    EXTRA=--network="host"
fi

DOCKER="osgeo/gdal:alpine-small-3.3.3"

declare -A tables
tables["lsoa_gis"]="Lower_Layer_Super_Output_Areas_(December_2011)_Boundaries_Super_Generalised_Clipped_(BSC)_EW_V3.geojson"
tables["lad_gis"]="Local_Authority_Districts_(December_2017)_Boundaries_in_the_UK_(WGS84).geojson"

for TABLE in "${!tables[@]}"; do
    GEOJSON="${tables[$TABLE]}"
    echo "creating '$TABLE' in '$PGDATABASE' on '$PGHOST'"
    docker run $EXTRA -v $PWD:$PWD $DOCKER ogr2ogr -f "PostgreSQL" PG:"host=$PGHOST user=$PGUSER dbname=$PGDATABASE password=$PGPASSWORD port=$PGPORT" "$PWD/$GEOJSON" -nln "$TABLE" --config PG_USE_COPY YES -lco GEOM_TYPE=geometry
    psql -c "VACUUM ANALYZE $TABLE"
done
