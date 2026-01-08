package errors

import (
	"errors"
	"strings"

	"uni/bridge/tools"
	F "uni/bridge/tools/format"
)

type MultiError interface {
	Unwrap() []error
}

type multiError struct {
	errors []error
}

func (e *multiError) Error() string {
	return "(" + strings.Join(F.MapToString(e.errors), " | ") + ")"
}

func (e *multiError) Unwrap() []error {
	return e.errors
}

func Errors(errors ...error) error {
	errors = tools.FilterNotNil(errors)
	errors = ExpandAll(errors)
	errors = tools.FilterNotNil(errors)
	errors = tools.UniqBy(errors, error.Error)
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}
	return &multiError{
		errors: errors,
	}
}

func Expand(err error) []error {
	if err == nil {
		return nil
	} else if multiErr, isMultiErr := err.(MultiError); isMultiErr {
		return ExpandAll(tools.FilterNotNil(multiErr.Unwrap()))
	} else {
		return []error{err}
	}
}

func ExpandAll(errs []error) []error {
	return tools.FlatMap(errs, Expand)
}

func Append(err error, other error, fn func(error) error) error {
	if other == nil {
		return err
	}
	return Errors(err, fn(other))
}

// IsMulti checks if the given error (err) matches any error in the targetList.
//
// It returns true if err is one of the errors in targetList, otherwise false.
func IsMulti(err error, targetList ...error) bool {
	for _, target := range targetList {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
