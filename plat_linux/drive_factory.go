//go:build linux
// +build linux

package plat_linux

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jc-lab/go-dparm/common"

	"golang.org/x/sys/unix"
)

var skipPat *regexp.Regexp = regexp.MustCompile(`nvme-\w*\.|-part|scsi|wwn`)

type LinuxDriveFactory struct {
	drivers []LinuxDriver
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}
	factory.drivers = []LinuxDriver{
		//linux.NewSamsungNvmeDriver
		NewLinuxNvmeDriver(),
		NewSgDriver(),
	}
	return factory
}

func NewSystemDriveFactory() common.DriveFactory {
	return NewLinuxDriveFactory()
}

func (f *LinuxDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	fd, err := OpenDevice(path)
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
	var err error

	impl := &common.DriveHandleImpl{}
	impl.Info.DrivingType = common.DrivingUnknown
	impl.Info.DevicePath = path
	impl.Info.Removable = -1

	basicInfo, err := ReadBasicInfo(fd, path)

	if err == nil {
		impl.Info.PartitionStyle = basicInfo.PartitionStyle
		impl.Info.GptDiskId = basicInfo.GptDiskId
		impl.Info.MbrDiskSignature = basicInfo.MbrSignature
	}

	// Try to get incomplete data first in case of inquiry failure..
	impl.Info.Model, impl.Info.Serial, impl.Info.VendorId, impl.Info.FirmwareRevision, err = getIdInfo(path)
	_ = err

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByFd(fd)
		if err == nil {
			impl.Dh = driverHandle
			impl.Info.DrivingType = driverHandle.GetDrivingType()
			impl.Info.DriverName = driverHandle.GetDriverName()
			impl.Init()
		}
	}

	return impl, nil
}

func (f *LinuxDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	var results []common.DriveInfo

	var s unix.Stat_t

	dir, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, common.NewNestedError("open /sys/block failed", err)
	}

	for _, ent := range dir {
		name := ent.Name()
		devPath := "/dev/"
		devPath += name
		if (!strings.Contains(name, "loop")) && (unix.Stat(devPath, &s) == nil) {
			if (s.Mode & unix.S_IFMT) == unix.S_IFBLK {
				// As CD-ROM is not supported, exclude cd-rom from probing
				if strings.Contains(name, "sr") {
					continue
				}

				driveHandle, err := f.OpenByPath(devPath)
				if err != nil {
					log.Println(err)
					continue
				}

				defer driveHandle.Close()
				results = append(results, *driveHandle.GetDriveInfo())
			}
		}
	}
	return results, nil
}

func (f *LinuxDriveFactory) EnumVolumes() (common.EnumVolumeContext, error) {
	return EnumVolumes(f)
}

func getIdInfo(path string) (string, string, string, string, error) {
	// Get model, serial from /dev/disk/by-id, has dependency to udev..?
	idPath := "/dev/disk/by-id"
	var model, serial, vendor, rev string

	dir, err := os.ReadDir(idPath)
	if err != nil {
		return "", "", "", "", common.NewNestedError("readdir /dev/disk/by-id failed", err)
	}

	devMap := make(map[string]string)

	for _, ent := range dir {
		name := ent.Name()
		if skipPat.MatchString(name) {
			continue
		}

		actualPath, err := filepath.EvalSymlinks(idPath + "/" + name)
		if err != nil {
			return "", "", "", "", common.NewNestedError("EvalSymlinks "+idPath+"/"+name+" failed", err)
		}
		devMap[name] = actualPath
	}

	for id, devPath := range devMap {
		if path == devPath {
			var temp string
			delimit := strings.LastIndex(id, "_")
			if delimit == -1 {
				continue
			}
			temp, serial = id[:delimit], id[delimit+1:]
			_, model, _ = strings.Cut(temp, "-")
			break
		}
	}

	// Get vendor name and rev version from /sys/block/{device name}/device?
	soleDev := path[strings.LastIndex(path, "/")+1:]
	b, err := os.ReadFile("/sys/block/" + soleDev + "/device/vendor")
	if err == nil {
		s := string(b)
		vendor = strings.TrimSpace(s)
	}

	b, err = os.ReadFile("/sys/block/" + soleDev + "/device/rev")
	if err == nil {
		s := string(b)
		rev = strings.TrimSpace(s)
	}

	return model, serial, vendor, rev, nil
}
