package where

import (
	"reflect"
	"testing"
)

func TestParseMultiArgs_Errors(t *testing.T) {
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
		_, err := ParseMultiArgs(test.args)
		if err == nil {
			t.Errorf("%s: expected error, got nil", test.desc)
		}
	}
}

func TestParseMultiArgs_OK(t *testing.T) {
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
		got, err := ParseMultiArgs(test.args)
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: %+v, want %+v", test.desc, got, test.want)
		}
	}
}
