#! /bin/sh

for f in \
    'Lower_Layer_Super_Output_Areas_(December_2011)_Boundaries_Super_Generalised_Clipped_(BSC)_EW_V3.geojson' \
    'Middle_Layer_Super_Output_Areas_(December_2011)_Boundaries_Super_Generalised_Clipped_(BSC)_EW_V3.geojson' \
    'Local_Authority_Districts_(December_2017)_Boundaries_in_the_UK_(WGS84).geojson'
do
    aws --profile development s3 cp s3://find-insights-input-data-files/geojson/"$f" .
done
