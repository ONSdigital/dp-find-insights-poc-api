package where

import (
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
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
					return nil, fmt.Errorf("%w: value must not be empty", sentinel.ErrInvalidParams)
				}
				set.AddSingle(token)

			} else {
				r := strings.Split(token, "...")
				if len(r) != 2 {
					return nil, fmt.Errorf("%w: range must be low...high", sentinel.ErrInvalidParams)
				}
				if r[0] == "" || r[1] == "" {
					return nil, fmt.Errorf("%w: range values must not be empty", sentinel.ErrInvalidParams)
				}
				set.AddRange(r[0], r[1])
			}
		}
	}
	return set, nil
}
