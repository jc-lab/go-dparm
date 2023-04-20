//go:build linux
// +build linux

package plat_linux

import (
	"bufio"
	"os"
	"strings"
	_ "unsafe"

	"github.com/jc-lab/go-dparm/common"
)

type VolumeInfoImpl struct {
	Path string
	Filesystem string
	MountPoints []string
}

type EnumVolumeContextImpl struct {
	factory common.DriveFactory
	volumes map[string]*VolumeInfoImpl
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
	for _, volume := range ctx.volumes {
		results = append(results, volume.ToVolumeInfo())
	}
	return results
}

func (ctx *EnumVolumeContextImpl) FindVolumesByDrive(driveInfo *common.DriveInfo) []common.VolumeInfo {
	results := []common.VolumeInfo{}
	for _, volume := range ctx.volumes {
		if volume.Path == driveInfo.DevicePath {
			results = append(results, volume.ToVolumeInfo())
		}
	}
	return results
}

func (ctx *EnumVolumeContextImpl) OpenDriveByVolumePath(volumePath string) (common.DriveHandle, error) {
	volumePath = strings.TrimSuffix(volumePath, "/")
	for _, volume := range ctx.volumes {
		if strings.TrimSuffix(volume.Path, "/") == volumePath {
			return ctx.factory.OpenByPath(volumePath)
		}
		return nil, nil
	}
	return nil, nil
}

func EnumVolumes(factory common.DriveFactory) (*EnumVolumeContextImpl, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	impl := &EnumVolumeContextImpl{
		factory: factory,
		volumes: make(map[string]*VolumeInfoImpl),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokenCtx := scanner.Text()
		devicePath, tokenCtx, _ := strings.Cut(tokenCtx, " ")
		mountPath, tokenCtx, _ := strings.Cut(tokenCtx, " ")
		filesystem, tokenCtx, _ := strings.Cut(tokenCtx, " ")
		mountOptions, tokenCtx, _ := strings.Cut(tokenCtx, " ")

		_ = mountOptions //not used?

		volumeInfo, exist := impl.volumes[devicePath]
		if !exist {
			impl.volumes[devicePath] = &VolumeInfoImpl{}
			volumeInfo = impl.volumes[devicePath]
		}

		if volumeInfo.Path == "" {
			volumeInfo.Path = devicePath
			volumeInfo.Filesystem = filesystem
		}
		volumeInfo.MountPoints = append(volumeInfo.MountPoints, mountPath)

		if err = scanner.Err(); err != nil {
			return nil, err
		}
	}

	return impl, nil
}