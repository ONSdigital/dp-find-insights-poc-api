package handlers

import (
	"net/http"
	"testing"
)

func Test_noCache(t *testing.T) {
	var tests = map[string]struct {
		headers []string
		want    bool
	}{
		"no Cache-Control header": {
			nil,
			false,
		},
		"single Cache-Control header w/o no-cache": {
			[]string{"no-store"},
			false,
		},
		"single Cache-Control header with no-cache": {
			[]string{"no-cache"},
			true,
		},
		"multiple Cache-Control headers w/o no-cache": {
			[]string{
				"no-store",
				"max-age=0",
			},
			false,
		},
		"multiple Cache-Control headers with no-cache": {
			[]string{
				"no-cache",
				"max-age=0",
				"must-revalidate",
			},
			true,
		},
	}

	for name, test := range tests {
		req := &http.Request{}
		req.Header = map[string][]string{
			"Cache-Control": test.headers,
		}
		got := noCache(req)
		if got != test.want {
			t.Errorf("%s: %t, want %t", name, got, test.want)
		}
	}
}
