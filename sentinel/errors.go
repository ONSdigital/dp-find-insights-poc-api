package sentinel

type Sentinel string

const (
	ErrMissingParams     = Sentinel("missing parameter")
	ErrInvalidParams     = Sentinel("invalid parameter")
	ErrTooManyMetrics    = Sentinel("too many metrics")
	ErrNoContent         = Sentinel("no data found")
	ErrPartialContent    = Sentinel("insufficient data found")
	ErrNotSupported      = Sentinel("not supported")
	ErrTableName         = Sentinel("empty table name")
	ErrInconsistentTypes = Sentinel("inconsistent property types")
	ErrUnusableType      = Sentinel("unusable property type")
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
