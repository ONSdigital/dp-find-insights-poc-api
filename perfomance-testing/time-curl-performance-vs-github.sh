#!/bin/bash

echo "flat-file,api" > perf.csv

while IFS= read -r dataset;
do
  if [[ "$dataset" == "" ]]; then continue; fi
  dataset_len=$((${#dataset}))
  data_table=$(echo ${dataset:0:$(($dataset_len-3))} | tr '[:upper:]' '[:lower:]')
  data_col_raw=${dataset:$(($dataset_len-3)):3}
  data_col=_$((10#${data_col_raw}))
  flat_file_url="https://jtrim-ons.github.io/census-2011-data/data-with-lad-and-ew-rows/${dataset}.csv"
  time_flat_file=$(curl --output output_flat_file.csv --silent --write-out '%{time_total}' $flat_file_url)
  api_url="https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.${data_table}?cols=geography_code,total,${data_col}"
  time_api=$(curl --output output_api.csv --silent --write-out '%{time_total}' $api_url)
  echo "${time_flat_file},${time_api}" >> flat_file_vs_api_perf.csv
done < datasets.txt
