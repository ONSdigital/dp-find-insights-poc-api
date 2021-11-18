# dp-find-insights-poc-api
Experimental application for developing CI/CD pipeline

### Getting started

* Run `make debug`

To enable postgres and census queries, set `ENABLE_DATABASE`, and the postgres environment variables.
To lookup the postgres password in Secrets Manager, set `FI_PG_SECRET_ID` instead of `PGPASSWORD`,
and make sure to set `AWS_REGION` and any other AWS environment variables.

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

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

