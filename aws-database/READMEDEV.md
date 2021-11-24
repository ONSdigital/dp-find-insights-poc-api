# Provision Tiny dev DB locally

These instructions are for the "new" or "skinny" database.

TODO automate

Assumes postgres is running on developer desktop natively or via docker.

## Create "insights" Postgres user if needed

```
$ psql postgres
postgres=# CREATE USER insights WITH PASSWORD "insights";
postgres=# ALTER USER insights WITH CREATEDB;
```

## 
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
$ cd dp-find-insights-poc-api
$ make update-schema
```

Answer "y" to questions in last step. "Record not found warning" can be
ignored.

## Get source data


```
$ cd nomis-bulk-to-postgres
$ ./download-data-qs101.sh
```

Note the last step downloads CSV under the "data" directory.  A tiny DB can be
populated by restricting the download to "qs101" but usually
"download-data-qs.sh" would be used.

## populate database

Assumes python venv is setup and "secrets.json" in place as described in [nomis-bulk-to-postgres README](https://github.com/ONSdigital/nomis-bulk-to-postgres/blob/main/README.md)

TODO python should pick up from the env like everything else

The following imports most data
```
$ . ./bulk/bin/activate
$ python add_to_db.py
```

Fix up "geo" table data (setup using README.md in following directory)

```
$ cd dp-find-insights-poc-api/dataingest/geo
$ go run .
```

## test SQL queries

Result from python process

```
$ psql censustiny
[..]
censustiny=> select metric from geo_metric where id=1;
   metric   
------------
 56075912.0
```

Result from geo fixup

```
censustiny=> select * from geo where id=1;
 id | type_id |   code    |       name        
----+---------+-----------+-------------------
  1 |       1 | K04000001 | England and Wales
```

