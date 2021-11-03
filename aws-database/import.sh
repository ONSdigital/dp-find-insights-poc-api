#!/bin/bash -x

DB="insights"

cd data
for i in `ls *.csv`; do
    table=${i%\.csv}
    echo $table
    ../pgfutter/pgfutter --host localhost -db "$DB" --table $table --schema atlas2011 -user "$DB" --pw "$DB" csv $i
done

cd ..

pg_dump -c -C "$DB" > "$DB.sql"
