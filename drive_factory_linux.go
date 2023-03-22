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

func NewSystemDriveFactory() DriveFactory {
	return NewLinuxDriveFactory()
}

func (l LinuxDriveFactory) OpenByPath(path string) (DriveHandle, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinuxDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinuxDriveFactory) EnumVolumes() (EnumVolumeContext, error) {
	//TODO implement me
	panic("implement me")
}
