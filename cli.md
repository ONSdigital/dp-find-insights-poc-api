## CLI

To use the cli, you need to set the following environment variables, and provide `--dataset` on the command line.
`--rows` and `--cols` are optional:

| Environment Variable | Example Value |
|---|---|
| `AWS_REGION` | `eu-central-1` |
| `PGHOST` | `fi-database-1.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com` |
| `PGPORT` | `54322` |
| `PGUSER` | `insights` |
| `PASSWORD` | the insights password |
| `PGDATABASE` | `insights` |

To build the cli:

```
$ make build-cli
```

And it ends up as `build/geodata`

So you can do this:

```
$ build/geodata --dataset atlas2011.qs119ew
```

Or:

```
$ build/geodata --dataset atlas2011.qs119ew --rows geography_code:E01000001 --cols geography_code
```
