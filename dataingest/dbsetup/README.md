# LIVE DATABASE UPDATE

This is how to export a local census database and import it into another postgres instance.

These are beta quality instructions and a knowledge of the system is still
needed.

TODO investigate AWS RDS migration.  Maybe possible to upload a compressed census.sql type file more rapidly than using `creatdb.sh` (?)

## create new database locally

The ingest runs locally via `indigestion.sh`.
The prerequisites are:
* The source data files must be downloaded into the various ingest directories.
  See their respective `README.md` and `download.sh` files.
  For example, populate `dataingest/addtodb/data` from `2i.txt` via `download-data-2i.sh`.

* populate the environment with (local non-live) PG_ vars (see docker.md), 
  eg. `export PGDATABASE=census_new` etc

* run `indigestion.sh`

## create a local dump of the new database

(`pg_dump` dumps the database named in `$PGDATABASE`)

```
cd dataingest/dbsetup
pg_dump --exclude-table=spatial_ref_sys > census-$(date +%Y%m%d).sql
```

## create an SSH tunnel to the RDS instance

Skip this step if your workstation can contact the postgres instance directly.

If you are importing into an RDS instance in the `develop` environment, not visible directly, then set up an SSH tunnel.

```
dp ssh develop web 1 --port 5555:RDS:5432
```

where `RDS` is the hostname of the RDS instance.

This creates a tunnel connecting localhost:5555 to RDS:5432.

## import local dump onto live system using a different db name to the live database

* populate the environment with (LIVE) PG_ vars (see docker.md) but with PGDATABASE=census_new 

* you will also need to set `export POSTGRES_PASSWORD=XXX` to the live `postgres` (admin) accn password

* If you are using the SSH tunnel above, then set:

```
PGHOST=localhost
PGPORT=5555
```

(be careful at this point! triple check your `PG_*` and `POSTGRES_PASSWORD` variables)

```
./creatdb.sh census-$(date +%Y%m%d).sql
```

* this will load data into `$PGDATABASE` on the live database server (there may be permissions errors which can ignored at this point since they are fixed up by `make update-schema` automatically in the final line of the script.  This can take quite some time.

## cut over 'census_new' to 'census' on EC2 instances

(to cut over an RDS instance, see the step below)

* Announce back-end downtime to slack channel(s) to alert front-end devs.

* Anyone with db client connections to the live db server should disconnect. It should be possible to restart RDS via the AWS interface to drop connections if this isn't possible.

If you are importing into a database on one of our temporary EC2 servers, do one of `make ssh-dev` or `make ssh-int`:

```
make ssh-int
```

* sudo to root and
```
systemctl stop dp-find-insights-poc-api.service
```

* Then into another terminal connected to the live db
```
alter database census rename to census_to20220301 -- (or whatever today's date is!)
alter database census_new rename to census;
```

* Back to the SSH server.

```
systemctl start dp-find-insights-poc-api.service
systemctl status dp-find-insights-poc-api.service -l

```
* Back on your local host, check the API response looks sane

```
make health-int
```

## cut over 'census_new' to 'census' on RDS instance

We do not yet have any APIs pointing to the RDS instance, so you don't need to bounce apps.

* set `PG_*` and `POSTGRES_PASSWORD` variables, noting that you want to contact `localhost:5555`

* connect to the RDS instance as the postgres superuser

```
PGPASSWORD=$POSTGRES_PASSWORD psql -U postgres -d postgres
```

* rotate the databases as above:
Then into another terminal connected to the live db
```
alter database census rename to census_to20220301 -- (or whatever today's date is!)
alter database census_new rename to census;
```
