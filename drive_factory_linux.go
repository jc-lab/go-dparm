//go:build linux
// +build linux

package go_dparm

import "github.com/jc-lab/go-dparm/common"

type LinuxDriveFactory struct {
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}

	return factory
}

func NewSystemDriveFactory() common.DriveFactory {
	return NewLinuxDriveFactory()
}

func (l LinuxDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinuxDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinuxDriveFactory) EnumVolumes() (common.EnumVolumeContext, error) {
	//TODO implement me
	panic("implement me")
}
