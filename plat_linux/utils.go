//go:build linux
// +build linux

package plat_linux

import (
	"log"
	"strings"
	"syscall"
	"unsafe"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/partition"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/diskfs/go-diskfs/partition/mbr"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal/direct_mbr"

	"golang.org/x/sys/unix"
)

type LinuxBasicInfo struct {
	PartitionStyle  common.PartitionStyle
	PartitionTable  partition.Table
	DiskGeometry    unix.HDGeometry
	MbrSignature    uint32
	GptDiskId       string
	BlockTotalBytes int64
	BlockSectorSize int
}

func OpenDevice(path string) (int, error) {
	return unix.Open(path, unix.O_RDWR|unix.O_NONBLOCK,
		uint32(unix.S_IRUSR|unix.S_IWUSR|unix.S_IRGRP|unix.S_IWGRP|unix.S_IROTH|unix.S_IWOTH))
}

func ReadBasicInfo(fd int, path string) (*LinuxBasicInfo, error) {
	result := &LinuxBasicInfo{}

	dev, err := diskfs.Open(path, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return nil, common.NewNestedError(path+" diskfs.open failed", err)
	}

	// Get sector size
	var sectorSize int
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), unix.BLKSSZGET, uintptr(unsafe.Pointer(&sectorSize))); errno == 0 {
		result.BlockSectorSize = sectorSize
	}

	// Get device size in bytes
	var size uint64
	if err := ioctl(fd, unix.BLKGETSIZE64, uintptr(unsafe.Pointer(&size))); err == nil {
		result.BlockTotalBytes = int64(size)
	}

	if err = ioctl(fd, unix.HDIO_GETGEO, uintptr(unsafe.Pointer(&result.DiskGeometry))); err != nil {
		return nil, common.NewNestedError(path+" HDIO_GETGEO failed", err)
	}

	result.PartitionTable = dev.Table

	switch pt := dev.Table.(type) {
	case *gpt.Table:
		result.PartitionStyle = common.PartitionStyleGpt
		result.GptDiskId = strings.ToLower(pt.GUID)
	case *mbr.Table:
		result.PartitionStyle = common.PartitionStyleMbr
		tableEx, err := direct_mbr.Read(dev.File, int(dev.LogicalBlocksize), int(dev.PhysicalBlocksize))
		if err == nil {
			result.MbrSignature = tableEx.MbrIdentifier
		}
	default:
		log.Printf("%s: %T\n", "Unknown partition type", pt)
	}

	return result, nil
}

func ioctl(fd int, op, arg uintptr) error {
	ret, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		op,
		uintptr(unsafe.Pointer(&arg)),
	)

	if err != 0 {
		return err
	} else if ret != 0 {
		return unix.Errno(ret)
	}

	return nil
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
