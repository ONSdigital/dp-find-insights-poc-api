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

The RDS password is encrypted in `PGPASSWORD.env.asc`.
Decrypt this:

```shell
gpg -d < PGPASSWORD.env.asc > PGPASSWORD.env
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
