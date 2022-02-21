# Tactical EC2 solution 

http://ec2-18-193-78-190.eu-central-1.compute.amazonaws.com:25252/swaggerui

This is a Temporary Fix (TM) until we move fully to ONS EC2.

A micro EC2 instance on the free tier was provisioned via the web with an
encrypted version of the private key in this directory.

TODO: should be replaced by terraform

Access via

```
$ ssh -i frank-ec2-dev0.pem ubuntu@ec2-18-193-78-190.eu-central-1.compute.amazonaws.com
```

An encrypted copy of the private key is in this repo.

dp-find-insights-poc-api is running under systemd as the ubuntu user, using the
environment for config.

## Update

* build binary locally and scp the freshly compiled binary to the EC2 instance

    *Do not copy to `dp-find-insights-poc-api` itself*

        make build-linux-amd
        scp build/dp-find-insights-poc-api frank:dp-find-insights-poc-api.new

* ssh into the remote system and run `./deploy.sh` to install the new binary; hit `^C`  or `q` after log displayed

        ssh frank
        ./deploy.sh dp-find-insights-poc-api.new

## Rollback

You can rollback to the previously installed binary:

        ssh frank
        ./deploy.sh previous

## Log monitoring

```
$ journalctl -fu dp-find-insights-poc-api
```

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
Environment="BIND_ADDR=ec2-18-193-78-190.eu-central-1.compute.amazonaws.com:25252"
Environment="PGUSER=insights"
Environment="PGPASSWORD=XXXXXXXXXXX"
Environment="PGHOST=fi-database-2.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"
Environment="PGPORT=54322"
Environment="PGDATABASE=census"
```
