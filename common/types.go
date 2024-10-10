package common

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/nvme"
	"github.com/jc-lab/go-dparm/tcg"
)

type PartitionStyle int

const (
	PartitionStyleRaw PartitionStyle = 0 + iota
	PartitionStyleMbr
	PartitionStyleGpt
)

type MbrPartitionInfo struct {
	PartitionType byte
	BootIndicator bool
}

type GptPartitionInfo struct {
	// PartitionType GUID Upper String Format {...}
	PartitionType string
	// PartitionId GUID Upper String Format {...}
	PartitionId string
}

type Partition interface {
	// GetStart bytes
	GetStart() uint64
	// GetEnd bytes
	GetSize() uint64

	GetPartitionStyle() PartitionStyle

	GetMbrInfo() *MbrPartitionInfo
	GetGptInfo() *GptPartitionInfo
}

type DriveType int

const (
	DriveTypeUnknown DriveType = 0 + iota
	DriveTypeNoRootDir
	DriveTypeRemovable
	DriveTypeFixed
	DriveTypeRemote
	DriveTypeCdrom
	DriveTypeRamdisk
	DriveTypeTape
)

type VolumeInfo struct {
	Path        string
	Filesystem  string
	MountPoints []string
	Partitions  []Partition
	DriveType   DriveType
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

	VendorId        string
	ProductRevision string

	WindowsDevNum   int
	SmartEnabled    bool
	AtaIdentity     *ata.IdentityDeviceData
	AtaIdentityRaw  []byte
	NvmeIdentity    *nvme.IdentifyController
	NvmeIdentityRaw []byte

	IsSsd          bool
	SsdCheckWeight int
	TotalCapacity  int64

	TcgSupport int
	tcg.TcgLevel0Info

	Removable  int
	DriveType  DriveType
	Partitions []Partition
}
