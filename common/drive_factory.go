package common

import (
	"github.com/jc-lab/go-dparm/ata"
)

type DriveHandle interface {
	GetDriverHandle() DriverHandle
	Close()
	GetDevicePath() string
	GetDrivingType() DrivingType
	GetDriverName() string

	GetDriveInfo() *DriveInfo

	// ATA
	AtaDoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error

	// NVME
	NvmeGetLogPage(nsid uint32, logId uint32, rae bool, size int) ([]byte, error)

	// COMMON
	SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error

	TcgDiscovery0() error
}

type EnumVolumeContext interface {
	GetList() []VolumeInfo
	FindVolumesByDrive(driveInfo *DriveInfo) []VolumeInfo
	OpenDriveByVolumePath(volumePath string) (DriveHandle, error)
	OpenDriveByPartition(partition Partition) (DriveHandle, error)
}

type DriveFactory interface {
	OpenByPath(path string) (DriveHandle, error)
	EnumDrives() ([]DriveInfo, error)
	EnumVolumes() (EnumVolumeContext, error)
}
