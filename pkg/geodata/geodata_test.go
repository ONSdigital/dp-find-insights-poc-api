package geodata_test

import (
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
 -- nomis_desc conditions:
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
			desc: "censustable condition with single geography",
			args: geodata.CensusQuerySQLArgs{Geos: []string{"E01000001"}, Censustable: "QS101EW"},
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
 -- nomis_desc conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND nomis_category.nomis_desc_id = nomis_desc.id
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
			desc: "censustable condition with multiple geographies",
			args: geodata.CensusQuerySQLArgs{Geos: []string{"E01000001", "E01000002", "E01000003"}, Censustable: "QS101EW"},
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
 geo.code IN ( 'E01000001', 'E01000002', 'E01000003' )
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- nomis_desc conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND nomis_category.nomis_desc_id = nomis_desc.id
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
			desc: "censustable condition with range of geographies",
			args: geodata.CensusQuerySQLArgs{Geos: []string{"E01000001...E01000005"}, Censustable: "QS101EW"},
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
 geo.code BETWEEN 'E01000001' AND 'E01000005'
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- nomis_desc conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND nomis_category.nomis_desc_id = nomis_desc.id
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
			desc: "censustable condition with multiple and range of geographies",
			args: geodata.CensusQuerySQLArgs{Geos: []string{"E01000001", "E01000002", "E01000003", "E01000005...E010000010"}, Censustable: "QS101EW"},
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
 geo.code IN ( 'E01000001', 'E01000002', 'E01000003' )
 OR
 geo.code BETWEEN 'E01000005' AND 'E010000010'
)
AND geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- nomis_desc conditions:
AND nomis_desc.short_nomis_code = 'QS101EW'
AND nomis_category.nomis_desc_id = nomis_desc.id
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
 -- category conditions:
`,
		},
	}
	for _, test := range tests {
		gotSQL, gotInclude, gotErr := geodata.CensusQuerySQL(test.args)
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
