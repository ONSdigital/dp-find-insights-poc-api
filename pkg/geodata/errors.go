package geodata

type Sentinel string

const (
	ErrMissingParams  = Sentinel("missing parameter")
	ErrTooManyMetrics = Sentinel("too many metrics")
	ErrInvalidTable   = Sentinel("invalid table")
	ErrNoContent      = Sentinel("no data found")
)

func (e Sentinel) Error() string {
	return string(e)
}

func (e Sentinel) Is(err error) bool {
	sentinel, ok := err.(Sentinel)
	if !ok {
		return false
	}
	return sentinel == e
}
