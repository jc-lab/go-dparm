package common

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/nvme"
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
