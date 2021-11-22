package table_test

import (
	"strings"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
)

func TestNew_Error(t *testing.T) {
	_, err := table.New("")
	if err == nil {
		t.Fatal("expected error when New primary is empty")
	}
}

func TestGenerate(t *testing.T) {
	type row struct {
		geo string
		cat string
		val float64
	}

	var tests = []struct {
		desc  string
		input []row
		want  string
	}{
		{
			desc:  "no input",
			input: nil,
			want:  "geography_code\n",
		},
		{
			desc: "single category",
			input: []row{
				{"geo", "cat", 1.23},
			},
			want: "geography_code,cat\ngeo,1.23\n",
		},
		{
			desc: "multiple categories",
			input: []row{
				{"geo", "cat2", 0},
				{"geo", "cat1", 45.6},
			},
			want: "geography_code,cat1,cat2\ngeo,45.6,0\n",
		},
		{
			desc: "multiple geographies",
			input: []row{
				{"geo2", "cat", 7.8},
				{"geo1", "cat", 9.101112},
			},
			want: "geography_code,cat\ngeo1,9.101112\ngeo2,7.8\n",
		},
		{
			desc: "several geographies and categories",
			input: []row{
				{"sun", "mass", 1.98847e+30},
				{"earth", "mass", 5.9722e+24},
				{"moon", "mass", 7.348e+22},
				{"sun", "diameter", 1392000},
				{"earth", "diameter", 12756},
				{"moon", "diameter", 3471},
			},
			want: "geography_code,diameter,mass\nearth,12756,5.9722e+24\nmoon,3471,7.348e+22\nsun,1392000,1.98847e+30\n",
		},
		// specific numeric formatting tests
		{
			desc: "zero should be printed with no decimals or spaces",
			input: []row{
				// zero should be printed as just "0", no decimals or spaces
				{"here", "zero", 0},
			},
			want: "geography_code,zero\nhere,0\n",
		},
		{
			desc: "large integers should be printed as integers, no decimals or spaces",
			input: []row{
				{"here", "millions", 123456789},
			},
			want: "geography_code,millions\nhere,123456789\n",
		},
		{
			desc: "decimals should be printed without trailing zeros",
			input: []row{
				{"here", "decimal", 0.123000},
			},
			want: "geography_code,decimal\nhere,0.123\n",
		},
		{
			desc: "12 digit decimals should be printed",
			input: []row{
				{"here", "12digits", 1.012345678901},
			},
			want: "geography_code,12digits\nhere,1.012345678901\n",
		},
	}

	for _, test := range tests {
		tbl, err := table.New("geography_code")
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}

		// populate table from test.input
		for _, r := range test.input {
			tbl.SetCell(r.geo, r.cat, r.val)
		}

		// generate csv into buf
		var buf strings.Builder
		if err := tbl.Generate(&buf); err != nil {
			t.Fatalf("can't happen: %s", err)
		}

		if buf.String() != test.want {
			t.Errorf("%s:\n%s\nwant:\n%s\n", test.desc, buf.String(), test.want)
		}
	}
}
