//go:build linux
// +build linux

package plat_linux

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func OpenDevice(path string) (int, error) {
	return unix.Open(path, unix.O_RDWR | unix.O_NONBLOCK, uint32(unix.S_IRUSR | unix.S_IWUSR | unix.S_IRGRP | unix.S_IWGRP))
}

func copyToPointer(dest unsafe.Pointer, src []byte, len int) {
	destRef := unsafe.Slice((*byte)(dest), len)
	copy(destRef, src[:len])
}

func copyFromPointer(dest []byte, src unsafe.Pointer, len int) {
	srcRef := unsafe.Slice((*byte)(src), len)
	copy(dest, srcRef[:len])
}

func copyFromAsciiToBuffer(dest []byte, text string) {
	c := len(text)
	for i := 0; i < c; i++ {
		dest[i] = text[i]
	}
}