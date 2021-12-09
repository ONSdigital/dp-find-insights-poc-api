package geodata_test

import (
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

func Test_String(t *testing.T) {
	var tests = []struct {
		c    geodata.Point
		want string
	}{
		{geodata.Point{}, "0 0"},
		{geodata.Point{Lon: 0.1214, Lat: 51.4910}, "0.1214 51.491"},
	}

	for _, test := range tests {
		got := test.c.String()
		if got != test.want {
			t.Errorf("%s, want %s", got, test.want)
		}
	}
}

func Test_ParsePolygonErrors(t *testing.T) {
	var tests = []struct {
		desc string
		s    string
	}{
		{
			desc: "empty input",
			s:    "",
		},
		{
			desc: "odd number of numbers",
			s:    "1,2,3,4,5,6,7,8,9",
		},
		{
			desc: "not enough coordinates in polygon",
			s:    "1,2",
		},
		{
			desc: "first and last don't match",
			s:    "1,2,3,4,5,6,7,8",
		},
		{
			desc: "lon not a number",
			s:    "a,2,3,4,5,6,1,2",
		},
		{
			desc: "lat not a number",
			s:    "1,a,3,4,5,6,1,2",
		},
	}

	for _, test := range tests {
		_, err := geodata.ParsePolygon(test.s)
		if err == nil {
			t.Errorf("%s: expected error", test.desc)
		}
	}
}

func Test_ParsePolygon(t *testing.T) {
	var tests = []struct {
		desc string
		s    string
		want []geodata.Point
	}{
		{
			desc: "sanity check",
			s:    "0.1214,51.4910,0.1017,51.4647,0.1338,51.4635,0.1214,51.4910",
			want: []geodata.Point{
				{0.1214, 51.4910},
				{0.1017, 51.4647},
				{0.1338, 51.4635},
				{0.1214, 51.4910},
			},
		},
	}

	for _, test := range tests {
		got, err := geodata.ParsePolygon(test.s)
		if err != nil {
			t.Errorf("%s: %s", test.desc, err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: %+v, want %+v", test.desc, got, test.want)
		}
	}
}
