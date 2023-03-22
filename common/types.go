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

	WindowsDevNum   int
	SmartEnabled    bool
	AtaIdentity     *ata.IdentityDeviceData
	AtaIdentityRaw  []byte
	NvmeIdentity    *nvme.IdentifyController
	NvmeIdentityRaw []byte

	IsSsd          bool
	SsdCheckWeight int
	TotalCapacity  int64

	TcgSupport           int
	TcgTper              bool
	TcgLocking           bool
	TcgGeometryReporting bool
	TcgOpalSscV100       bool
	TcgOpalSscV200       bool
	TcgEnterprise        bool
	TcgSingleUser        bool
	TcgDataStore         bool

	TcgRawFeatures map[uint16][]byte
}
