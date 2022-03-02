# LIVE DATABASE UPDATE

These are beta quality instructions and a knowledge of the system is still
needed.

TODO investigate AWS RDS migration

## create a local dump of the new database

* populate `dataingest/addtodb/data` from `2i.txt` via `download-data-2i.sh`

* populate the environment with (local non-live) PG_ vars (see docker.md), 
  eg. `export PGDATABASE=census_new` etc

* run `indigestion.sh`

```
cd dataingest/dbsetup
pg_dump census_new > census-20220301.sql (or whatever today's date is!)
```

## import local dump onto live system using a different db name to the live database

* populate the environment with (LIVE) PG_ vars (see docker.md) but with PGDATABASE=census_new 

* you will also need to set `export POSTGRES_PASSWORD=XXX` to the live `postgres` (admin) accn password

(be careful at this point!)

```
./creatdb.sh census-20220301.sql    
```

* this will load data into `census_new` on the live database server (there may be permissions errors which can ignored at this point since they are fixed up by `make update-schema` automatically in the final line of the script.

## cut over 'census_new' to 'census'

* Announce download to slack channel(s)

* Anyone with db client connections to the live db server should disconnect. It should be possible to restart RDS via the AWS interface to drop connections if this isn't possible.

ssh to the EC2 server as described in TACTICALEC2.md to drop app db connections.

sudo to root and

```
systemctl stop dp-find-insights-poc-api.service
```

Then into another terminal connected to the live db
```
alter database census rename to census_to20220301 -- (or whatever today's date is!)
alter database census_new rename to census;
```

Back to the SSH server.

```
systemctl start dp-find-insights-poc-api.service
systemctl status  dp-find-insights-poc-api.service -l

```

## check responses look sane

```
curl -s  "http://ec2-18-193-78-190.eu-central-1.compute.amazonaws.com:25252/health"| jq . 
curl -s  "http://ec2-18-193-78-190.eu-central-1.compute.amazonaws.com:25252/metadata/2011"|jq .
```
