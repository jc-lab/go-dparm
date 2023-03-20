//go:build windows
// +build windows

package windows

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
)

type NvmeWinDriver struct {
	common.Driver
}

func NewNvmeWinDriver() *NvmeWinDriver {
	return &NvmeWinDriver{}
}

func (d *NvmeWinDriver) OpenByPath(path string) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}

func (d *NvmeWinDriver) OpenByWindowsPhysicalDrive(path *common.WindowsPhysicalDrive) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}
