package errors

import (
	"context"
	E "errors"
	"guard/bridge/tools"
	F "guard/bridge/tools/format"
	"io"
	"net"
	"os"
	"syscall"
)

// New returns an error created from the given message arguments.
func New(msg ...any) error {
	return E.New(F.ToString(msg))
}

type extendedError struct {
	msg   string
	cause error
}

func (e *extendedError) Error() string {
	if e.cause == nil {
		return e.msg
	}
	return e.cause.Error() + ": " + e.msg
}

func (e *extendedError) Unwrap() error {
	return e.cause
}

func Extend(cause error, msg ...any) error {
	if cause == nil {
		panic("can't extend a nil error")
	}

	return &extendedError{F.ToString(msg...), cause}
}

func Cause(cause error, msg ...any) error {
	if cause == nil {
		panic("can't cause on a nil error")
	}
	return &causeError{F.ToString(msg...), cause}
}

func Cause1(err error, cause error) error {
	if cause == nil {
		panic("cause on an nil error")
	}
	return &causeError1{err, cause}
}

func IsClosed(err error) bool {
	return IsMulti(err, io.EOF, net.ErrClosed, io.ErrClosedPipe, os.ErrClosed, syscall.EPIPE, syscall.ECONNRESET, syscall.ENOTCONN)
}

func IsCanceled(err error) bool {
	return IsMulti(err, context.Canceled, context.DeadlineExceeded)
}

func Cast[T any](err error) (T, bool) {
	if err == nil {
		return tools.DefaultValue[T](), false
	}

	for {
		interfaceError, isInterface := err.(T)
		if isInterface {
			return interfaceError, true
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err == nil {
				return tools.DefaultValue[T](), false
			}
		case interface{ Unwrap() []error }:
			for _, innerErr := range x.Unwrap() {
				if interfaceError, isInterface = Cast[T](innerErr); isInterface {
					return interfaceError, true
				}
			}
			return tools.DefaultValue[T](), false
		default:
			return tools.DefaultValue[T](), false
		}
	}
}
