//go:build !windows
// +build !windows

package go_dparm

import (
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_linux"
)

func NewSystemDriveFactory() common.DriveFactory {
	return plat_linux.NewLinuxDriveFactory()
}
