package proxy

import "errors"

const (
	maxInt          = int(^uint(0) >> 1) // maxInt is max buffer length limit
	smallBufferSize = 128                // smallBufferSize is an initial allocation minimal capacity.
)

var BufferTooLarge = errors.New("wasm plugin Buffer: too large")

type ConfigMap Header

type Header interface {
	// Get value of key
	// If multiple values associated with this key, first one will be returned.
	Get(key string) (string, bool)

	// Set key-value pair in header map, the previous pair will be replaced if exists
	Set(key, value string)

	// Add value for given key.
	// Multiple headers with the same key may be added with this function.
	// Use Set for setting a single header for the given key.
	Add(key, value string)

	// Del delete pair of specified key
	Del(key string)

	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(key, value string) bool)
}

type Buffer interface {
	Bytes() []byte
	Len() int
	Cap() int
	Grow(n int)
	Reset()
	Peek(n int) []byte
	Drain(len int)
	Mark()
	ResetMark()

	WriteByte(value byte) error
	WriteUint16(value uint16) error
	WriteUint32(value uint32) error
	WriteUint(value uint) error
	WriteUint64(value uint64) error
	WriteInt16(value int16) error
	WriteInt32(value int32) error
	WriteInt(value int) error
	WriteInt64(value int64) error

	PutByte(index int, value byte) error
	PutUint16(index int, value uint16) error
	PutUint32(index int, value uint32) error
	PutUint(index int, value uint) error
	PutUint64(index int, value uint64) error
	PutInt16(index int, value int16) error
	PutInt32(index int, value int32) error
	PutInt(index int, value int) error
	PutInt64(index int, value int64) error

	Write(p []byte) (n int, err error)

	ReadByte() (byte, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint() (uint, error)
	ReadUint64() (uint64, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt() (int, error)
	ReadInt64() (int64, error)

	GetByte(index int) (byte, error)
	GetUint16(index int) (uint16, error)
	GetUint32(index int) (uint32, error)
	GetUint(index int) (uint, error)
	GetUint64(index int) (uint64, error)
	GetInt16(index int) (int16, error)
	GetInt32(index int) (int32, error)
	GetInt(index int) (int, error)
	GetInt64(index int) (int64, error)
}

//func NewBuffer(size int) Buffer {
//	//return GetIoBuffer(size)
//}

type byteBuffer struct {
	buf  []byte // contents are the bytes buf[pos : len(buf)]
	pos  int    // read at &buf[pos], write at &buf[len(buf)]
	mark int
}

func (b *byteBuffer) Bytes() []byte {
	return b.buf[b.pos:]
}

func (b *byteBuffer) Len() int {
	return len(b.buf) - b.pos
}

func (b *byteBuffer) Cap() int {
	return cap(b.buf)
}

// Grow grows the buffer to guarantee space for n more bytes.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *byteBuffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

func (b *byteBuffer) Reset() {
	b.buf = b.buf[:0]
	b.pos = 0
}

func (b *byteBuffer) Peek(n int) []byte {

}

func (b *byteBuffer) Drain(len int) {

}

func (b *byteBuffer) Mark() {

}

func (b *byteBuffer) ResetMark() {

}

func (b *byteBuffer) WriteByte(value byte) error {

}

func (b *byteBuffer) WriteUint16(value uint16) error {

}

func (b *byteBuffer) WriteUint32(value uint32) error {

}

func (b *byteBuffer) WriteUint(value uint) error {

}

func (b *byteBuffer) WriteUint64(value uint64) error {

}

func (b *byteBuffer) WriteInt16(value int16) error {

}

func (b *byteBuffer) WriteInt32(value int32) error {

}

func (b *byteBuffer) WriteInt(value int) error {

}

func (b *byteBuffer) WriteInt64(value int64) error {

}

func (b *byteBuffer) PutByte(index int, value byte) error {

}

func (b *byteBuffer) PutUint16(index int, value uint16) error {

}

func (b *byteBuffer) PutUint32(index int, value uint32) error {

}

func (b *byteBuffer) PutUint(index int, value uint) error {

}

func (b *byteBuffer) PutUint64(index int, value uint64) error {

}

func (b *byteBuffer) PutInt16(index int, value int16) error {

}

func (b *byteBuffer) PutInt32(index int, value int32) error {

}

func (b *byteBuffer) PutInt(index int, value int) error {

}

func (b *byteBuffer) PutInt64(index int, value int64) error {

}

func (b *byteBuffer) Write(p []byte) (n int, err error) {

}

func (b *byteBuffer) ReadByte() (byte, error) {

}

func (b *byteBuffer) ReadUint16() (uint16, error) {

}

func (b *byteBuffer) ReadUint32() (uint32, error) {

}

func (b *byteBuffer) ReadUint() (uint, error) {

}

func (b *byteBuffer) ReadUint64() (uint64, error) {

}

func (b *byteBuffer) ReadInt16() (int16, error) {

}

func (b *byteBuffer) ReadInt32() (int32, error) {

}

func (b *byteBuffer) ReadInt() (int, error) {

}

func (b *byteBuffer) ReadInt64() (int64, error) {

}

func (b *byteBuffer) GetByte(index int) (byte, error) {

}

func (b *byteBuffer) GetUint16(index int) (uint16, error) {

}

func (b *byteBuffer) GetUint32(index int) (uint32, error) {

}

func (b *byteBuffer) GetUint(index int) (uint, error) {

}

func (b *byteBuffer) GetUint64(index int) (uint64, error) {

}

func (b *byteBuffer) GetInt16(index int) (int16, error) {

}

func (b *byteBuffer) GetInt32(index int) (int32, error) {

}

func (b *byteBuffer) GetInt(index int) (int, error) {

}

func (b *byteBuffer) GetInt64(index int) (int64, error) {

}

// ======================== private method impl ========================

// empty reports whether the unread portion of the buffer is empty.
func (b *byteBuffer) empty() bool { return b.pos >= len(b.buf) }

func (b *byteBuffer) grow(n int) int {
	m := b.Len()

	// If buffer is empty, reset to recover space.
	if m == 0 && b.pos != 0 {
		b.Reset()
	}

	// Try to grow by means of a re-slice.
	if i, ok := b.tryGrowBySlice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}

	c := cap(b.buf)
	if n <= c/2-m {
		// The current position can be moved.
		if b.pos > 0 {
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= c to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(b.buf, b.buf[b.pos:])
		}
	} else if c > maxInt-c-n {
		panic(BufferTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := make([]byte, m+n, 2*c+n)
		copy(buf, b.buf[b.pos:])
		b.buf = buf
	}
	// Restore b.pos and len(b.buf).
	b.pos = 0
	b.buf = b.buf[:m+n]
	return m
}

func (b *byteBuffer) tryGrowBySlice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}
