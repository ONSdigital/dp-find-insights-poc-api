#! /usr/bin/env bash

set -e

tables='
lsoa_gis|Lower_Layer_Super_Output_Areas_(December_2011)_Boundaries_Super_Generalised_Clipped_(BSC)_EW_V3.geojson
lad_gis|LAD2011ish.geojson
msoa_gis|Middle_Layer_Super_Output_Areas_(December_2011)_Boundaries_Super_Generalised_Clipped_(BSC)_EW_V3.geojson
'

while read line
do
    if test -z "$line"
    then
        continue
    fi
    oIFS=$IFS
    IFS='|'
    set -- $line
    IFS=$oIFS

    TABLE=$1
    GEOJSON=$2

    echo "creating '$TABLE' in '$PGDATABASE' on '$PGHOST'"
    ./geo2sql -t "$TABLE" -f "$GEOJSON" | psql -f -
    psql -c "VACUUM ANALYZE $TABLE"
done <<-EOF
$tables
EOF

psql -c "ALTER TABLE lad_gis ADD CONSTRAINT uq_lad17cd UNIQUE(lad17cd)"
psql -c "ALTER TABLE lsoa_gis ADD CONSTRAINT uq_lsoa11cd UNIQUE(lsoa11cd)"
psql -c "ALTER TABLE msoa_gis ADD CONSTRAINT uq_msoa11cd UNIQUE(msoa11cd)"

# copy LAD data into geo
psql <<EOT
\x
UPDATE geo SET wkb_geometry=lad_gis.wkb_geometry, long=lad_gis.long, lat=lad_gis.lat, name=lad_gis.lad17nm, welsh_name=lad_gis.lad17nmw
FROM lad_gis  
WHERE geo.code=lad_gis.lad17cd AND geo.type_id=4
EOT

# copy LSOA data into geo
psql <<EOT2
\x
UPDATE geo SET wkb_geometry=lsoa_gis.wkb_geometry, long=lsoa_gis.long, lat=lsoa_gis.lat, name=lsoa_gis.lsoa11nm,  welsh_name=lsoa11nmw 
FROM lsoa_gis  
WHERE geo.code=lsoa_gis.lsoa11cd AND geo.type_id=6
EOT2

# copy MSOA data into geo
# DON'T SET name we use another source (House of Commons library) of better ones for MSOA!
# but populate welsh_name
psql <<EOT3
\x
UPDATE geo SET wkb_geometry=msoa_gis.wkb_geometry, long=msoa_gis.long, lat=msoa_gis.lat, welsh_name=msoa11nmw
FROM msoa_gis  
WHERE geo.code=msoa_gis.msoa11cd AND geo.type_id=5
EOT3

# clean up temp tables
psql -c "DROP TABLE lsoa_gis"
psql -c "DROP TABLE lad_gis"
psql -c "DROP TABLE msoa_gis"

# set "English Welsh" names (which aren't actually Welsh) to be HoC MSOA names
psql -c "update geo set welsh_name=name where code like 'E%' and type_id=5;"

psql -c "VACUUM ANALYZE geo"
