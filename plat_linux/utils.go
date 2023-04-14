//go:build linux
// +build linux

package plat_linux

import (
	"log"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/diskfs"
	"github.com/jc-lab/go-dparm/diskfs/partition/gpt"
	"github.com/jc-lab/go-dparm/diskfs/partition/mbr"

	"golang.org/x/sys/unix"
)

type LinuxBasicInfo struct {
	PartitionStyle	common.PartitionStyle      
	DiskGeometry 	unix.HDGeometry
	MbrSignature 	uint32
	GptDiskId 		string
}

func OpenDevice(path string) (int, error) {
	return unix.Open(path, unix.O_RDWR | unix.O_NONBLOCK, 
		uint32(unix.S_IRUSR | unix.S_IWUSR | unix.S_IRGRP | unix.S_IWGRP | unix.S_IROTH | unix.S_IWOTH))
}

func ReadBasicInfo(fd int, path string) *LinuxBasicInfo {
	result := &LinuxBasicInfo{}

	dev, err := diskfs.Open(path, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		log.Fatalln(err)
	}

	_, _, err = unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		unix.HDIO_GETGEO,
		uintptr(unsafe.Pointer(&result.DiskGeometry)),
	)
	if err != unix.Errno(0) {
		log.Fatalln(err)
	}

	switch pt := dev.Table.(type) {
	case *gpt.Table:
		result.PartitionStyle = common.PartitionStyleGpt
		result.GptDiskId = pt.GUID
	case *mbr.Table:
		result.PartitionStyle = common.PartitionStyleMbr
		result.MbrSignature = pt.MbrIdentifier
	default:
		log.Printf("%s: %T\n", "Unknown partition type", pt)
	}

	return result
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
