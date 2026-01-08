package errors

import (
	"context"
	E "errors"
	"io"
	"net"
	"os"
	"syscall"

	"uni/bridge/tools"
	F "uni/bridge/tools/format"
)

// New returns an error created from the given message arguments.
func New(msg ...any) error {
	return E.New(F.ToString(msg))
}

type extendedError struct {
	msg   string
	cause error
}

// Error returns error info
func (e *extendedError) Error() string {
	if e.cause == nil {
		return e.msg
	}
	return e.cause.Error() + ": " + e.msg
}

// Unwrap returns cause of the error
func (e *extendedError) Unwrap() error {
	return e.cause
}

// Extend creates a new extendedError by wrapping an existing error (cause)
// with an additional message (msg).
//
// If the cause is nil, it panics, as an error cannot be extended without
// a valid underlying cause.
func Extend(cause error, msg ...any) error {
	if cause == nil {
		panic("can't extend a nil error")
	}

	return &extendedError{F.ToString(msg...), cause}
}

// IsClosed checks if the error indicates a closed resource, such as io.EOF,
// net.ErrClosed, io.ErrClosedPipe, os.ErrClosed, or various syscall errors like
// EPIPE, ECONNRESET, and ENOTCONN.
func IsClosed(err error) bool {
	return IsMulti(err, io.EOF, net.ErrClosed, io.ErrClosedPipe, os.ErrClosed, syscall.EPIPE, syscall.ECONNRESET, syscall.ENOTCONN)
}

// IsCanceled checks if the error indicates an operation was canceled or timed out.
// Such as, context.Canceled, context.DeadlineExceeded
func IsCanceled(err error) bool {
	return IsMulti(err, context.Canceled, context.DeadlineExceeded)
}

// Deprecated: 这个函数很奇怪
//
// Cast 看一个 err 是否包含指定的类型 T (T 可能是某种包含错误的类型)
//
// err 是 error 类型
//
// err.(T) 成功，说明 T 实现了 error 接口，err 是 T，返回 (interfaceError，true)
//
// err.(T) 失败，则检查 err 是否实现了 Unwrap() 方法（
//
// 如果 err 实现了 Unwrap() 则调用 Unwrap() 返回 error ，继续循环检查
//
// 如果在整个错误链中都没有找到匹配的类型 T，则返回 T 的默认值和 false
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
