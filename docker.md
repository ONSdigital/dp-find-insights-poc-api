# Docker Stack Goals

* let front end developers run a local API and database without Go or postgres installed
* let backend developers work on the API and ingest processes locally
* minimise differences between the API running locally and in EC2
* minimise the differences between postgres running locally and in RDS

# Environment

Everything is controlled with environment variables.
You should be able to set a handful of environment variables and treat
components the same no matter where they are running.

## Variables

The normal postgres variables are used by clients.
These variables are used if you run utilities like `psql` locally, and when
you run ingest scripts.

These variables are also passed on to containers.
So when you run the API in a container, it picks up these variables from your
current environment.

* `PGHOST`
* `PGPORT`
* `PGDATABASE`
* `PGUSER`
* `PGPASSWORD`

In addition, `POSTGRES_PASSWORD` holds the postgres superuser password.
This is used in two ways:
* The postgis image uses it when initialising a new database.
* Ingest scripts use it when they need `postgres` credentials.

Two other variables are used by processes running in containers who want to
talk to a postgres also running in a container.
These are necessary because of the way docker does networking.
* `PGHOST_INTERNAL`
* `PGPORT_INTERNAL`

Usually `PGHOST_INTERNAL=host.docker.internal` works.

## RDS Password

The password for the RDS `insights` user is in `secrets/PGPASSWORD.env.asc`.

You need to decrypt this file if you are going to point the API at RDS:

    cd secrets
    gpg -d PGPASSWORD.env.asc > PGPASSWORD.env

Since docker `.env` files look almost like shell scripts, in most cases you
can do this:

    . secrets/PGPASSWORD.env
    export PGPASSWORD

## Environment Files

Two example files can be sourced to set environment variables for common cases.
These are used in the examples below.

* `api-docker.env` -- use when you are using postgis in a container
* `api-rds.env` -- use when you are talking to RDS

# Building Images and Binaries

You can build the API image without a local Go compiler:

    make image

The image is named `dp-find-insights-poc-api:latest`.

You can also build an API binary to run as a local process.
This requires a Go compiler.

    make build

The binary is `build/dp-find-insights-poc-api`.

# Running the API

You can run the api as a local process or within a container.
In both cases they will stay in the foreground so you can see logs.

Hit '^C' to stop.

As a local process:

    . api-rds.env
    make debug

In a container:

    . api-rds.env
    docker compose up api

The API listens on port 25252.

If you use `api-docker.env` instead of `api-rds.env`, you can access a local
postgres instance that has been populated by a dump file or through the
ingest process.

# Sanity Checking the API

You can run a quick sanity check on the API you just started:

    curl http://localhost:25252/health

    curl http://localhost:25252/metadata/2011

You should get JSON back.
If you pipe the healthcheck output through `jq` you can clearly see if it is ok.

```
$ curl -s http://localhost:25252/health | jq
{
  "status": "OK",
  "version": {
    "build_time": "2022-01-28T07:20:58Z",
    "git_commit": "73f0856c9b8e9d7e2e5dd1dfa00f4f86997d54c7",
    "language": "go",
    "language_version": "go1.17.6",
    "version": ""
  },
  "uptime": 49303,
  "start_time": "2022-01-28T11:46:38.110117Z",
  "checks": [
    {
      "name": "postgres",
      "status": "OK",
      "message": "pgx healthy",
      "last_checked": "2022-01-28T11:47:07.961363Z",
      "last_success": "2022-01-28T11:47:07.961363Z",
      "last_failure": null
    },
    {
      "name": "gorm",
      "status": "OK",
      "message": "gorm healthy",
      "last_checked": "2022-01-28T11:47:08.212101Z",
      "last_success": "2022-01-28T11:47:08.212101Z",
      "last_failure": null
    }
  ]
}
```

# Running postgis in a container

To start a local database:

    . api-docker.env
    docker compose up db

Persistent data is held in `./dbdata`.
This directory will be created if it doesn't exist.
Zap this directory when you want to start completey from scratch.

When a database is created from scratch, the superuser password is set to the
current value of `POSTGRES_PASSWORD`.

Postgres will listen on `PGPORT` on the host's locahost interface.
But it always listens on 5432 internally.

# Running psql from a container

You can run `psql` without installing postgres locally with the `psql.sh`
wrapper.
This scripts invokes `psql` from within a postgres container with the following
settings:

* The normal `PG*` variables are passed along to the container
* `PGHOST_INTERNAL` and `PGPORT_INTERNAL` are used if they are present
* The current working directory is mapped to /tmp in the container, and the
  container starts with workdir set to /tmp.

So in most cases you an pass `-f ./file.sql` and it will work right.

# Running update-schema from a container

You can run `update-schema` against a database without a local Go compiler.

1. Build the update-schema image

        make update-schema-image

2. Setup environment

        . api-docker.env

3. Run the update-schema image

        make run-update-schema

# Importing a DB dump

You can set up a local database with a dump taken from another database.
This is a "quick" way to setup a local stack for front end development.

1. Download the dump file

    A compressed database dump is held in S3.
    Download the gzip and uncompress.

        aws --profile development --region eu-central-1 s3 cp s3://find-insights-db-dumps/census-20220118.sql.gz .
        gunzip census-20220118.sql.gz

    The file unzips to >500MB.

2. Shutdown any locally running postgres and remove the `./dbdata` directory

3. Set up the environment

        . api-docker.env

4. Start postgis

        docker compose up db

5. In another terminal, create the `insights` user and `census` database

        . api-docker.env
        PGPASSWORD=$POSTGRES_PASSWORD ./psql.sh --dbname postgres -U postgres -f sql/pre-restore.sql

6. Import the dump

    Run the restore as the superuser because object ownership is set within
    the dump file.

        PGPASSWORD=$POSTGRES_PASSWORD ./psql.sh -U postgres -f census-20220118.sql

Smoke test the database by starting an API and running the sanity tests.

# Running the ingest processes

The full ingest process can be run against a local postgis instance.


1. Shutdown any locally running postgres and remove the `./dbdata` directory

2. Set up the environment

        . api-docker.env

3. Start postgis

        docker compose up db

4. Follow the normal [dataingest](dataingest/README.md) instructions.

Once the database is imported, you should be able start an API container or
local API process using the same environment variables.
