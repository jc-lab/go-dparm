package common

//TODO: DRIVER interface

type WindowsPhysicalDrive struct {
	DeviceIndex        int
	PhysicalDiskPath   string
	SetupApiDevicePath string
}

type DrivingType int

const (
	kDrivingUnknown DrivingType = 0 + iota
	kDrivingAtapi
	kDrivingNvme
)

type DriverHandle interface {
	GetDriverName() string
	MergeDriveInfo(data DriveInfo)
	GetDrivingType() DrivingType
	ReopenWritable() error
	Close()
}

type Driver interface {
	OpenByPath(path string) (DriverHandle, error)
	OpenByWindowsPhysicalDrive(path *WindowsPhysicalDrive) (DriverHandle, error)
}
