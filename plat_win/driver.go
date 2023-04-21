//go:build windows
// +build windows

package plat_win

import (
	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
)

type WinDriver interface {
	OpenByHandle(handle windows.Handle) (common.DriverHandle, error)
}
