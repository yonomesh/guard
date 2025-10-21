//go:build !windows

package buf

import "golang.org/x/sys/unix"

// 该方法的作用是将 Buffer 类型的缓冲区数据包装成一个 unix.Iovec 结构体，这个结构体
// 可以被用于 Unix 系统上的低级 I/O 操作。通过设置缓冲区的起始地址和长度，
// 可以在系统调用中传递给操作系统进行 I/O 操作。
func (b *Buffer) Iovec(length int) unix.Iovec {
	var iov unix.Iovec
	iov.Base = &b.data[b.start]
	iov.SetLen(length)
	return iov
}
