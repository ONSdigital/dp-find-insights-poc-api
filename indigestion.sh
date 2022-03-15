#!/usr/bin/env bash
# temporary script to help steve nom the data
# this will be replaced (probably by go)

echo -n "Drop and recreate '$PGDATABASE' (y/n)?"
read -r a
if [[ $a != "y" ]]; then
    exit 1
fi

otime=$SECONDS
set -e -x
dropdb $PGDATABASE
yes | make update-schema
go run ./dataingest/addtodb
(yes | make update-schema) 
cd dataingest/geoname && go run .  
cd ../spatial && go build ./geo2sql.go && ./import.sh
cd longlatgeom  && go run .    
echo "sec(s) elapsed: " $(($SECONDS-$otime))
cd ../../../dataingest && make test
