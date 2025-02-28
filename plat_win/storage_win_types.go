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
	_                    uint32
}

type PARTITION_INFORMATION_MBR struct {
	PartitionType       byte
	BootIndicator       bool
	RecognizedPartition bool
	HiddenSectors       uint32
	PartitionId         windows.GUID
}

type PARTITION_INFORMATION_GPT struct {
	PartitionType windows.GUID
	PartitionId   windows.GUID
	Attributes    uint64
	Name          [36]uint16
}

func (p *PARTITION_INFORMATION_EX) GetMbr() *PARTITION_INFORMATION_MBR {
	if p.PartitionStyle == PartitionStyleMbr {
		return (*PARTITION_INFORMATION_MBR)(unsafe.Pointer(&p.PartitionInfo[0]))
	}
	return nil
}

func (p *PARTITION_INFORMATION_EX) GetGpt() *PARTITION_INFORMATION_GPT {
	if p.PartitionStyle == PartitionStyleGpt {
		return (*PARTITION_INFORMATION_GPT)(unsafe.Pointer(&p.PartitionInfo[0]))
	}
	return nil
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
