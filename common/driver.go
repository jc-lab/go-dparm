package common

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/nvme"
)

//TODO: DRIVER interface

type WindowsPhysicalDrive struct {
	DeviceIndex        int
	PhysicalDiskPath   string
	SetupApiDevicePath string
}

type DrivingType int

const (
	DrivingUnknown DrivingType = 0 + iota
	DrivingAtapi
	DrivingNvme
)

type DriverHandle interface {
	GetDriverName() string
	GetDrivingType() DrivingType
	ReopenWritable() error
	Close()
}

type AtaDriverHandle interface {
	GetIdentity() *ata.IdentityDeviceData
	DoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error
}

type NvmeDriverHandle interface {
	GetIdentity() *nvme.IdentifyController
}
