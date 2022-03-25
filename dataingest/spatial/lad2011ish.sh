#!/bin/bash
# fixes SELECT * FROM geo where type_id=4 and lat is null and long is null
# nomis is 2011 and geojson is 2017. Hack 2017 to be 2011ish

DEST=LAD2011ish.geojson
cp 'Local_Authority_Districts_(December_2017)_Boundaries_in_the_UK_(WGS84).geojson' "$DEST"
sed -i s/E06000057/E06000048/ "$DEST"
sed -i s/E07000240/E07000100/ "$DEST"
sed -i s/E07000241/E07000104/ "$DEST"
sed -i s/E07000242/E07000097/ "$DEST"
sed -i s/E07000243/E07000101/ "$DEST"
sed -i s/E08000037/E08000020/ "$DEST"