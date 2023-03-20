//go:build windows
// +build windows

package go_dparm

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_win"
	"golang.org/x/sys/windows"
)

type WindowsDriveFactory struct {
	drivers []plat_win.WinDriver
}

func NewWindowsDriveFactory() *WindowsDriveFactory {
	factory := &WindowsDriveFactory{}
	factory.drivers = []plat_win.WinDriver{
		plat_win.NewNvmeWinDriver(),
		//windows.NewSamsungNvmeDriver(),
		plat_win.NewWindowsNvmeDriver(),
		plat_win.NewScsiDriver(),
		plat_win.NewAtaDriver(),
	}
	return factory
}

func NewSystemDriveFactory() DriveFactory {
	return NewWindowsDriveFactory()
}

func (f *WindowsDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	handle, err := plat_win.OpenDevice(path)
	if err != nil {
		return nil, err
	}

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByHandle(handle)
		if err == nil {
			return driverHandle, nil
		}
	}

	_ = windows.CloseHandle(handle)

	return nil, errors.New("Not supported device")
}

func (f *WindowsDriveFactory) EnumDrives() ([]EnumDriveItem, error) {
	return nil, errors.New("Not supoported yet")
}

func (f *WindowsDriveFactory) EnumVolumes() ([]EnumVolumeItem, error) {
	return nil, errors.New("Not supoported yet")
}
