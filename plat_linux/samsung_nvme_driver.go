//go:build linux
// +build linux

package plat_linux

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
	_ "golang.org/x/sys/unix"
)

type SamsungNvmeDriver struct {
	LinuxDriver
}

func NewSamsungNvmeDriver() *SamsungNvmeDriver {
	return &SamsungNvmeDriver{}
}

func (d *SamsungNvmeDriver) OpenByHandle(handle int) (common.DriverHandle, error) {
	return nil, errors.New("not supported yet")
}
