#!/bin/bash

set -eu -o pipefail

mkdir -p data

cat 2i.txt | while read f; do
    b=$(basename $f)
    curl "$f" > data/$b
    (cd data && unzip $b)
done
