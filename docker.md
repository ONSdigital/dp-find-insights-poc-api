# Running API in local docker

You can run a a local instance of the API for development.
The Go compiler does not need to be installed locally.

The docker image is configured to talk to our AWS RDS instance as the backend database.

## Build the image

```shell
make image
```

You need to rebuild the image if you change any of the Go code in the project.

## Decrypt the database password

The RDS password is encrypted in `secrets/PGPASSWORD.env.asc`.
Decrypt this:

```shell
cd secrets
gpg -d PGPASSWORD.env.asc > PGPASSWORD.env
```
You need to have the `ons-develop` key in your keychain.

## Start the container

In a dedicated terminal:

```shell
make run-api
```

This will leave the container running in the foreground so you can see logs.

To stop the container, hit ^C in the terminal running the container.

## Test

```shell
curl http://localhost:12550/health
```

You should get JSON back.

# Running API and Postgres in local docker

You can run a local instance of the API along with a local instance of Postgres (postgis).
You do not need to install Go or Postgres locally.

## Build the image

As above, remember to build a new image if the Go sources change.

```
make image
```

## Download data

A compressed database dump is held in S3.
Download the gzip into `docker-entrypoint-initdb.d` and uncompress.

```
cd docker-entrypoint-initdb.d
aws --profile development s3 cp s3://find-insights-db-dumps/census2i-20220112.sql.gz .
gunzip census2i-20220112.sql.gz
```

The file unzips to >500MB.

(The files in `docker-entrypoint-initdb.d` are run in alphabetical order.
There is currently a single sql script that creates the `insight` user, and this must be run before the dump is run.
So the user script starts with a number so it sorts earlier.)

## Clear postgres data directory

Postgres holds its data files in a persistent volume mounted from `dbdata`.
When this directory is empty when postgres starts up, sql files and shell scripts in `docker-entrypoint-initdb.d` are run.

If `dbdata` is not empty when postgres starts, postgres will simply start the database as-is.

So when you want to setup a database from the dump file, make sure `dbdata` is empty:

```
rm -rf dbdata
mkdir dbdata
```

## Start the containers

```
docker compose up
```

It will take several minutes for postgres to import the db dump the first time.
But subsequent restarts will not have to reload, and postgres will come up right away.

The API and postgres containers will be left running in the foreground so you can follow logs.
Hit ^C to stop both containers.

## Test

Run a test query and watch the container logs:

```
curl --include http://localhost:12550/dev/hello/census?rows=K04000001\&cols=geography_code,geotype,QS208EW0001
```

## psql

If you want to use psql to connect to the containerised postgres:

```
docker exec -it dp-find-insights-poc-api-db-1 psql -U postgres
```
