#! /bin/sh

# not sure where these files are used:

for f in ChangeHistory.csv Equivalents.csv
do
    aws --profile development s3 cp s3://find-insights-input-data-files/geoname/$f .
done
curl -o MSOA-Names-1.16.csv https://houseofcommonslibrary.github.io/msoanames/MSOA-Names-1.16.csv
