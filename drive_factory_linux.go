//go:build linux
// +build linux

package go_dparm

type LinuxDriveFactory struct {
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}

	return factory
}

func NewSystemDriveFactory() DriveFactory {
	return NewLinuxDriveFactory()
}
