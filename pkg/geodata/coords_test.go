package geodata

import (
	"errors"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	"github.com/stretchr/testify/assert"
)

func Test_ParseCoords_Err(t *testing.T) {
	var tests = map[string]string{
		"empty string":          "",
		"non-numeric string":    "x",
		"non-numeric component": "1,x",
		"empty component":       ",1",
	}

	for name, s := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := parseCoords(s)
			if !errors.Is(err, sentinel.ErrInvalidParams) {
				t.Errorf("%s, want %s", err, sentinel.ErrInvalidParams)
			}
		})
	}
}

func Test_ParseCoords(t *testing.T) {
	var tests = map[string]struct {
		s      string
		coords []float64
	}{
		"single coord": {
			"3.14",
			[]float64{3.14},
		},
		"two coords": {
			"1.23,4.56",
			[]float64{1.23, 4.56},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			coords, err := parseCoords(test.s)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, coords, test.coords)
		})
	}
}

func Test_CheckValidCoords_Err(t *testing.T) {
	var tests = map[string][]float64{
		"only one coord":       {1.23},
		"odd number of coords": {1.2, 3.4, 5.6},
		"long out of range":    {181, 0, 0, 0},
		"lat out of range":     {0, 91, 0, 0},
	}

	for name, coords := range tests {
		t.Run(name, func(t *testing.T) {
			err := checkValidCoords(coords)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func Test_CheckValidCoords(t *testing.T) {
	var tests = map[string][]float64{
		"two coords":  {1.23, 4.56},
		"four coords": {1, 2, 3, 4},
	}

	for name, s := range tests {
		t.Run(name, func(t *testing.T) {
			err := checkValidCoords(s)
			if !assert.NoError(t, err) {
				return
			}
		})
	}
}

func Test_asLineString_Err(t *testing.T) {
	_, err := asLineString([]float64{0})
	if !errors.Is(err, sentinel.ErrInvalidParams) {
		t.Errorf("got %s, expected %s", err, sentinel.ErrInvalidParams)
	}
}

func Test_asLineString_OK(t *testing.T) {
	var tests = []struct {
		coords []float64
		want   string
	}{
		{[]float64{1, 2}, "1 2"},
		{[]float64{1, 2, 3, 4}, "1 2,3 4"},
	}

	for i, test := range tests {
		got, err := asLineString(test.coords)
		if err != nil {
			t.Errorf("Test %d: %s", i, err)
			continue
		}
		if got != test.want {
			t.Errorf("Test %d: %s, want %s", i, got, test.want)
		}
	}
}
