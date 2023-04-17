//go:build linux
// +build linux

package go_dparm

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_linux"

	"golang.org/x/sys/unix"
)

const (
	DIR_MAX_NUM = (1 << 32) - 1 // the max number of entry which directory can hold
)

type LinuxDriveFactory struct {
	drivers []plat_linux.LinuxDriver
}

func NewLinuxDriveFactory() *LinuxDriveFactory {
	factory := &LinuxDriveFactory{}
	factory.drivers = []plat_linux.LinuxDriver{
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

	basicInfo := plat_linux.ReadBasicInfo(fd, path)

	impl.Info.PartitionStyle = basicInfo.PartitionStyle
	impl.Info.GptDiskId = basicInfo.GptDiskId
	impl.Info.MbrDiskSignature = basicInfo.MbrSignature

	impl.Info.Model, impl.Info.Serial, impl.Info.VendorId, impl.Info.FirmwareRevision = getIdInfo(path)

	return impl, nil
}

func (f *LinuxDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	var results []common.DriveInfo

	var names []string
	var s unix.Stat_t

	dfd, err := unix.Open("/sys/block", unix.O_RDONLY | unix.O_DIRECTORY, 0o666)
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
				// As CD-ROM is not supported, exclude cd-rom from probing
				if strings.Contains(name, "sr") || strings.Contains(name, "nvme") {
					continue
				}

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

func getIdInfo(path string) (string, string, string, string) {
	// Get model, serial from /dev/disk/by-id, has dependency to udev..? 
	idPath := "/dev/disk/by-id"
	var model, serial, vendor, rev string
	_, _, _, _ = model, serial, vendor, rev // error handling for when value not found

	fd, err := unix.Open(idPath, unix.O_RDONLY | unix.O_DIRECTORY, 0o666)
	if err != nil {
		log.Fatalln(err)
	}
	defer unix.Close(fd)

	devBuf := make([]byte, 65536)
	_, err = unix.ReadDirent(fd, devBuf)
	if err != nil {
		log.Fatalln(err)
	}

	entNames := make([]string, 0)
	_, _, entNames = unix.ParseDirent(devBuf, DIR_MAX_NUM, entNames)

	devMap := make(map[string]string)

	for _, name := range entNames {
		if strings.Contains(name, "-part") || strings.Contains(name, "wwn") { // wwn not used..?
			continue
		}

		actualPath, err := filepath.EvalSymlinks(idPath + "/" + name)
		if err != nil {
			log.Fatalln(err)
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
	soleDev := path[strings.LastIndex(path, "/") + 1:]
	b, err := os.ReadFile("/sys/block/" + soleDev + "/device/vendor")
	if err != nil {
		log.Fatalln(err)
	}
	vendor = string(b)

	b, err = os.ReadFile("/sys/block/" + soleDev + "/device/rev")
	if err != nil {
		log.Fatalln(err)
	}
	rev = string(b)

	return model, serial, vendor, rev
}
