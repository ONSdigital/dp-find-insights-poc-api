package where

import (
	"testing"
)

func TestWherePart_OK(t *testing.T) {
	var tests = []struct {
		desc string
		args []string
		want string
	}{
		{
			"no values",
			[]string{},
			"",
		},
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
		set, err := ParseMultiArgs(test.args)
		if err != nil {
			t.Errorf("%s: %s\n", test.desc, err)
			continue
		}
		got := WherePart("col", set)
		if got != test.want {
			t.Errorf("%s: %s, want %s", test.desc, got, test.want)
		}
	}
}
