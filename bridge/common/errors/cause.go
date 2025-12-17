package errors

import F "guard/bridge/tools/format"

type causeError struct {
	message string
	cause   error
}

// Error returns error info
func (e *causeError) Error() string {
	return e.message + ": " + e.cause.Error()
}

// Unwrap returns cause of the error
func (e *causeError) Unwrap() error {
	return e.cause
}

// Cause returns a *causeError
func Cause(cause error, msg ...any) error {
	if cause == nil {
		panic("can't cause on a nil error")
	}
	return &causeError{F.ToString(msg...), cause}
}

type causeError1 struct {
	// This field stores the reason that caused the current error.
	error
	cause error
}

// Error returns error info
func (e *causeError1) Error() string {
	return e.error.Error() + ": " + e.cause.Error()
}

// Unwrap returns cause of the error
func (e *causeError1) Unwrap() []error {
	return []error{e.error, e.cause}
}

// Cause returns a *causeError1
func Cause1(err error, cause error) error {
	if cause == nil {
		panic("cause on an nil error")
	}
	return &causeError1{err, cause}
}
