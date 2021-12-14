package geodata_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

type getCensusQuerySQLArgs struct {
	ctx         context.Context
	geos        []string
	censustable string
	bbox        string
	location    string
	radius      int
	polygon     string
	geotypes    []string
	cols        []string
}

func TestGetCensusQuery(t *testing.T) {
	var tests = []struct {
		desc        string
		args        getCensusQuerySQLArgs
		wantSQL     string
		wantInclude []string
		wantErr     error
	}{
		{
			desc:    "no arguments",
			args:    getCensusQuerySQLArgs{},
			wantErr: errors.New("must specify a condition (rows,bbox,location/radius, or polygon)"),
		},
		{
			desc: "rows condition only",
			args: getCensusQuerySQLArgs{geos: []string{"TestCensusRow"}},
			wantSQL: geodata.NormSQL(`
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
	geo.code IN ( 'TestCensusRow' )
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
			`),
		},
	}
	for _, test := range tests {
		gotSQL, gotInclude, gotErr := geodata.GetCensusQuerySQL(
			test.args.ctx,
			test.args.geos,
			test.args.bbox,
			test.args.location,
			test.args.radius,
			test.args.polygon,
			test.args.geotypes,
			test.args.cols,
		)
		if !reflect.DeepEqual(gotSQL, test.wantSQL) {
			t.Errorf("%s: got this SQL\n %s, wanted this SQL\n %s", test.desc, gotSQL, test.wantSQL)
		}
		if !reflect.DeepEqual(gotInclude, test.wantInclude) {
			t.Errorf("%s: got these geography column values %s, wanted %s", test.desc, gotInclude, test.wantInclude)
		}
		if test.wantErr != nil {
			if !reflect.DeepEqual(gotErr.Error(), test.wantErr.Error()) {
				t.Errorf("%s: got this error %s, wanted %s", test.desc, gotErr, test.wantErr)
			}
		}
	}
}
