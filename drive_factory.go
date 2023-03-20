package go_dparm

import "github.com/jc-lab/go-dparm/common"

type EnumDriveItem interface {
}

type EnumVolumeItem interface {
}

type DriveFactory interface {
	OpenByPath(path string) (common.DriveHandle, error)
	EnumDrives() ([]EnumDriveItem, error)
	EnumVolumes() ([]EnumVolumeItem, error)
}
