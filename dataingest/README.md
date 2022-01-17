# Data Ingest Processes

WARNING these processes are in rapid flux as of Jan 2022 and subject to change.

Various scripts to provision & load data into a AWS RDS (Postgres 13.4)
instance used by the Find Insights back-end team.

Most are dependent on the existance of Postgres client utilities being
installed and also a configured aws command line client

## Environment Setup

Environment should contain similar to:

```
PGUSER=insights
PGPASSWORD=insights
PGHOST=localhost
PGPORT=5432
PGDATABASE=census
```

Confirm with 

```
$ env |grep PG
```

# Provision dev DB locally from scratch (source data) & create prod DB dump

TODO automate further. Rewrite in Go.

Assumes postgres is running on developer desktop natively or via docker with pg
client tools (eg. psql, pgdump, createdb, dropdb etc.) in the PATH

## Create database and schema

```
$ cd dp-find-insights-poc-api && make update-schema
```

Answer "y" to questions in last step. "Record not found warning" can be
ignored.

## Download source data 

```
$ git clone git@github.com:ONSdigital/nomis-bulk-to-postgres.git
$ cd nomis-bulk-to-postgres && ./download-data-2i.sh
```

Note the last step downloads CSV under the "data" directory. 

"download-data-2i.sh" downloads some QS and a few KS rows

"download-data-qs.sh" would be used for a large dev db (eg. for DataVis) with
all QS rows and no KS rows

## Populate database from downloaded source data

Assumes python venv is setup and "secrets.json" in place as described in [nomis-bulk-to-postgres README](https://github.com/ONSdigital/nomis-bulk-to-postgres/blob/main/README.md)

TODO python should pick up from the env like everything else (?) or just rewrite
TODO migrate this to go code enable to run under AWS (Lambda or Fargate etc.) &
read from S3

* The following imports data

```
$ . ./bulk/bin/activate
$ python add_to_db.py
```

## Clean up Database

*  Run 
```
$ cd dp-find-insights-poc-api && make update-schema
$ ./dataingest/dbsetup/cleandb.sh
 ```

## Non LSOA geo.name population

* unzip "Code_History_Database_(June_2021)_UK.zip"
* see dataingest/geo/README.md & fi-census-data/geo

```
$ cd dataingest/geo && go run .
```

## Geo data import

* Import data from GeoJSON files
  * This puts POLYGON in geo.wkb_geometry
  * *.geojson from https://github.com/ONSdigital/fi-census-data/tree/main/spatial should be in place

```
$ cd dp-find-insights-poc-api/dataingest/spatial
$ ./import.sh linux-localhost

TODO remove "K04000001 name TODO" etc with null lat and long

```
* Populate geo.wkb_long_lat_geom with long, lat POINT

```
$ cd dataingest/spatial/longlatgeom    
$ go run .
```

## Sanity checks

Various database sanity checks to ensure steps aren't missed.

```
$ cd dp-find-insights-poc-api/dataingest
$ make test
```

Finis!

# More information

## Getting the live Postgres password

* `../../secrets/PGPASSWORD.env.asc`
  * Encrypted version of postgres password - currently same for "postgres" (admin
user) & "insights" (app user)

Decrypt within `secrets/` directory
```
gpg -d PGPASSWORD.env.asc
```

to create "PGPASSWORD.env" which can be sourced to populate PGPASSWORD

## producing dump of new prod type database

* "pg_dump census > census.sql" can be used to create a dump

* "creatdb.sh" is then run on that dump to import into a Postgres instance

## Initial Setup
* awscreate.sh
  * Tactical solution to create AWS RDS instance and security group opening (non-standard) postgres port of 54322.
  * XXX should be migrated to Terraform.
