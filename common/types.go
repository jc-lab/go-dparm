package common

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/nvme"
)

type PartitionStyle int

const (
	PartitionStyleRaw PartitionStyle = 0 + iota
	PartitionStyleMbr
	PartitionStyleGpt
)

type VolumeInfo struct {
	Path        string
	Filesystem  string
	MountPoints []string
}

type DriveInfo struct {
	DevicePath  string
	DrivingType DrivingType
	DriverName  string

	PartitionStyle   PartitionStyle
	GptDiskId        string // uuid format
	MbrDiskSignature uint32

	Model            string
	Serial           string
	FirmwareRevision string
	RawSerial        [20]byte

	WindowsDevNum int
	SmartEnabled  bool
	AtaIdentity   *ata.IdentityDeviceData
	NvmeIdentity  *nvme.IdentifyController

	IsSsd          bool
	SsdCheckWeight int
	TotalCapacity  int64
}
