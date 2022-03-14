# import postcode data

CSV from 

https://geoportal.statistics.gov.uk/datasets/a8d42df48f374a52907fe7d4f804a662/about

is also stored at `s3://find-insights-input-data-files/postcode/PCD_OA_LSOA_MSOA_LAD_MAY20_UK_LU.csv`

```
go run .
```

Note: this takes quite a long time ~40mins to populate about 2.3 million rows
and adds a FK relationship with existing `geo` records (and names).

## CSV

Contains three forms of postcode

pcd7: 7-character version of the postcode (e.g. 'BT1 1AA', 'BT486PL')

pcd8: 8-character version of the postcode (e.g. 'BT1  1AA', 'BT48 6PL')

pcds: one space between the district and sector-unit part of the postcode (e.g.
'BT1 1AA', 'BT48 6PL') - possibly the most common formatting of postcodes.

We use the latter and our API works with no space and any number of spaces.  If
spaces are used in the request they should be url encoded.

```
$ curl "http://localhost:25252/msoa/ox4%205aa"  
E02005956, Blackbird Leys

```
