# Data Ingest Processes

WARNING these processes are in rapid flux as of Dec 2021 and subject to change.

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
PGDATABASE=censustiny
```

Confirm with 

```
$ env |grep PG
```

## dbsetup

* `../../secrets/PGPASSWORD.env.asc`
  * Encrypted version of postgres password - currently same for "postgres" (admin
user) & "insights" (app user)

Decrypt within `secrets/` directory
```
gpg -d PGPASSWORD.env.asc
```

to create "PGPASSWORD.env"

## One off processes
* awscreate.sh
  * Tactical solution to create AWS RDS instance and security group opening (non-standard) postgres port of 54322.
  * XXX should be migrated to Terraform.

* creatdb.sh
  * used to import an existing db dump & enable PostGIS

# Provision Tiny dev DB locally & create prod DB dump

These instructions are for the "new" or "skinny" database.

TODO automate further

Assumes postgres is running on developer desktop natively or via docker with pg
client tools (eg. psql, pgdump, createdb, dropdb etc.) in the PATH


## Create database and schema

```
$ createdb censustiny
$ cd dp-find-insights-poc-api
$ make update-schema
```

Answer "y" to questions in last step. "Record not found warning" can be
ignored.

## Get source data (dev)

```
$ cd nomis-bulk-to-postgres
$ ./download-data-qs101.sh
```

Note the last step downloads CSV under the "data" directory.  A tiny DB can be
populated by restricting the download to "qs101" but usually
"download-data-qs.sh" would be used for production (see below).

## Populate database

Assumes python venv is setup and "secrets.json" in place as described in [nomis-bulk-to-postgres README](https://github.com/ONSdigital/nomis-bulk-to-postgres/blob/main/README.md)

TODO python should pick up from the env like everything else (?) or just rewrite
TODO migrate this to go code enable to run under AWS (Lambda or Fargate etc.) &
read from S3

* The following imports most data

```
$ . ./bulk/bin/activate
$ python add_to_db.py
```

## Clean up Database

*  Run 
```
$ cd dp-find-insights-poc-api
$ make update-schema
$ ./dataingest/dbsetup/cleandb.sh
 ```

## Geo data import

* Import data from GeoJSON files
  * This puts POLYGON in geo.wkb_geometry

```
$ ./dataingest/spatial/import.sh linux-localhost
```
* Populate geo.wkb_long_lat_geom with long, lat POINT

```
$ cd spacial/longlatgeom
$ go run .
```

## Sanity checks

Various database sanity checks to ensure steps aren't missed.

```
$ make test
```

## producing dump of new prod type database

* Previous steps should be followed but "download-data-qs.sh" would be used

* "pg_dump census > census.sql" can be used to create a dump

* "creatdb.sh" is then run on that dump
