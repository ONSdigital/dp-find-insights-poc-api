#!/bin/bash -x

DB="insights"
USER="$DB"
PW="$DB"
#DATA="data"
DATA="fi-census-data/2011/atlas/viv"

dropdb "$B"
createdb "$B"

if [ ! -d "pgfutter" ]; then
    git clone -b develop git@github.com:stmuk/pgfutter.git
    cd pgfutter && go build
    cd ..
fi

if [ ! -d "$DATA" ]; then
    echo "data dir must exist"
    exit 1
fi

cd "$DATA"
for i in `ls *.csv`; do
    table=${i%\.csv}
    echo $table
    ATLAS=1 ../pgfutter/pgfutter --host localhost -db "$DB" --table $table --schema atlas2011 -user "$USER" --pw "$PW" csv "$i"
done

cd ..

pg_dump -c -C "$DB" > "$DB.sql"
