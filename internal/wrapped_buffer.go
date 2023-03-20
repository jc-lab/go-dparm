package internal

import (
	"io"
)

type WrappedBuffer struct {
	io.Reader

	buffer   []byte
	capacity int

	refBuf    []byte
	limit     int
	writerPos int
	readerPos int
}

func NewWrappedBuffer(buffer []byte) *WrappedBuffer {
	b := &WrappedBuffer{
		buffer:    buffer,
		capacity:  len(buffer),
		limit:     len(buffer),
		writerPos: 0,
		readerPos: 0,
		refBuf:    buffer,
	}

	return b
}

func (b *WrappedBuffer) GetCapacity() int {
	return b.capacity
}

func (b *WrappedBuffer) SetLimit(limit int) {
	b.limit = limit
	b.refBuf = b.buffer[:b.limit]
}

func (b *WrappedBuffer) ResetWrite() {
	b.writerPos = 0
}

func (b *WrappedBuffer) ResetRead() {
	b.readerPos = 0
}

func (b *WrappedBuffer) Write(p []byte) (n int, err error) {
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

func (b *WrappedBuffer) Read(p []byte) (n int, err error) {
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
