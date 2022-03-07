#!/usr/bin/env bash

set -eu -o pipefail

mkdir -p data

while read -r f; do
    if [[ $f =~ ^# ]]; then
        continue
    fi
    b=$(basename "$f")
    curl "$f" > data/"$b"
    (cd data && unzip "$b")
done < "$1"
