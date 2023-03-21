//go:build linux
// +build linux

package go_dparm

import (
	"errors"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_linux"
)

type LinuxDriveFactory struct {
	drivers []plat_linux.LinuxDriver
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}

	return factory
}

func NewSystemDriveFactory() common.DriveFactory {
	return NewLinuxDriveFactory()
}

func (f *LinuxDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	handle, err := plat_linux.OpenDevice(path)
	if err != nil {
		return nil, err
	}

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByHandle(handle)
		if err != nil {
			return driverHandle, nil 
		}
	}

	_ = unix.Close(handle)

	return nil, errors.New("not supported device")
}

func (f *LinuxDriveFactory) EnumDrives() ([]EnumDriveItem, error) {
	return nil, nil // not implmented
}

func (f *LinuxDriveFactory) EnumVolumes() ([]EnumVolumeItem, error) {
	return nil, nil
}
