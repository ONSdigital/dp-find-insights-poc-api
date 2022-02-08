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
POSTGRES_PASSWORD=superuser-passwd
```

```
PGHOST_INTERNAL=host.docker.internal
PGPORT_INTERNAL=5432
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

## Download NOMIS source data 

XXX actual 2021 data will probably look different and need further processing

```
$ git clone git@github.com:ONSdigital/nomis-bulk-to-postgres.git
$ cd nomis-bulk-to-postgres && ./download-data-2i.sh
```

The last step downloads CSV under the "data" directory. 

"download-data-2i.sh" downloads some QS and a few KS rows

"download-data-qs.sh" would be used for a large dev db (eg. for DataVis) with
all QS rows and no KS rows

## Populate database from downloaded source data

Assumes python venv is setup and "secrets.json" in place as described in [nomis-bulk-to-postgres README](https://github.com/ONSdigital/nomis-bulk-to-postgres/blob/main/README.md)

TODO python should pick up from the env like everything else (?) or just rewrite
TODO migrate this to go code enable to run under AWS (Lambda or Fargate etc.) &
read from S3

* The following imports data and takes ~20mins

```
$ . ./bulk/bin/activate
$ python add_to_db.py
```

TODO skipping import of MSOA would be faster

## Clean up Database

*  Run 
```
$ cd dp-find-insights-poc-api && make update-schema
$ ./dataingest/dbsetup/cleandb.sh
 ```

## geo.name population

* see geo/README.md

```
$ cd dataingest/geo
$ aws --region eu-central-1 s3 sync s3://find-insights-input-data-files/geoname/ .
$ go run .
```

## Geo data import

* Import data from GeoJSON files
  * This puts POLYGON in geo.wkb_geometry
  * *.geojson from https://github.com/ONSdigital/fi-census-data/tree/main/spatial should be in place

```
$ cd ../spatial
$ aws --region eu-central-1 s3 sync s3://find-insights-input-data-files/geojson/ .
```

```
go build ./geo2sql.go
```

On linux:
```
$ ./import.sh
```

On mac:
```
$ /usr/local/bin/bash import.sh
```
(`/bin/bash` on mac is too old, so use `brew install bash` or equivalent to get a bash 4+.)

* Populate geo.wkb_long_lat_geom with long, lat POINT

```
$ cd longlatgeom    
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
