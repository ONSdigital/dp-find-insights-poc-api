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
PGPASSWORD=$POSTGRES_PASSWORD dropdb --username postgres --if-exists "$PGDATABASE"
yes | make update-schema
go run ./dataingest/addtodb
(yes | make update-schema) 
cd dataingest/geoname && go run .  
cd ../spatial && ./lad2011ish.sh && go build ./geo2sql.go && ./import.sh
cd longlatgeom  && go run .    
cd ../../postcode  && go run . 
delta=$((SECONDS-otime))
echo "about" $((delta/60)) "min(s) elapsed"
cd ../../dataingest && make test
