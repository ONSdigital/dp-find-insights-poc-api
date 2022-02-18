# geodata CLI

The `geodata` cli is a thin layer over the `geodata` and `metadata` packages.
It does everything the API does, but takes query parameters on the command line instead of through an http interface, and exits after printing results. It works against a census database like curl works against the API.

The cli can be used to test the `geodata` and `metadata` packages without building and firing up an API somewhere.
Results are printed exactly as they are returned from the library functions.

To use the cli, you must set the standard postgres environment variables to point to your target database.
The files `api-rds.env` and `api-docker.env` are good starting points.

To build the cli:

    $ make build-cli

And it ends up as `build/geodata`

So you can do this, for example:

    $ build/geodata metadata -year 2011

cli subcommands access different methods within the `geodata` and `metadata` packages.
subcommand arguments more or less correspond to method arguments.

subcommand | method
--|--
`ckmeans` | `geodata.CKmeans`
`metadata` | `metadata.Get`
`query` | `geodata.Query`

You can get help for each subcommand with `-h`, eg:

    $ build/geodata query -h
    
## Examples

    $ geodata query -year 2011 \
        -rows K04000001 \
        -cols geography_code,geotype,QS208EW0001

    $ geodata ckmeans -year 2011 -cat QS101EW0001 -geotype LSOA -k 5

    $ geodata ckmeansratio -year 2011 \
        -cat1 QS101EW0002 -cat2 QS101EW0001 \
        -geotype LSOA \
        -k 5
