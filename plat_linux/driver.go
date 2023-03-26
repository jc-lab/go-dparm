//go:build linux
// +build linux

package plat_linux

import (
	"github.com/jc-lab/go-dparm/common"
)

type LinuxDriver interface {
	OpenByFd(fd int) (common.DriverHandle, error)
}