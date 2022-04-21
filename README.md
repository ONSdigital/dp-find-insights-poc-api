# dp-find-insights-poc-api

Census Atlas geodata API and database.

### Further Docs

* Environment Variables

    Database connection details and feature flags are controlled with environment variables.
    Almost all processes depend on these variables being set correctly, so be careful to verify before running anything risky.

    * [common](docker.md#environment)
    * [some extra API-related variables](#configuration) (below)

* Postgres/PostGIS

    The system is designed around Postgres with PostGIS extensions.
    The database may be run as a native local process, within a container, or as an RDS instance.

    * [local container](docker.md#running-postgis-in-a-container)
    * local native instance (docs to be written)
    * [RDS provisioning](../dp-setup/terraform/dp-geodata-api-postgres/) (see dp-setup repo)

* Data Ingest

    The "ingest" process takes source data files and loads them into the database.
    For best results, the ingest should be applied to a local native postgres instance.
    Running against RDS directly or against a containerised postgres instance is too slow.

    Source data files should be in place, and the postgres environment variables must be exported before running the ingest script (`indigestion.sh`).

    * source data files
        * [S3 buckets](dataingest/S3-BUCKETS.md)
        * [postcode data](dataingest/postcode/README.md)
        * [ONS geo codes](dataingest/geoname/README.md)
        * [nomis data](dataingest/addtodb/README.md)
        * [spatial data](dataingest/spatial/README.md)
    * [running](dataingest/dbsetup/README.md)

* Export/Import

    The export/import processes copies a database.
    Our main application is to export a locally ingested database for import into RDS.

    * [export/import procedure RDS](dataingest/dbsetup/README.md)
    * [export/import local db](docker.md#importing-a-db-dump)

* API

    The API presents a specific interface for querying the database.
    It is mainly used by the front end application.

    The API may be run as a local process, within a container, or on some temporary EC2 instances we have in place for the moment.
    The plan is that instance of the API will run within the standard develop and prod environments.

    * [building](docker.md#building-images-and-binaries)
    * [deploying to EC2](TACTICALEC2.md)
    * [running locally](docker.md#running-the-api)
    * [sanity checking](docker.md#sanity-checking-the-api)

* CLIs

    These clis are mostly for developers, and may not be fully caught up with features of the API.

    * [geodata cli](cli.md)
    * [cantabular cli](cmd/cantabular/README.md)
    * [geobb](geobb/README.md) (to generate static `geoLookup.json` file)

* Testing
    * [api unit tests](Makefile)
    * [component tests](Makefile)
    * [integration tests](inttests/README.md)

* Terraform
    * [S3 and misc](terraform/README.md)
    * [RDS](../dp-setup/terraform/dp-geodata-api-postgres/) (see dp-setup repo)

* How-Tos
    * [quick start running local API](docker.md#a-id"quick-starts"a-quick-starts)
### Getting started

* Run `make debug`

To enable postgres and census queries, set `ENABLE_DATABASE`, and the postgres environment variables.
To lookup the postgres password in Secrets Manager, set `FI_PG_SECRET_ID` instead of `PGPASSWORD`,
and make sure to set `AWS_REGION` and any other AWS environment variables.

### Auto generated code

`swagger.yaml` is used to generate code via `make generate`

Particularly api/api.go (and similar files) shouldn't be directly edited.

### <a id="configuration"></a> Configuration ###

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :25252    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| ENABLE_DATABASE              | false     | Enable postgres and census query functionality
| AWS_REGION                   |           | used by AWS SDK when ENABLE_DATABASE is true and PGPASSWORD is empty
| PGHOST                       |           | postgres host when ENABLE_DATABASE is true
| PGPORT                       |           | postgres port when ENABLE_DATABASE is true
| PGUSER                       |           | postgres user when ENABLE_DATABASE is true
| PGPASSWORD                   |           | postgres password when ENABLE_DATABASE is true (also see FI_PG_SECRET_ID)
| PGDATABASE                   |           | postgres database when ENABLE_DATABASE is true
| FI_PG_SECRET_ID              |           | ARN of key holding postgres password if PGPASSWORD is empty

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

