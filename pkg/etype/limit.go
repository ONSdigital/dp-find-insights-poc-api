package etype

import "fmt"

// Limit errors happen when an internal limit is exceeded.
type LimitError struct {
	msg string
	e   error
}

func Limit(format string, args ...interface{}) *LimitError {
	return &LimitError{
		msg: fmt.Sprintf(format, args...),
	}
}

func LimitWrap(err error, format string, args ...interface{}) *LimitError {
	msg := fmt.Sprintf(format, args...)
	if err != nil && len(format) > 0 && format[len(format)-1] == ':' {
		msg += err.Error()
	}
	return &LimitError{
		msg: msg,
		e:   err,
	}
}

func (e LimitError) Error() string {
	return e.msg
}

func (e LimitError) Is(err error) bool {
	_, ok := err.(*ParamError)
	return ok
}

func (e LimitError) Unwrap() error {
	return e.e
}
