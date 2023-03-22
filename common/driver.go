package common

import (
	"github.com/jc-lab/go-dparm/ata"
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

	SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error
}

type AtaDriverHandle interface {
	GetIdentity() []byte
	DoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error
}

type NvmeDriverHandle interface {
	GetIdentity() []byte
	NvmeGetLogPage(nsid uint32, logId uint32, rae bool, size int) ([]byte, error)
}
