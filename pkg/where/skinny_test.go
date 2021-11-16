package where

import (
	"reflect"
	"testing"
)

func TestGetValues_Errors(t *testing.T) {
	var tests = []struct {
		desc string
		args []string
	}{
		{"empty value", []string{""}},
		{"empty low", []string{"...high"}},
		{"empty high", []string{"low..."}},
		{"too many ellipses", []string{"low...med...hi"}},
	}

	for _, test := range tests {
		_, err := GetValues(test.args)
		if err == nil {
			t.Errorf("%s: expected error, got nil", test.desc)
		}
	}
}

func TestGetValues_OK(t *testing.T) {
	var tests = []struct {
		desc string
		args []string
		want *ValueSet
	}{
		{
			"single value",
			[]string{"val"},
			&ValueSet{
				Singles: []string{"val"},
			},
		},
		{
			"singles separated with comma",
			[]string{"val1,val2"},
			&ValueSet{
				Singles: []string{"val1", "val2"},
			},
		},
		{
			"two separate singles",
			[]string{"val1", "val2"},
			&ValueSet{
				Singles: []string{"val1", "val2"},
			},
		},
		{
			"a range",
			[]string{"lo...hi"},
			&ValueSet{
				Ranges: []*ValueRange{
					{
						Low:  "lo",
						High: "hi",
					},
				},
			},
		},
		{
			"two ranges separated with comma",
			[]string{"lo1...hi1,lo2...hi2"},
			&ValueSet{
				Ranges: []*ValueRange{
					{
						Low:  "lo1",
						High: "hi1",
					},
					{
						Low:  "lo2",
						High: "hi2",
					},
				},
			},
		},
		{
			"two separate ranges",
			[]string{"lo1...hi1", "lo2...hi2"},
			&ValueSet{
				Ranges: []*ValueRange{
					{
						Low:  "lo1",
						High: "hi1",
					},
					{
						Low:  "lo2",
						High: "hi2",
					},
				},
			},
		},
		{
			"mix",
			[]string{"val1,lo1...hi1", "val2,lo2...hi2"},
			&ValueSet{
				Singles: []string{"val1", "val2"},
				Ranges: []*ValueRange{
					{
						Low:  "lo1",
						High: "hi1",
					},
					{
						Low:  "lo2",
						High: "hi2",
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := GetValues(test.args)
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: %+v, want %+v", test.desc, got, test.want)
		}
	}
}

func TestWherePart_Errors(t *testing.T) {
	var tests = []struct {
		desc string
		args []string
	}{
		{
			"error in values",
			[]string{""},
		},
	}

	for _, test := range tests {
		_, err := WherePart("geo", test.args)
		if err == nil {
			t.Errorf("%s: expected error", test.desc)
		}
	}
}

func TestWherePart_OK(t *testing.T) {
	var tests = []struct {
		desc string
		args []string
		want string
	}{
		{
			"a single value",
			[]string{"val"},
			"    col IN ( 'val' )\n",
		},
		{
			"two single values",
			[]string{"val1", "val2"},
			"    col IN ( 'val1', 'val2' )\n",
		},
		{
			"a range",
			[]string{"lo...hi"},
			"    col BETWEEN 'lo' AND 'hi'\n",
		},
		{
			"two ranges",
			[]string{"lo1...hi1", "lo2...hi2"},
			"    col BETWEEN 'lo1' AND 'hi1'\n    OR\n    col BETWEEN 'lo2' AND 'hi2'\n",
		},
		{
			"singles and ranges",
			[]string{"val1,lo1...hi1", "val2,lo2...hi2"},
			"    col IN ( 'val1', 'val2' )\n    OR\n    col BETWEEN 'lo1' AND 'hi1'\n    OR\n    col BETWEEN 'lo2' AND 'hi2'\n",
		},
	}

	for _, test := range tests {
		got, err := WherePart("col", test.args)
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: %s, want %s", test.desc, got, test.want)
		}
	}
}

func TestSkinnyWhere_Error(t *testing.T) {
	var tests = []struct {
		desc string
		rows []string
		cols []string
	}{
		{
			"missing rows or cols",
			nil,
			nil,
		},
		{
			"malformed row",
			[]string{"...range"},
			[]string{"col"},
		},
		{
			"malformed col",
			[]string{"row"},
			[]string{"...range"},
		},
	}

	for _, test := range tests {
		_, err := SkinnyWhere(test.rows, test.cols)
		if err == nil {
			t.Errorf("expected error")
		}
	}
}

func TestSkinnyWhere_OK(t *testing.T) {
	var tests = []struct {
		desc string
		rows []string
		cols []string
		want string
	}{
		{
			"singles row and col",
			[]string{"geo1"},
			[]string{"cat1"},
			`WHERE (
    geography_code IN ( 'geo1' )
) AND (
    category_code IN ( 'cat1' )
)
`,
		},
		{
			"some ranges",
			[]string{"geo1...geo2"},
			[]string{"cat1...cat2"},
			`WHERE (
    geography_code BETWEEN 'geo1' AND 'geo2'
) AND (
    category_code BETWEEN 'cat1' AND 'cat2'
)
`,
		},
		{
			"singles and ranges",
			[]string{"geo1...geo2", "geo3"},
			[]string{"cat1...cat2", "cat3"},
			`WHERE (
    geography_code IN ( 'geo3' )
    OR
    geography_code BETWEEN 'geo1' AND 'geo2'
) AND (
    category_code IN ( 'cat3' )
    OR
    category_code BETWEEN 'cat1' AND 'cat2'
)
`,
		},
	}

	for _, test := range tests {
		clause, err := SkinnyWhere(test.rows, test.cols)
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}
		if clause != test.want {
			t.Errorf("%s:\n%s, want\n%s", test.desc, clause, test.want)
		}
	}
}
