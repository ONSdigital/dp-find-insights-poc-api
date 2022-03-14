# Tactical EC2 solution 

## EC2 instances

We are currently running the API on two EC2 instances.
This is a Temporary Fix (TM) until we move fully to ONS EC2.

Purpose | Architecture | Instance Hostname
--|--|--
backend dev | amd64 | ec2-18-193-6-194.eu-central-1.compute.amazonaws.com
f/e integration | aarch64 | ec2-35-158-105-228.eu-central-1.compute.amazonaws.com

Unfortunately the instance hostnames will change if the instances are rebooted.
When this happens, the hostnames in the following locations must be updated:

* this doc
* [Makefile](./Makefile)
* the systemd environment for the dp-find-insights-poc-api service (see below)

dp-find-insights-poc-api is running under systemd as the `ubuntu` user, using the
environment for config.

TODO: these instances have been provisioned manually; we should use terraform


## Shell access

The ssh private key is encrypted in `swaggerui/frank-ec2-dev0.pem.gpg`.
Decrypt this into `swaggerui/frank-ec2-dev0.pem`.

To get a shell into the dev instance:

        make ssh-dev

And for the integration instance:

        make ssh-int


## HTTP access

Both instances listen on HTTP port 25252, so you can hit

dev: http://ec2-18-193-6-194.eu-central-1.compute.amazonaws.com:25252/swaggerui

int: http://ec2-35-158-105-228.eu-central-1.compute.amazonaws.com:25252/swaggerui


## How to update the dev EC2 instance

* build binary locally

        make build-linux-amd

* push binary and deploy

        make deploy-dev

* run integration tests

        make test-dev

If tests do not pass, you can rollback to the previous version:

        make rollback-dev


## How to update the f/e integration EC2 instance

* build binary locally

        make build-linux-arm

* push binary and deploy

        make deploy-int

* run integration tests

        make test-int

If tests do not pass, you can rollback to the previous version:

        make rollback-int

## Log monitoring

        make ssh-dev
        journalctl -fu dp-find-insights-poc-api

### SERVICE SETUP

## Systemd service

/etc/systemd/system/dp-find-insights-poc-api.service

added via `systemctl enable dp-find-insights-poc-api`

```
[Unit]Description=Find Insights APIWants=network-online.target
After=network-online.target

[Service]
Type=simple
User=ubuntuSyslogIdentifier=dp-find-insights-poc-api
ExecStart=/home/ubuntu/dp-find-insights-poc-api
Restart=always

[Install]
WantedBy=multi-user.target
```

Environment added via `systemctl edit dp-find-insights-poc-api`

/etc/systemd/system/dp-find-insights-poc-api.service.d/override.conf

```
[Service]
Environment="ENABLE_DATABASE=1
Environment="BIND_ADDR=ec2-18-193-6-194.eu-central-1.compute.amazonaws.com:25252"
Environment="PGUSER=insights"
Environment="PGPASSWORD=XXXXXXXXXXX"
Environment="PGHOST=fi-database-2.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"
Environment="PGPORT=54322"
Environment="PGDATABASE=census"
Environment="ENABLE_CANTABULAR=true"
Environment="CANT_URL=https://ftb-api-ext.ons.sensiblecode.io/graphql"
Environment="CANT_USER=XXXXXXXXXXX"
Environment="CANT_PW=XXXXXXXXXXX"
```
