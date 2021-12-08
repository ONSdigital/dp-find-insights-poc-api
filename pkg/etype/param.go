package etype

import "fmt"

// Param errors are those that result from missing or malformed parameters.
type ParamError struct {
	msg string
	e   error
}

func Param(format string, args ...interface{}) *ParamError {
	return &ParamError{
		msg: fmt.Sprintf(format, args...),
	}
}

func ParamWrap(err error, format string, args ...interface{}) *ParamError {
	msg := fmt.Sprintf(format, args...)
	if err != nil && len(format) > 0 && format[len(format)-1] == ':' {
		msg += err.Error()
	}
	return &ParamError{
		msg: msg,
		e:   err,
	}
}

func (e ParamError) Error() string {
	return e.msg
}

func (e ParamError) Is(err error) bool {
	_, ok := err.(*ParamError)
	return ok
}

func (e ParamError) Unwrap() error {
	return e.e
}
