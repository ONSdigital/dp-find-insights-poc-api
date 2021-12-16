# SWAGGERUI

This is a tactical swaggerui url which allows access to the Find Insights
Backend Team Census API via a web form.

http://ec2-18-193-78-190.eu-central-1.compute.amazonaws.com:25252/swaggerui

With the caveats that this is a work in progress which we are iterating on
rapidly and that some of the naming is currently confusing and subject to
change we thought it useful to give some visibility of the sort of queries
possible.

The part of the interface which is useful is the "/dev/hello/{dataset}"
endpoint. Expanding this will give you a list of parameters starting with
"dataset".  This isn't actually a dataset anymore but specifies the query type.

For actual use currently

https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/census

should be used (but has no swaggerui endpoint)

## Query Type ("dataset")

Possible values are "skinny" and "census" (these names will change)

* "skinny" refers to an old query type

* "census" refers to an new query type

The new "census" query type probably should be avoided for the time being.

## Skinny Queries

We can use the integration tests as a crib for the parameters (eg. rows, cols etc.)

https://github.com/ONSdigital/dp-find-insights-poc-api/blob/develop/inttests/main.go

"rows" are one, more or ranges of ONS geographical codes (eg. E01000001). LSOA
and LAD are currently available.  These can be specified in the "geocode"
parameter as well.

"cols" are the long NOMIS category codes from the NOMIS Bulk system (eg. QS118EW0011). The QS
(quick statistics are available).

 Warning these are different NOMIS codes to those used in the NOMIS API. There
is an extra 0 in these NOMIS codes and 0001 refers to total rather than the
first metric.

## Ranges

eg. QS118EW0001...QS118EW0011
eg. E01001111...E01001211

## Census Queries

These aren't fully tested and under development but will support PostGIS
geography types in the database which allow things like long,lat to be used in
bounding box queries etc. which resemble the "skinny" ones.

eg. bbox=-0.370947083400182,51.3624781092781,0.17687729439413147,51.67377813346024

There is some description of these queries in

https://docs.google.com/presentation/d/1rJYG9JIKFsFsgXU-JW16-d3Exw97551EgthULbV44mE/edit#slide=id.g105f3797443_0_0s://docs.google.com/presentation/d/1rJYG9JIKFsFsgXU-JW16-d3Exw97551EgthULbV44mE/edit#slide=id.g105f3797443_0_0
