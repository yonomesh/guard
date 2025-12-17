package errors

import (
	"errors"
	"net"
)

type TimeoutError interface {
	Timeout() bool
}

// Deprecated: use IsTimeoutError instead
func IsTimeout(err error) bool {
	var netErr net.Error

	if errors.As(err, &netErr) {
		//nolint:staticcheck
		return netErr.Temporary() && netErr.Timeout() // 兼容一些老代码
	}

	if timeoutErr, isTimeout := Cast[TimeoutError](err); isTimeout {
		return timeoutErr.Timeout()
	}
	return false
}

func IsTimeoutError(err error) bool {
	var te TimeoutError
	if errors.As(err, &te) {
		return te.Timeout()
	}
	return false
}
