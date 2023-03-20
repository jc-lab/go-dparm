//go:build windows
// +build windows

package go_dparm

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
)

type WindowsDriveFactory struct {
	drivers []common.Driver
}

func NewWindowsDriveFactory() *WindowsDriveFactory {
	factory := &WindowsDriveFactory{}
	factory.drivers = []common.Driver{
		//windows.NewNvmeWinDriver(),
		//windows.NewSamsungNvmeDriver(),
		//windows.NewWindowsNvmeDriver(),
		//windows.NewScsiDriver(),
		//windows.NewAtaDriver(),
	}
	return factory
}

func NewSystemDriveFactory() DriveFactory {
	return NewWindowsDriveFactory()
}

func (f *WindowsDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	return nil, errors.New("Not supoported yet")
}

func (f *WindowsDriveFactory) EnumDrives() ([]EnumDriveItem, error) {
	return nil, errors.New("Not supoported yet")
}

func (f *WindowsDriveFactory) EnumVolumes() ([]EnumVolumeItem, error) {
	return nil, errors.New("Not supoported yet")
}
