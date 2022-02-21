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
