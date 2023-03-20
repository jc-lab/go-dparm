//go:build windows
// +build windows

package windows

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
)

type SamsungNvmeDriver struct {
	common.Driver
}

func NewSamsungNvmeDriver() *SamsungNvmeDriver {
	return &SamsungNvmeDriver{}
}

func (d *SamsungNvmeDriver) OpenByPath(path string) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}

func (d *SamsungNvmeDriver) OpenByWindowsPhysicalDrive(path *common.WindowsPhysicalDrive) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}
