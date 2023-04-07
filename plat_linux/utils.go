//go:build linux
// +build linux

package plat_linux

import (
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/unix"
)

const (
	
)

type LinuxBasicInfo struct {
	StorageDeviceNumber *STORAGE_DEVICE_NUMBER
	DiskGeometryEx      *DISK_GEOMETRY_EX
	PartitionStyle      common.PartitionStyle
	MbrSignature        uint32
	GptDiskId           string
}

func OpenDevice(path string) (int, error) {
	return unix.Open(path, unix.O_RDWR | unix.O_NONBLOCK, uint32(unix.S_IRUSR | unix.S_IWUSR | unix.S_IRGRP | unix.S_IWGRP))
}

func readNullTerminatedAscii(buf []byte, offset int) string {
	if offset <= 0 {
		return ""
	}
	buf = buf[offset:]
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return string(buf[:i])
		}
	}
	return ""
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

func zerofill(buf []uint16) {
	for i := range buf {
		buf[i] = 0
	}
}

func wcslen(buf []uint16) int {
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return i
		}
	}
	return len(buf)
}
