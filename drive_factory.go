package go_dparm

import "github.com/jc-lab/go-dparm/common"

type EnumVolumeContext interface {
	GetList() []common.VolumeInfo
	FindVolumesByDrive(driveInfo *common.DriveInfo) []common.VolumeInfo
}

type DriveFactory interface {
	OpenByPath(path string) (DriveHandle, error)
	EnumDrives() ([]common.DriveInfo, error)
	EnumVolumes() (EnumVolumeContext, error)
}
