package buf

import (
	"crypto/rand"
	"io"
	"net"
	"sync/atomic"

	E "guard/bridge/common/errors"
	"guard/bridge/tools"
	F "guard/bridge/tools/format"
)

type Buffer struct {
	data     []byte
	start    int
	end      int
	capacity int
	refs     atomic.Int32
	managed  bool
}

func New() *Buffer {
	return &Buffer{
		data:     Get(BufferSize),
		capacity: BufferSize,
		managed:  true,
	}
}

func NewPacket() *Buffer {
	return &Buffer{
		data:     Get(UDPBufferSize),
		capacity: UDPBufferSize,
		managed:  true,
	}
}

func NewSize(size int) *Buffer {
	if size == 0 {
		return &Buffer{}
	} else if size > 65535 { // 2^16=1024*64-1
		return &Buffer{
			data:     make([]byte, size),
			capacity: size,
		}
	}

	return &Buffer{
		data:     Get(size),
		capacity: size,
		managed:  true,
	}
}

func As(data []byte) *Buffer {
	return &Buffer{
		data:     data,
		end:      len(data),
		capacity: len(data),
	}
}

func With(data []byte) *Buffer {
	return &Buffer{
		data:     data,
		capacity: len(data),
	}
}

func (b *Buffer) IsFull() bool {
	return b.end-b.start == b.capacity
}

func (b *Buffer) IsEmpty() bool {
	return b.end-b.start == 0
}

// FreeBytes
func (b *Buffer) AvailableSpace() []byte {
	return b.data[b.end:b.capacity]
}

func (b *Buffer) Byte(index int) byte {
	return b.data[b.start+index]
}

func (b *Buffer) SetByte(index int, value byte) {
	b.data[b.start+index] = value
}

func (b *Buffer) Bytes() []byte {
	return b.data[b.start:b.end]
}

// PeekN returns a slice of the first n available bytes in the buffer.
// fun (b *Buffer)To
func (b *Buffer) PeekN(n int) []byte {
	return b.data[b.start : b.start+n]
}

// PeekAfter returns a slice of the bytes starting after the first n available bytes.
// fun (b *Buffer)From
func (b *Buffer) PeekAfterN(n int) []byte {
	return b.data[b.start+n : b.end]
}

// Extend reserves 'n' bytes at the end of the buffer and returns the writable slice.
// Panics if capacity is insufficient. This is a zero-copy operation.
func (b *Buffer) Extend(n int) []byte {
	end := b.end + n
	if end > b.capacity {
		panic(F.ToString("buffer overflow: capacity ", b.capacity, ",end ", b.end, ", need ", n))
	}
	ext := b.data[b.end:end]
	b.end = end
	return ext
}

// func Advance
func (b *Buffer) Shift(from int) {
	b.start += from
	if b.end < b.start {
		b.end = b.start
	}
}

func (b *Buffer) Resize(start, length int) {
	b.start = start
	b.end = b.start + length
}

// 剩余容量是否足够插入 n byte
func (b *Buffer) HasSpace(n int) bool {
	if b.end-b.start+n > b.capacity {
		return false
		// panic(F.ToString("buffer overflow: capacity is ", b.capacity, ", but need ", n))
	}
	return true
}

// Deprecated: The Reset modifies the capacity, which is not recommended.
func (b *Buffer) Reset() {
	b.start = 0
	b.end = 0
	b.capacity = len(b.data)
}

// Increment Reference Count
func (b *Buffer) IncRef() {
	b.refs.Add(1)
}

// Decrement Reference Count
func (b *Buffer) DecRef() {
	b.refs.Add(-1)
}

// Maybe it can be renamed to Recycle()
func (b *Buffer) Release() {
	if b == nil || !b.managed {
		return
	}
	if b.refs.Load() > 0 {
		return
	}
	tools.Must(Put(b.data))
	*b = Buffer{}
}

func (b *Buffer) Truncate(to int) {
	b.end = b.start + to
}

func (b *Buffer) ExtendHeader(n int) []byte {
	if b.start < n {
		panic(F.ToString("buffer overflow: capacity ", b.capacity, ",start", b.start, ", need", n))
	}
	b.start -= n
	return b.data[b.start : b.start+n]
}

// n 写入多少 byte
func (b *Buffer) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return
	}

	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}

	n = copy(b.data[b.end:b.capacity], data)
	b.end += n
	return
}

// WriteTo writes all currently available data in the buffer to the writer w.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.Bytes())
	return int64(n), err
}

func (b *Buffer) WriteRandom(size int) []byte {
	buf := b.Extend(size)
	tools.Must1(io.ReadFull(rand.Reader, buf))
	return buf
}

func (b *Buffer) WriteByte(d byte) error {
	if b.IsFull() {
		return io.ErrShortBuffer
	}
	b.data[b.end] = d
	b.end++
	return nil
}

func (b *Buffer) WriteRune(s rune) (int, error) {
	return b.Write([]byte{byte(s)})
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return
	}
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}

	n = copy(b.data[b.end:b.capacity], s)
	b.end += n
	return
}

func (b *Buffer) WriteZero() error {
	if b.IsFull() {
		return io.ErrShortBuffer
	}
	b.data[b.end] = 0
	b.end++
	return nil
}

func (b *Buffer) WriteZeroN(n int) error {
	if b.end+n > b.capacity {
		return io.ErrShortBuffer
	}
	tools.ClearSlice(b.Extend(n))
	return nil
}

// Read reads up to len(p) bytes from the Buffer into p.
func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.IsEmpty() {
		return 0, io.EOF
	}

	n = copy(p, b.data[b.start:b.end])
	b.start += n
	return
}

// read data from an io.Reader directly into the available free space of the Buffer
func (b *Buffer) ReadOnceFrom(r io.Reader) (int, error) {
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n, err := r.Read(b.AvailableSpace())
	b.end += n
	return n, err
}

// net.PacketConn which is typically used for connectionless protocols like UDP.
func (b *Buffer) ReadPacketFrom(pc net.PacketConn) (int64, net.Addr, error) {
	if b.IsFull() {
		return 0, nil, io.ErrShortBuffer
	}
	n, addr, err := pc.ReadFrom(b.AvailableSpace())
	b.end += n
	return int64(n), addr, err
}

func (b *Buffer) ReadAtLeastFrom(r io.Reader, min int) (int64, error) {
	if min <= 0 {
		n, err := b.ReadOnceFrom(r)
		return int64(n), err
	}
	if b.IsFull() {
		return 0, io.ErrShortBuffer
	}
	n, err := io.ReadAtLeast(r, b.AvailableSpace(), min)
	b.end += n
	return int64(n), err
}

func (b *Buffer) ReadFullFrom(r io.Reader, size int) (n int, err error) {
	if b.end+size > b.capacity {
		return 0, io.ErrShortBuffer
	}
	n, err = io.ReadFull(r, b.data[b.end:b.end+size])
	b.end += n
	return
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		if b.IsFull() {
			return 0, io.ErrShortBuffer
		}
		var readN int
		readN, err = r.Read(b.AvailableSpace())
		b.end += readN
		n += int64(readN)
		if err != nil {
			if E.IsMulti(err, io.EOF) {
				err = nil
			}
			return
		}
	}
}

// Buffer 类型如果使用 ReadByte，就可以自动满足 io.ByteReader 接口
func (b *Buffer) ReadByte() (byte, error) {
	if b.IsEmpty() {
		return 0, io.EOF
	}
	nextByte := b.data[b.start]
	b.start++
	return nextByte, nil
}

func (b *Buffer) ReadBytes(n int) ([]byte, error) {
	if b.end-b.start < n {
		return nil, io.ErrUnexpectedEOF
	}

	data := b.data[b.start : b.start+n]
	b.start += n
	return data, nil
}
