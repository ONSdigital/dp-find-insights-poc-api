package table_test

import (
	"strings"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
)

func TestGenerate(t *testing.T) {
	type row struct {
		geo     string
		geotype string
		cat     string
		val     float64
	}

	inputSingleCategory := []row{
		{"geo", "type", "cat", 1.23},
	}
	inputMultipleCategories := []row{
		{"geo", "type", "cat2", 0},
		{"geo", "type", "cat1", 45.6},
	}
	inputMultipleGeographies := []row{
		{"geo2", "type", "cat", 7.8},
		{"geo1", "type", "cat", 9.101112},
	}

	var tests = []struct {
		desc    string
		input   []row
		include []string
		want    string
	}{
		{
			desc: "no input and no cols",
			want: "\n",
		},
		{
			desc:    "no input, cols listed",
			include: []string{table.ColGeographyCode, table.ColGeotype},
			want:    "geography_code,geotype\n",
		},
		{
			desc:    "unrecognized include columns just ignored XXX make this an error",
			include: []string{"anything"},
			want:    "\n",
		},

		// single category

		{
			desc:  "single category, no extra cols wanted",
			input: inputSingleCategory,
			want:  "cat\n1.23\n",
		},
		{
			desc:    "single category, geography_code col wanted",
			input:   inputSingleCategory,
			include: []string{table.ColGeographyCode},
			want:    "geography_code,cat\ngeo,1.23\n",
		},
		{
			desc:    "single category, geotype col wanted",
			input:   inputSingleCategory,
			include: []string{table.ColGeotype},
			want:    "geotype,cat\ntype,1.23\n",
		},
		{
			desc:    "single category, geography_code and geotype cols wanted",
			input:   inputSingleCategory,
			include: []string{table.ColGeographyCode, table.ColGeotype},
			want:    "geography_code,geotype,cat\ngeo,type,1.23\n",
		},

		// multiple categories

		{
			desc:  "multiple categories, no extra columns wanted",
			input: inputMultipleCategories,
			want:  "cat1,cat2\n45.6,0\n",
		},
		{
			desc:    "multiple categories, geography_code col wanted",
			input:   inputMultipleCategories,
			include: []string{table.ColGeographyCode},
			want:    "geography_code,cat1,cat2\ngeo,45.6,0\n",
		},
		{
			desc:    "multiple categories, geotype col wanted",
			input:   inputMultipleCategories,
			include: []string{table.ColGeotype},
			want:    "geotype,cat1,cat2\ntype,45.6,0\n",
		},
		{
			desc:    "multiple categories, geography_code and geotype cols wanted",
			input:   inputMultipleCategories,
			include: []string{table.ColGeographyCode, table.ColGeotype},
			want:    "geography_code,geotype,cat1,cat2\ngeo,type,45.6,0\n",
		},

		// multiple geographies

		{
			desc:  "multiple geographies, no extra columns wanted",
			input: inputMultipleGeographies,
			want:  "cat\n9.101112\n7.8\n",
		},
		{
			desc:    "multiple geographies, geography_col wanted",
			input:   inputMultipleGeographies,
			include: []string{table.ColGeographyCode},
			want:    "geography_code,cat\ngeo1,9.101112\ngeo2,7.8\n",
		},
		{
			desc:    "multiple geographies, geotype col wanted",
			input:   inputMultipleGeographies,
			include: []string{table.ColGeotype},
			want:    "geotype,cat\ntype,9.101112\ntype,7.8\n",
		},
		{
			desc:    "multiple geographies, geography_code and geotype columns wanted",
			input:   inputMultipleGeographies,
			include: []string{table.ColGeographyCode, table.ColGeotype},
			want:    "geography_code,geotype,cat\ngeo1,type,9.101112\ngeo2,type,7.8\n",
		},

		{
			desc: "several geographies and categories",
			input: []row{
				{"sun", "gas", "mass", 1.98847e+30},
				{"earth", "solid", "mass", 5.9722e+24},
				{"moon", "solid", "mass", 7.348e+22},
				{"sun", "gas", "diameter", 1392000},
				{"earth", "solid", "diameter", 12756},
				{"moon", "solid", "diameter", 3471},
			},
			include: []string{table.ColGeographyCode, table.ColGeotype},
			want:    "geography_code,geotype,diameter,mass\nearth,solid,12756,5.9722e+24\nmoon,solid,3471,7.348e+22\nsun,gas,1392000,1.98847e+30\n",
		},
		// specific numeric formatting tests
		{
			desc: "zero should be printed with no decimals or spaces",
			input: []row{
				// zero should be printed as just "0", no decimals or spaces
				{"here", "", "zero", 0},
			},
			want: "zero\n0\n",
		},
		{
			desc: "large integers should be printed as integers, no decimals or spaces",
			input: []row{
				{"here", "", "millions", 123456789},
			},
			want: "millions\n123456789\n",
		},
		{

			desc: "decimals should be printed without trailing zeros",
			input: []row{
				{"here", "", "decimal", 0.123000},
			},
			want: "decimal\n0.123\n",
		},
		{
			desc: "12 digit decimals should be printed",
			input: []row{
				{"here", "", "12digits", 1.012345678901},
			},
			want: "12digits\n1.012345678901\n",
		},
	}

	for _, test := range tests {
		tbl := table.New()

		// populate table from test.input
		for _, r := range test.input {
			tbl.SetCell(r.geo, r.geotype, r.cat, r.val)
		}

		// generate csv into buf
		var buf strings.Builder
		if err := tbl.Generate(&buf, test.include); err != nil {
			t.Fatalf("can't happen: %s", err)
		}

		if buf.String() != test.want {
			t.Errorf("%s:\n%s\nwant:\n%s\n", test.desc, buf.String(), test.want)
		}
	}
}
