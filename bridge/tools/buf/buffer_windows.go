package buf

import "golang.org/x/sys/windows"

// windows.WSABuf 和 unix.Iovec 实际上是涉及 IO 操作的缓冲区 包含缓冲区和指向缓冲区的指针
// Iovec 是将一个 Buffer 转换成 Unix 或 Windows 上的和 IO 相关的 Buf 结构体
func (b *Buffer) Iovec(length int) windows.WSABuf {
	return windows.WSABuf{
		Buf: &b.data[b.start],
		Len: uint32(length),
	}
}
