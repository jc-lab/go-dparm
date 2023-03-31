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
		plat_linux.NewNvmeLinuxDriver(),
		//linux.NewSamsungNvmeDriver
		plat_linux.NewLinuxNvmeDriver(),
		plat_linux.NewSgDriver(),
		plat_linux.NewAtaDriver(),
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

	basicInfo := plat_linux.ReadBasicInfo(fd)

	impl.Info.PartitionStyle = basicInfo.PartitionStyle
	impl.Info.GptDiskId = basicInfo.GptDiskId
	impl.Info.MbrDiskSignature = basicInfo.MbrSignature

	if basicInfo.StorageDeviceNumber != nil {
		impl.Info.LinuxDevNum = int(basicInfo.StorageDeviceNumber.DeviceNumber)
	}

	if basicInfo.DiskGeometryEx != nil {
		impl.Info.TotalCapacity = int64(basicInfo.DiskGeometryEx.DiskSize)
	}

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByFd(fd)
		if err == nil {
			impl.dh = driverHandle
			impl.Info.DrivingType = driverHandle.GetDrivingType()
			impl.Info.DriverName = driverHandle.GetDriverName()
			impl.init()
			return impl, nil
		}
	}

	storageQueryResp, err := plat_linux.ReadStorageQuery(fd)
	if err == nil && storageQueryResp != nil {
		if impl.Info.Model == "" {
			impl.Info.Model = storageQueryResp.ProductId
		}
		if impl.Info.Serial == "" {
			impl.Info.Serial = storageQueryResp.SerialNumber
		}

		impl.Info.VendorId = storageQueryResp.VendorId
		impl.Info.ProductRevision = storageQueryResp.ProductRevision
	}

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
	_, _, names = unix.ParseDirent(buf, entNum, names)

	for _, name := range names {
		devPath := "/dev/"
		devPath += name
		if (!strings.Contains(name, "loop")) && (unix.Stat(devPath, &s) == nil) {
			if (s.Mode & unix.S_IFMT == unix.S_IFBLK) {
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
