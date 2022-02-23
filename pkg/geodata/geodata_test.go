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
	"github.com/stretchr/testify/assert"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
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
			wantErr: geodata.ErrMissingParams,
		},
		// Rows
		{
			desc: "rows condition only",
			args: geodata.CensusQuerySQLArgs{
				Year: 2011,
				Geos: []string{"E01000001"},
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
    geo.code IN ( 'E01000001' )
)
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
 -- category conditions:
`,
		},
		// Bounding Box
		{
			desc: "bbox condition only",
			args: geodata.CensusQuerySQLArgs{
				Year: 2011,
				BBox: "-0.370947083400182,51.3624781092781,0.17687729439413147,51.673778133460246",
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
geo.wkb_geometry && ST_GeomFromText(
 'MULTIPOINT(-0.370947 51.362478, 0.176877 51.673778)',
 4326
)
)
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
 -- category conditions:
`,
		},
		{
			desc:    "bbox error - non-numeric data",
			args:    geodata.CensusQuerySQLArgs{BBox: "X,51.3624781092781,0.17687729439413147,51.673778133460246"},
			wantErr: geodata.ErrInvalidParams,
		},
		{
			desc:    "bbox error - too few coords",
			args:    geodata.CensusQuerySQLArgs{BBox: "-0.370947083400182,51.3624781092781"},
			wantErr: geodata.ErrInvalidParams,
		},
		{
			desc:    "bbox error - too many coords",
			args:    geodata.CensusQuerySQLArgs{BBox: "-0.370947083400182,51.3624781092781,0.17687729439413147,51.673778133460246,-0.370947083400182,0.17687729439413147,51.673778133460246"},
			wantErr: geodata.ErrInvalidParams,
		},
		// Columns
		{
			desc: "single col condition",
			args: geodata.CensusQuerySQLArgs{
				Year: 2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
 -- category conditions:
AND (
 nomis_category.long_nomis_code IN ( 'QS119EW0002' )
)
			`,
		},
		{
			desc: "censustable condition with single geography",
			args: geodata.CensusQuerySQLArgs{
				Year:        2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
 -- category conditions:
AND (
 nomis_category.nomis_desc_id = nomis_desc.id
)
 `,
		},
		{
			desc: "censustable condition with single col",
			args: geodata.CensusQuerySQLArgs{
				Year:        2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
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
				Year:        2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
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
				Year:        2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
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
				Year:        2011,
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND (
 geo.code IN ( 'E01000001' )
)
AND nomis_desc.short_nomis_code = 'QS101EW'
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
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
		{
			desc:    "all rows, too many tokens",
			args:    geodata.CensusQuerySQLArgs{Geos: []string{"x", "all"}},
			wantErr: geodata.ErrInvalidParams,
		},
		{
			desc: "all rows, all categories",
			args: geodata.CensusQuerySQLArgs{
				Year: 2011,
				Geos: []string{"all"},
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
WHERE geo.valid
AND geo_type.id = geo.type_id
 -- geotype conditions:
 -- geo conditions:
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
 -- category conditions:
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
			if !errors.Is(gotErr, test.wantErr) {
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

func TestExtractSpecialCols_Error(t *testing.T) {
	// catch error when a special column is named in a range
	var tests = map[string]struct {
		cols []string // input cols as received from query string
	}{
		"special col at beginning of range": {
			[]string{"foo," + table.ColGeotype + "...high"},
		},
		"special col at end of range": {
			[]string{"foo,low..." + table.ColGeotype},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse cols into a ValueSet
			set, err := where.ParseMultiArgs(test.cols)
			if !assert.NoError(t, err) {
				return
			}

			// split into special columns and new ValueSet
			_, _, err = geodata.ExtractSpecialCols(set)
			assert.Error(t, err)
		})
	}
}
func TestExtractSpecialCols_OK(t *testing.T) {
	var tests = map[string]struct {
		cols         []string // input cols as received from query string(s)
		wantIncludes []string // expected list of special cols found
		wantCols     []string // expected list of non-special cols
	}{
		"no query strings": {
			cols:         []string{},
			wantIncludes: nil,
			wantCols:     nil,
		},
		"single non-special col": {
			cols:         []string{"foo"},
			wantIncludes: nil,
			wantCols:     []string{"foo"},
		},
		"single special col": {
			cols:         []string{table.ColGeotype},
			wantIncludes: []string{table.ColGeotype},
			wantCols:     nil,
		},
		"special and non-special cols": {
			cols:         []string{"foo", table.ColGeotype},
			wantIncludes: []string{table.ColGeotype},
			wantCols:     []string{"foo"},
		},
		"special and non-special cols with range": {
			cols:         []string{"foo", table.ColGeotype, "low...high"},
			wantIncludes: []string{table.ColGeotype},
			wantCols:     []string{"foo,low...high"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse cols into a ValueSet
			set, err := where.ParseMultiArgs(test.cols)
			if !assert.NoError(t, err) {
				return
			}

			// split into special columns and new ValueSet
			includes, newset, err := geodata.ExtractSpecialCols(set)
			if !assert.NoError(t, err) {
				return
			}

			// make sure extracted special column list is correct
			if !assert.Equal(t, test.wantIncludes, includes) {
				return
			}

			// parse wantcols into a ValueSet so we can compare with newset
			wantset, err := where.ParseMultiArgs(test.wantCols)
			if !assert.NoError(t, err) {
				return
			}

			// verify resulting ValueSet after extraction
			assert.Equal(t, wantset, newset)
		})
	}
}

func TestValidateAllToken_Error(t *testing.T) {
	// catch error when a ALL is given more than once or in a range
	var tests = map[string]struct {
		cols []string // input cols as received from query string
	}{
		"ALL is given more than once": {
			[]string{"ALL", "ALL}"},
		},
		"ALL is low part of a range": {
			[]string{"ALL...high"},
		},
		"ALL is high part of a range": {
			[]string{"low...ALL"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse cols into a ValueSet
			set, err := where.ParseMultiArgs(test.cols)
			if !assert.NoError(t, err) {
				return
			}

			// call validator
			err = geodata.ValidateAllToken(set)
			assert.Error(t, err)
		})
	}
}

func TestValidateAllToken_OK(t *testing.T) {
	var tests = map[string]struct {
		cols []string // input cols as received from query string
	}{
		"ALL is not given at all": {
			[]string{"foo"},
		},
		"ALL is only token": {
			[]string{"ALL"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse cols into a ValueSet
			set, err := where.ParseMultiArgs(test.cols)
			if !assert.NoError(t, err) {
				return
			}

			// call validator
			err = geodata.ValidateAllToken(set)
			assert.NoError(t, err)
		})
	}
}

func TestFixgeotype_Error(t *testing.T) {
	// catch invalid geotypes
	_, err := geodata.FixGeotype("foo")
	assert.Error(t, err)
}

func TestFixgeotype_OK(t *testing.T) {
	// verify geotypes are fixed
	geotype, err := geodata.FixGeotype("country")
	assert.NoError(t, err)
	assert.Equal(t, geotype, "Country")
}

func TestMapGeotypes_Error(t *testing.T) {
	// catch error when a geotype isn't valid
	var tests = map[string]struct {
		geos []string // input geos as received from query string
	}{
		"invalid single geotype": {
			[]string{"foo"},
		},
		"geotype as low part of range": {
			[]string{"lsoa...foo"},
		},
		"geotype as high part of range": {
			[]string{"foo...lsoa"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse geos into a ValueSet
			set, err := where.ParseMultiArgs(test.geos)
			if !assert.NoError(t, err) {
				return
			}

			// expect validation error
			_, err = geodata.MapGeotypes(set)
			assert.Error(t, err)
		})
	}
}

func TestMapGeotypes_OK(t *testing.T) {
	var tests = map[string]struct {
		geos []string // input geos as received from query string
		want []string // expected fixed geos
	}{
		"single geotype": {
			[]string{"lsoa"},
			[]string{"LSOA"},
		},
		"multiple geotypes": {
			[]string{"eW", "country", "LaD"},
			[]string{"EW,Country,LAD"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// parse geos into a ValueSet
			set, err := where.ParseMultiArgs(test.geos)
			if !assert.NoError(t, err) {
				return
			}

			// expect validation ok
			fixedset, err := geodata.MapGeotypes(set)
			if !assert.NoError(t, err) {
				return
			}

			// parse our expected string to get a ValueSet
			wantset, err := where.ParseMultiArgs(test.want)
			if !assert.NoError(t, err) {
				return
			}

			// expect
			assert.Equal(t, wantset, fixedset)
		})
	}
}
