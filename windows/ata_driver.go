//go:build windows
// +build windows

package windows

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
)

type AtaDriver struct {
	common.Driver
}

func NewAtaDriver() *AtaDriver {
	return &AtaDriver{}
}

func (d *AtaDriver) OpenByPath(path string) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}

func (d *AtaDriver) OpenByWindowsPhysicalDrive(path *common.WindowsPhysicalDrive) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}
