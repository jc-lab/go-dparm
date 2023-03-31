//go:build linux
// +build linux

package plat_linux

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	_ "unsafe"

	"github.com/jc-lab/go-dparm/common"
)

type VolumeInfoImpl struct {
	Path string
	Filesystem string
	MountPoints []string
	DiskExtents []DISK_EXTENT
}

type EnumVolumeContextImpl struct {
	factory common.DriveFactory
	volumes []*VolumeInfoImpl
}

func (item *VolumeInfoImpl) ToVolumeInfo() common.VolumeInfo {
	return common.VolumeInfo{
		Path: item.Path,
		Filesystem: item.Filesystem,
		MountPoints: item.MountPoints,
	}
}

func (ctx *EnumVolumeContextImpl) GetList() []common.VolumeInfo {
	results := []common.VolumeInfo{}
	for _, item := range ctx.volumes {
		results = append(results, item.ToVolumeInfo())
	}
	return results
}

func (ctx *EnumVolumeContextImpl) FindVolumesByDrive(driveInfo *common.DriveInfo) []common.VolumeInfo {
	results := []common.VolumeInfo{}
	for _, volume := range ctx.volumes {
		if len(volume.DiskExtents) <= 0 {
			continue
		}
		diskExtent := volume.DiskExtents[0]
		if int(diskExtent.DiskNumber) == driveInfo.LinuxDevNum {
			results = append(results, volume.ToVolumeInfo())
		}
	}
	return results
}

func (ctx *EnumVolumeContextImpl) OpenDriveByVolumePath(volumePath string) (common.DriveHandle, error) {
	volumePath = strings.TrimSuffix(volumePath, "/")
	for _, volume := range ctx.volumes {
		if strings.TrimSuffix(volume.Path, "/") == volumePath {
			if len(volume.DiskExtents) > 0 {
				return ctx.factory.OpenByPath(fmt.Sprintf("/dev/"))
			}
			return nil, nil
		}
	}
	return nil, nil
}

func EnumVolumes(factory common.DriveFactory) (*EnumVolumeContextImpl, error) {
	impl := &EnumVolumeContextImpl{
		factory: factory,
	}
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := &VolumeInfoImpl{}

		tokenCtx := scanner.Text()
		item.Path, tokenCtx, _ = strings.Cut(tokenCtx, " ")
		mountPath, tokenCtx, _ := strings.Cut(tokenCtx, " ")
		item.Filesystem, tokenCtx, _ = strings.Cut(tokenCtx, " ")
		mountOptions, tokenCtx, _ := strings.Cut(tokenCtx, " ")

		_ = mountOptions //not used?

		item.MountPoints = append(item.MountPoints, mountPath)

		impl.volumes = append(impl.volumes, item)

		if scanner.Err() != nil {
			return nil, err
		}
	}

	return impl, nil
}