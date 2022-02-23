package where

import (
	"errors"
	"strings"
)

// ParseMultiArgs generates a ValueSet from multi value query parameters.
//
// XXX return proper error types once errors are factored out of geodata
func ParseMultiArgs(args []string) (*ValueSet, error) {
	set := NewValueSet()
	for _, arg := range args {

		// each argument may have many values or ranges separated by commas
		tokens := strings.Split(arg, ",")
		for _, token := range tokens {

			// each token may be a single value or a range
			if !strings.Contains(token, "...") {
				if token == "" {
					return nil, errors.New("value must not be empty")
				}
				set.AddSingle(token)

			} else {
				r := strings.Split(token, "...")
				if len(r) != 2 {
					return nil, errors.New("range must be low...high")
				}
				if r[0] == "" || r[1] == "" {
					return nil, errors.New("range values must not be empty")
				}
				set.AddRange(r[0], r[1])
			}
		}
	}
	return set, nil
}
