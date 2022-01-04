package geodata_test

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/kylelemons/godebug/diff"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

func TestGetCensusQuery(t *testing.T) {
	var tests = []struct {
		desc        string
		args        geodata.CensusQuerySQLArgs
		wantSQL     string
		wantInclude []string
		wantErr     error
	}{
		{
			desc:    "no arguments",
			args:    geodata.CensusQuerySQLArgs{},
			wantErr: errors.New("must specify a condition (rows, bbox, location/radius, and/or polygon)"),
		},
		// Rows
		{
			desc: "rows condition only",
			args: geodata.CensusQuerySQLArgs{Geos: []string{"E01000001"}},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
`,
		},
		// Bounding Box
		{
			desc: "bbox condition only",
			args: geodata.CensusQuerySQLArgs{BBox: "-0.370947083400182,51.3624781092781,0.17687729439413147,51.673778133460246"},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
WHERE (
 -- geo conditions:
geo.wkb_geometry && ST_GeomFromText(
 'MULTIPOINT(-0.370947 51.362478, 0.176877 51.673778)',
 4326
)
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
`,
		},
		{
			desc:    "bbox error - non-numeric data",
			args:    geodata.CensusQuerySQLArgs{BBox: "X,51.3624781092781,0.17687729439413147,51.673778133460246"},
			wantErr: errors.New("error parsing bbox \"X,51.3624781092781,0.17687729439413147,51.673778133460246\": strconv.ParseFloat: parsing \"X\": invalid syntax"),
		},
		{
			desc:    "bbox error - too few coords",
			args:    geodata.CensusQuerySQLArgs{BBox: "-0.370947083400182,51.3624781092781"},
			wantErr: errors.New("valid bbox is 'lon,lat,lon,lat', received \"-0.370947083400182,51.3624781092781\": invalid parameter"),
		},
		{
			desc:    "bbox error - too many coords",
			args:    geodata.CensusQuerySQLArgs{BBox: "-0.370947083400182,51.3624781092781,0.17687729439413147,51.673778133460246,-0.370947083400182,0.17687729439413147,51.673778133460246"},
			wantErr: errors.New("valid bbox is 'lon,lat,lon,lat', received \"-0.370947083400182,51.3624781092781,0.17687729439413147,51.673778133460246,-0.370947083400182,0.17687729439413147,51.673778133460246\": invalid parameter"),
		},
		// Columns
		{
			desc: "single col condition",
			args: geodata.CensusQuerySQLArgs{
				Geos: []string{"E01000001"},
				Cols: []string{"QS119EW0002"},
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.long_nomis_code IN ( 'QS119EW0002' )
)
`,
		},
		{
			desc: "censustable condition with single geography",
			args: geodata.CensusQuerySQLArgs{
				Geos:        []string{"E01000001"},
				Censustable: "QS101EW",
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
 , nomis_desc
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.nomis_desc_id = nomis_desc.id
)
 `,
		},
		{
			desc: "censustable condition with single col",
			args: geodata.CensusQuerySQLArgs{
				Geos:        []string{"E01000001"},
				Censustable: "QS101EW",
				Cols:        []string{"QS119EW0002"},
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
 , nomis_desc
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.long_nomis_code IN ( 'QS119EW0002' )
 OR
 nomis_category.nomis_desc_id = nomis_desc.id
)
`,
		},
		{
			desc: "censustable condition with multiple col",
			args: geodata.CensusQuerySQLArgs{
				Geos:        []string{"E01000001"},
				Censustable: "QS101EW",
				Cols: []string{
					"QS119EW0001",
					"QS119EW0002",
					"QS119EW0003",
				},
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
 , nomis_desc
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.long_nomis_code IN ( 'QS119EW0001', 'QS119EW0002', 'QS119EW0003' )
 OR
 nomis_category.nomis_desc_id = nomis_desc.id
)
`,
		},
		{
			desc: "censustable condition with ranged col",
			args: geodata.CensusQuerySQLArgs{
				Geos:        []string{"E01000001"},
				Censustable: "QS101EW",
				Cols:        []string{"QS119EW0001...QS119EW0004"},
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
 , nomis_desc
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.long_nomis_code BETWEEN 'QS119EW0001' AND 'QS119EW0004'
 OR
 nomis_category.nomis_desc_id = nomis_desc.id
)
`,
		},
		{
			desc: "censustable condition with multiple col and range col",
			args: geodata.CensusQuerySQLArgs{
				Geos:        []string{"E01000001"},
				Censustable: "QS101EW",
				Cols: []string{
					"QS119EW0001",
					"QS119EW0002",
					"QS119EW0003",
					"QS117EW0001...QS117EW0003",
				},
			},
			wantSQL: `
SELECT
 geo.code AS geography_code,
 geo_type.name AS geotype,
 nomis_category.long_nomis_code AS category_code,
 geo_metric.metric AS value
FROM
 geo,
 geo_type,
 geo_metric,
 data_ver,
 nomis_category
 , nomis_desc
WHERE (
 -- geo conditions:
 geo.code IN ( 'E01000001' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
AND (
 nomis_category.long_nomis_code IN ( 'QS119EW0001', 'QS119EW0002', 'QS119EW0003' )
 OR
 nomis_category.long_nomis_code BETWEEN 'QS117EW0001' AND 'QS117EW0003'
 OR
 nomis_category.nomis_desc_id = nomis_desc.id
)
`,
		},
	}
	for _, test := range tests {
		ctx := context.Background()
		gotSQL, gotInclude, gotErr := geodata.CensusQuerySQL(ctx, test.args)
		normedGotSql := normSQL(gotSQL)
		normedWantSql := normSQL(test.wantSQL)
		if !reflect.DeepEqual(normedGotSql, normedWantSql) {
			t.Errorf("%s: returned SQL differs from expected:  %s", test.desc, diff.Diff(normedWantSql, normedGotSql))
		}
		if !reflect.DeepEqual(gotInclude, test.wantInclude) {
			t.Errorf("%s: got these geography column values - '%s', wanted '%s'", test.desc, gotInclude, test.wantInclude)
		}
		if test.wantErr != nil {
			if !reflect.DeepEqual(gotErr.Error(), test.wantErr.Error()) {
				t.Errorf("%s: got this error = '%s', wanted '%s'", test.desc, gotErr, test.wantErr)
			}
		} else if gotErr != nil {
			t.Errorf("%s: got this error - '%s', wanted nil", test.desc, gotErr)
		}
	}
}

// normalise sql to single spaces and newlines only, with no blank lines
func normSQL(sql string) string {
	// replace all whitespace except newlines with single space
	fieldDetector := func(c rune) bool {
		return unicode.IsSpace(c) && string(c) != "\n"
	}
	wsNormed := strings.Join(strings.FieldsFunc(sql, fieldDetector), " ")
	// strip blank lines and return
	multiNewlinePattern := regexp.MustCompile(`\n\s*\n+`)
	return multiNewlinePattern.ReplaceAllString(wsNormed, "\n")
}
