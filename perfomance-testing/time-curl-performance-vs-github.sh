#!/bin/bash
output_file=$1
do_gzip=${2-0}
if [[ "$do_gzip" -ne 0 ]];
then
  echo "flat-file,api,api_gz" > "$output_file"
else
  echo "flat-file,api" > "$output_file"
fi

while IFS= read -r dataset;
do
  # avoid firing on blank lines
  if [[ "$dataset" == "" ]]; then continue; fi

  # process dataset name into data table name and column number
  dataset_name_len=$((${#dataset}))
  data_table=$(echo ${dataset:0:$(($dataset_name_len-3))} | tr '[:upper:]' '[:lower:]')
  data_col_raw=${dataset:$(($dataset_name_len-3)):3}
  data_col=_$((10#${data_col_raw}))

  # get time to curl flat file for dataset
  flat_file_url="https://jtrim-ons.github.io/census-2011-data/data-with-lad-and-ew-rows/${dataset}.csv"
  time_flat_file=$(curl --output output_flat_file.csv --silent --write-out '%{time_total}' $flat_file_url)

  # get time to curl same data from api
  api_url="https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.${data_table}?cols=geography_code,total,${data_col}"
  time_api=$(curl --output output_api.csv --silent --write-out '%{time_total}' $api_url)

  # get time to curl same data from api with gzip, if doing that, and output. Just output if not doing gzip.
  if [[ "$do_gzip" -ne 0 ]];
  then
    time_api_gz=$(curl --output output_api_gz.csv --silent --write-out '%{time_total}' --header 'Accept-Encoding:gzip'  $api_url)
    echo "${time_flat_file},${time_api},${time_api_gz}" >> "$output_file"
  else
    echo "${time_flat_file},${time_api}" >> "$output_file"
  fi


  # output results

done < datasets.txt
