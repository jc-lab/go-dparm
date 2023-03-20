//go:build windows
// +build windows

package plat_win

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
)

type SamsungNvmeDriver struct {
	WinDriver
}

func NewSamsungNvmeDriver() *SamsungNvmeDriver {
	return &SamsungNvmeDriver{}
}

func (d *SamsungNvmeDriver) OpenByHandle(handle windows.Handle) (common.DriveHandle, error) {
	return nil, errors.New("Not supported yet")
}
