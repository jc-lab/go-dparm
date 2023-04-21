//go:build windows
// +build windows

package plat_win

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type DRIVE_LAYOUT_INFORMATION_GPT struct {
	DiskId               windows.GUID
	StartingUsableOffset uint64
	UsableLength         uint64
	MaxPartitionCount    uint32
}

func GetSizeOf_DRIVE_LAYOUT_INFORMATION() int {
	a := unsafe.Sizeof(DRIVE_LAYOUT_INFORMATION_GPT{})
	b := unsafe.Sizeof(DRIVE_LAYOUT_INFORMATION_MBR{})
	if a > b {
		return int(a)
	} else {
		return int(b)
	}
}
