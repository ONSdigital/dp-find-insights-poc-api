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
		val string
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
				{"geo", "cat", "val"},
			},
			want: "geography_code,cat\ngeo,val\n",
		},
		{
			desc: "multiple categories",
			input: []row{
				{"geo", "cat2", "val2"},
				{"geo", "cat1", "val1"},
			},
			want: "geography_code,cat1,cat2\ngeo,val1,val2\n",
		},
		{
			desc: "multiple geographies",
			input: []row{
				{"geo2", "cat", "val2"},
				{"geo1", "cat", "val1"},
			},
			want: "geography_code,cat\ngeo1,val1\ngeo2,val2\n",
		},
		{
			desc: "several geographies and categories",
			input: []row{
				{"sun", "mass", "1.98847e30"},
				{"earth", "mass", "5.9722e24"},
				{"moon", "mass", "7.348e22"},
				{"sun", "diameter", "1392000"},
				{"earth", "diameter", "12756"},
				{"moon", "diameter", "3471"},
			},
			want: "geography_code,diameter,mass\nearth,12756,5.9722e24\nmoon,3471,7.348e22\nsun,1392000,1.98847e30\n",
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
