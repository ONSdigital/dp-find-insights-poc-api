#!/bin/bash
# We don't actually need these env vars but it's more explicit to use them
# directly

set -e

# delete MSOA data from "geo_metric" ~9411707 rows
# XXX slow
psql -c "DELETE FROM geo_metric USING geo WHERE geo_metric.geo_id=geo.id AND geo.type_id=5"

# delete MSOA codes from "geo" ~7201 rows
psql -c "DELETE FROM geo WHERE type_id=5"

# delete Welsh from nomis_category.category_name
psql -c "DELETE FROM  nomis_category WHERE  category_name LIKE '%Cyfradd%ddeiliadaeth%'"

# now we can make it long_nomis_code unique 
# XXX migrate to schema when we can
psql -c "ALTER TABLE nomis_category ADD CONSTRAINT uq_long_nomis_code UNIQUE (long_nomis_code)"

# XXX migrate to schema
psql -c "ALTER TABLE nomis_desc ADD CONSTRAINT uq_short_nomis_code UNIQUE (short_nomis_code)"

# should latter two be indexes?
