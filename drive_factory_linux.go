//go:build linux
// +build linux

package go_dparm

import (
	"log"
	"strings"
	_ "unsafe"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_linux"

	"golang.org/x/sys/unix"
)

type LinuxDriveFactory struct {
	drivers []plat_linux.LinuxDriver
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}
	factory.drivers = []plat_linux.LinuxDriver{
		//plat_linux.NewNvmeLinuxDriver(),
		//linux.NewSamsungNvmeDriver
		plat_linux.NewLinuxNvmeDriver(),
		plat_linux.NewSgDriver(),
	}
	return factory
}

func NewSystemDriveFactory() common.DriveFactory {
	return NewLinuxDriveFactory()
}

func (f *LinuxDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	fd, err := plat_linux.OpenDevice(path)
	if err != nil {
		return nil, err
	}

	driveHandle, err := f.OpenByFd(fd, path)
	if err == nil {
		return driveHandle, nil
	}

	_ = unix.Close(fd)

	return nil, err
}

func (f *LinuxDriveFactory) OpenByFd(fd int, path string) (common.DriveHandle, error) {
	impl := &DriveHandleImpl{}
	impl.Info.DrivingType = common.DrivingUnknown
	impl.Info.DevicePath = path

	return impl, nil
}

func (f *LinuxDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	var results []common.DriveInfo

	var names []string
	var s unix.Stat_t

	dfd, err := unix.Open("/sys/block", unix.O_RDONLY | unix.O_DIRECTORY, 666)
	if err != nil {
		return nil, err
	}
	defer unix.Close(dfd)

	buf := make([]byte, 4096)
	entNum, err := unix.ReadDirent(dfd, buf)
	if err != nil {
		return nil, err
	}
	_, _, names = unix.ParseDirent(buf, entNum, names)

	for _, name := range names {
		devPath := "/dev/"
		devPath += name
		if (!strings.Contains(name, "loop")) && (unix.Stat(devPath, &s) == nil) {
			if ((s.Mode & unix.S_IFMT) == unix.S_IFBLK) {
				driveHandle, err := f.OpenByPath(devPath)
				if err != nil {
					log.Println(err)
					continue
				}
				results = append(results, *driveHandle.GetDriveInfo())
			}
		}
	}
	return results, nil
}

func (f *LinuxDriveFactory) EnumVolumes() (common.EnumVolumeContext, error) {
	return plat_linux.EnumVolumes(f)
}
