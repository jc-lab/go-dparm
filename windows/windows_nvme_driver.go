//go:build windows
// +build windows

package windows

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
)

type WindowsNvmeDriver struct {
	common.Driver
}

func NewWindowsNvmeDriver() *WindowsNvmeDriver {
	return &WindowsNvmeDriver{}
}

func (d *WindowsNvmeDriver) OpenByPath(path string) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}

func (d *WindowsNvmeDriver) OpenByWindowsPhysicalDrive(path *common.WindowsPhysicalDrive) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}
