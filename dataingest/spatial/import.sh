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

psql -c "ALTER TABLE lad_gis ADD CONSTRAINT uq_lad17cd UNIQUE(lad17cd)"
psql -c "ALTER TABLE lsoa_gis ADD CONSTRAINT uq_lsoa11cd UNIQUE(lsoa11cd)"

# copy LAD data into geo
psql <<EOT
\x
UPDATE geo SET wkb_geometry=lad_gis.wkb_geometry, long=lad_gis.long, lat=lad_gis.lat, name=lad_gis.lad17nm 
FROM lad_gis  
WHERE geo.code=lad_gis.lad17cd
EOT

# copy LSOA data into geo
psql <<EOT2
\x
UPDATE geo SET wkb_geometry=lsoa_gis.wkb_geometry, long=lsoa_gis.long, lat=lsoa_gis.lat, name=lsoa_gis.lsoa11nm 
FROM lsoa_gis  
WHERE geo.code=lsoa_gis.lsoa11cd
EOT2

psql -c "DROP TABLE lsoa_gis"
psql -c "DROP TABLE lad_gis"

# These were valid in 2011 (where our data comes from) but aren't anymore
# XXX probably incomplete list
psql <<EOT3
\x
UPDATE geo SET valid=false WHERE code IN ('E06000048','E07000097','E07000100','E07000101','E07000104','E08000020')
EOT3
