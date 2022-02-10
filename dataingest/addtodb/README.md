# Instructions to populate new database with (Nomis) Data
# 2021 synthetic data to come!

```
$ ./download-data-2i.sh
$ cd ../.. && make update-schema 
$ go run ./dataingest/addtodb

```
