package common

import (
	"io"
	"net"

	"uni/bridge/tools"
)

type WithUpstream interface {
	Upstream() any
}

type stdWithUpstreamNetConn interface {
	NetConn() net.Conn
}

// Cast 尝试将 obj 转换为指定类型 T。如果转换成功，返回转换后的值和 true
func Cast[T any](obj any) (T, bool) {
	if c, ok := obj.(T); ok {
		return c, true
	}
	if u, ok := obj.(WithUpstream); ok {
		// obj.(WithUpstream) 是一个 类型断言，它检查 obj 是否实现了 WithUpstream 接口。
		// 2. 调用 u.Upstream 返回一个 any 值
		// 3. Cast[T](x) 又做一次 Cast
		return Cast[T](u.Upstream())
	}
	if u, ok := obj.(stdWithUpstreamNetConn); ok {
		return Cast[T](u.NetConn())
	}
	return tools.DefaultValue[T](), false
}

func MustCast[T any](obj any) T {
	value, ok := Cast[T](obj)
	if !ok {
		// make panic
		return obj.(T)
	}
	return value
}

func Top(obj any) any {
	if u, ok := obj.(WithUpstream); ok {
		return Top(u.Upstream())
	}
	if u, ok := obj.(stdWithUpstreamNetConn); ok {
		return Top(u.NetConn())
	}
	return obj
}

func Close(closers ...any) error {
	var retErr error
	for _, closer := range closers {
		if closer == nil {
			continue
		}
		switch c := closer.(type) {
		case io.Closer:
			err := c.Close()
			if err != nil {
				retErr = err
			}
			continue
		case WithUpstream:
			err := Close(c.Upstream())
			if err != nil {
				retErr = err
			}
		}
	}
	return retErr
}
