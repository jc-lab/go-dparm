package internal

import (
	"io"
	"unsafe"
)

type AlignedBuffer struct {
	io.Reader
	io.Writer

	buffer   []byte
	pointer  *byte
	refBuf   []byte
	capacity int

	limit     int
	writerPos int
	readerPos int
}

func IsAlignedPointer(align int, pointer uintptr) bool {
	return (pointer % uintptr(align)) > 0
}

func NewAlignedBuffer(align int, size int) *AlignedBuffer {
	allocateSize := align + size
	b := &AlignedBuffer{
		buffer:    make([]byte, allocateSize),
		capacity:  size,
		limit:     size,
		writerPos: 0,
		readerPos: 0,
	}

	pointer := uintptr(unsafe.Pointer(&b.buffer[0]))
	tmp := pointer % uintptr(align)
	if tmp > 0 {
		pointer += uintptr(align) - tmp
	}

	b.pointer = (*byte)(unsafe.Pointer(pointer))
	b.refBuf = unsafe.Slice(b.pointer, b.limit)

	return b
}

func (b *AlignedBuffer) GetPointer() *byte {
	return b.pointer
}

func (b *AlignedBuffer) GetBuffer() []byte {
	return b.refBuf
}

func (b *AlignedBuffer) GetCapacity() int {
	return b.capacity
}

func (b *AlignedBuffer) GetPos() int {
	return b.writerPos
}

func (b *AlignedBuffer) SetLimit(limit int) {
	b.limit = limit
	b.refBuf = unsafe.Slice(b.pointer, b.limit)
}

func (b *AlignedBuffer) Reset() {
	b.ResetRead()
	b.ResetWrite()

	// memclr
	for i := range b.refBuf {
		b.refBuf[i] = 0
	}
}

func (b *AlignedBuffer) ResetWrite() {
	b.writerPos = 0
}

func (b *AlignedBuffer) ResetRead() {
	b.readerPos = 0
}

func (b *AlignedBuffer) WriteByte(p byte) (err error) {
	if b.writerPos >= b.limit {
		return io.EOF
	}

	b.refBuf[b.writerPos] = p
	b.writerPos++

	return nil
}

func (b *AlignedBuffer) Write(p []byte) (n int, err error) {
	remaining := b.limit - b.writerPos
	available := remaining
	if available > len(p) {
		available = len(p)
	}
	if available == 0 {
		return 0, io.EOF
	}

	copy(b.refBuf[b.writerPos:b.writerPos+available], p)
	b.writerPos += available

	return available, nil
}

func (b *AlignedBuffer) Read(p []byte) (n int, err error) {
	remaining := b.limit - b.readerPos
	available := remaining
	if available > len(p) {
		available = len(p)
	}
	if available == 0 {
		return 0, io.EOF
	}

	copy(p, b.refBuf[b.readerPos:b.readerPos+available])
	b.readerPos += available

	return available, nil
}
