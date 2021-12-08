# Data Ingest Processes

WARNING these processes are in rapid flux as of Dec 2021 and subject to change.

## aws-database

Various scripts to provision & load data into a AWS RDS (Postgres 13.4)
instance used by the Find Insights back-end team.

Most are dependent on the existance of Postgres client utilities being
installed and also a configured aws command line client.

* create.env.asc
  * Encrypted version of postgres password - currently same for "postgres" (admin
user) & "insights" (app user)

Decrypt
```
gpg -d create.env.asc
```

to create "create.env"

* awscreate.sh
  * Tactical solution to create AWS RDS instance and security group opening (non-standard) postgres port of 54322.
  * Probably should be migrated to Terraform.

* awsloaddata.sh
  * creates "insights" pg user & imports "census.sql" DB dump when ran with '-create-user' flag
  * ommitting the flag just imports census.sql"
    * 'dropdb census && createdb census' should be ran first in the latter case

# Provision Tiny dev DB locally & create prod DB dump

These instructions are for the "new" or "skinny" database.

TODO automate

Assumes postgres is running on developer desktop natively or via docker with pg
client tools (eg. psql, pgdump, createdb, dropdb etc.) in the PATH

## Create "insights" Postgres user if needed

This only has to be done once.

```
$ psql postgres
postgres=# CREATE USER insights WITH PASSWORD 'insights';
postgres=# ALTER USER insights WITH CREATEDB;
```

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

## Create database and schema

```
$ createdb censustiny
$ psql -U postgres -d censustiny -c "CREATE EXTENSION postgis"
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

## populate database

Assumes python venv is setup and "secrets.json" in place as described in [nomis-bulk-to-postgres README](https://github.com/ONSdigital/nomis-bulk-to-postgres/blob/main/README.md)

TODO python should pick up from the env like everything else (?) or just rewrite
TODO migrate this to go code enable to run under AWS (Lambda or Fargate etc.) &
read from S3

* The following imports most data

```
$ . ./bulk/bin/activate
$ python add_to_db.py
```

* Fix up "geo" table data 

Setup 2 source CSV using [dp-find-insights-poc-api/dataingest/geo/README.md](https://github.com/ONSdigital/dp-find-insights-poc-api/blob/develop/dataingest/geo/README.md)

```
$ cd dp-find-insights-poc-api/dataingest/geo
$ go run .
```

## test SQL queries

* Result from python process

```
$ psql censustiny
[..]
censustiny=> select metric from geo_metric where id=1;
   metric   
------------
 56075912.0
```

* Result from geo fixup

```
censustiny=> select * from geo where id=1;
 id | type_id |   code    |       name        
----+---------+-----------+-------------------
  1 |       1 | K04000001 | England and Wales
```

## producing dump of new prod type database

* Previous steps should be followed but "download-data-qs.sh" would be used

* "pg_dump census > census.sql" can be used to create a dump

* "awsloaddata.sh" is then run on that dump
