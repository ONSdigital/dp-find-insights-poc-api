package where

// ValueSet holds all the single values and ranges for a for a multi-valued column.
type ValueSet struct {
	Singles []string      // list of single values; becomes IN
	Ranges  []*ValueRange // list of value ranges; becomes BETWEEN
}

// ValueRange holds the low and high values for a range.
type ValueRange struct {
	Low  string
	High string
}

func NewValueSet() *ValueSet {
	return &ValueSet{}
}

func (set *ValueSet) AddSingle(s string) {
	set.Singles = append(set.Singles, s)
}

func (set *ValueSet) AddRange(low, high string) {
	r := &ValueRange{
		Low:  low,
		High: high,
	}
	set.Ranges = append(set.Ranges, r)
}

// A callback function operates on a Single or on a Range.
// single will be non-nil for a Single, and low and high will be non-nil for
// a Range.

// When operating on a Single, the value of s on return can be:
// 	nil - delete this Single
// 	ptr to a new value - changes this Single
// 	ptr to the current value - does not change this Single
//
// Similarly when operating on a Range, the value of l and h on return can be:
//	nil - delete this Range
//	ptrs to new values - change this Range
//	ptrs to current values - does not change this Range
//
// See the example changeCallback in where_test.go.
//
type Callback func(single, low, high *string) (s *string, l *string, h *string, e error)

// Walk calls the callback function for every Single and Range in set.
// The callback function can return values telling Walk to delete, change, or
// leave the value untouched.
// If the callback function returns an error, Walk will stop and return
// the error, returning a nil ValueSet.
// The original ValueSet will not be changed on error.
//
func (set *ValueSet) Walk(callback Callback) (*ValueSet, error) {
	newset := NewValueSet()

	for _, single := range set.Singles {

		newsingle, _, _, err := callback(&single, nil, nil)
		if err != nil {
			return nil, err
		}
		if newsingle != nil {
			newset.AddSingle(*newsingle)
		}
	}

	for _, vr := range set.Ranges {
		_, newlow, newhigh, err := callback(nil, &vr.Low, &vr.High)
		if err != nil {
			return nil, err
		}
		if newlow != nil && newhigh != nil {
			newset.AddRange(*newlow, *newhigh)
		}
	}

	return newset, nil
}
