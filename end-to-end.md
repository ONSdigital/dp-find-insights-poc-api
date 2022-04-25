# End-to-End How To

This is a complete step-by-step guide to all the steps related to building a database and testing an API against it.

Some steps do not have to be done every single time, but are included here for completeness.

Unless otherwise stated, `cd` commands are relative to a checked out `dp-find-insights-poc-api` repo.

Environment variables are important! Always verify the `PG_*` and `POSTGRES_PASSWORD` variables before a risky step.

## Prerequisites

* [dp-find-insights-poc-api](https://github.com/ONSdigitsl/dp-find-insights-poc-api) and
  [dp-setup](https://github.com/ONSdigital/dp-setup) repos checked out

* postgres cli tools installed

    (I use [Postgres.app](https://postgresapp.com) and add `/Applications/Postgres.app/Contents/Versions/latest`
    to `PATH` and `/Applications/Postgres.app/Contents/Versions/latest/share/man` to `MANPATH`.)

* aws cli installed and able to access AWS `develop` environment

    [AWS cli](https://aws.amazon.com/cli/)

    [AWS Credentials](https://github.com/ONSdigital/dp/blob/main/guides/AWS_CREDENTIALS.md)

* [go compiler](https://go.dev/dl/) installed

* [terraform](https://www.terraform.io/downloads) installed

* [develop gpg key](https://github.com/ONSdigital/dp-ci/tree/master/gpg-keys/environments) (must get passphrase from tech lead)

* [ssh access](https://github.com/ONSdigital/dp-cli) to the `develop` environment

## Contents

[Create local Postgres instance](#create-local-postgres-instance)

[Download source data files](#download-source-data-files)

[Run the ingest](#run-the-ingest)

[Export Database](#export-database)

[Provision RDS](#provision-rds)

[Set up SSH tunnel](#set-up-ssh-tunnel)

[Import into RDS](#import-into-rds)

[Start local API](#start-local-api)

[Run API tests against RDS](#run-api-tests-against-rds)

[Tear down RDS](#tear-down-rds)

[Tear down local postgres instance](#tear-down-local-postgres-instance)

## Create local Postgres instance

    export PGDATA=~/ons/postgres/pgdata

    # It is possible to create the instance that does not need a superuser
    # password from localhost, but our scripts are meant to be used against any
    # type of instance, so for consistency we create the local instance to
    # expect a superuser password.
    export POSTGRES_PASSWORD=supassword

    initdb --auth password --username postgres --pwfile <(echo $POSTGRES_PASSWORD)

    # start server and log to file instead of stdout
    pg_ctl -l ~/ons/postgres/logfile start

## Download source data files

    # test S3 access
    aws --profile development s3 ls s3://find-insights-input-data-files

    cd dataingest/addtodb
    ./download-data-2i.sh

    cd ../geoname
    ./download.sh

    cd ../spatial
    ./download.sh

    cd ../postcode
    ./download.sh

    cd ../..

## Run the ingest

The ingest scripts are meant to talk to a local postgres instance, so set the environment variables as below.

`PGPASSWORD` is the password for the `insights` user, and `POSTGRES_PASSWORD` is the postgres superuser password
(our convention copied from the Postgres docker image).

    export PGHOST=localhost
    export PGPORT=5432
    export PGDATABASE=census
    export PGUSER=insights
    export PGPASSWORD=apassword

    PGPASSWORD=$POSTGRES_PASSWORD psql --dbname postgres --username postgres -c "CREATE USER $PGUSER WITH PASSWORD '$PGPASSWORD' CREATEDB"

    ./indigestion.sh

## Export Database

This export step also talks to the local postgres instance, so keep the environment variables from the Ingest step above.

    cd dataingest/dbsetup
    today=$(date +%Y%m%d)
    pg_dump --exclude-table=spatial_ref_sys > census-${today}.sql

## Provision RDS

This step doesn't need to be done if there is already an RDS instance in the `develop` environment.

    cd <dp-setup-repo>
    cd terraform/dp-geodata-api-postgres
    gpg -d < develop.tfvars.asc > develop.tfvars

    terraform init
    terraform workspace select develop

    terraform plan -var-file=develop.tfvars

## Set up SSH tunnel

You need the RDS instance hostname and port when you set up the tunnel.
You can grab this information like this:

    cd <dp-setup-repo>
    cd terraform/dp-geodata-api-postgres 
    terraform output

Once you have the RDS instance endpoint details, and run this in a separate terminal on your workstation:

    dp ssh develop web 1 --port 5555:<endpoint-host>:<endpoint-port>

Now back to your working terminal, set the environment to point to the tunnel and to use login details for the RDS instance.
The superuser password is found in `develop.tfvars`.

        export PGHOST=localhost
        export PGPORT=5555
        export PGDATABASE=census
        export PGUSER=insights
        export PGPASSWORD=<insights-password>
        export POSTGRES_PASSWORD=<superuser-password>

        # test connectivity to db; this should list 4 databases
        $ PGPASSWORD=$POSTGRES_PASSWORD psql --dbname postgres --username postgres -c '\l'

## Import into RDS

Verify your environment variables point to the tunnel and use the RDS credentials, then:

    cd <dp-find-insights-poc-api repo>
    cd dataingest/dbsetup
    
    # this is the sql dump file created in the Export Database step above
    ./creatdb.sh census-20220422.sql

## Start local API

In another terminal on your workstation, set the environment to point to the tunnel and
use the RDS credentials.
Then build and start a local API:

    make build
    export ENABLE_DATABASE=1
    make debug

## Run API tests against RDS

This step doesn't need any special environment variables.
The only requirement is a running local API that talks to RDS.

Just run this in another terminal.  You should not get any errors.

    cd <dp-find-insights-poc-api repo>
    make test-local

## Tear down RDS

If you need to tear down the RDS instance:

* stop any running local APIs
* log out of the tunnel shell to web-1
* exit from any psql sessions to RDS

Then

    cd <dp-setup repo>
    cd terraform/dp-geodata-api-postgres

    terraform destroy -var-file=develop.tfvars

## Tear down local Postgres instance

    export PGDATA=~/ons/postgres/pgdata
    pg_ctl stop
    # optionally remove logfile and pgdata directory
