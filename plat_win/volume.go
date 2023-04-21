//go:build windows
// +build windows

package plat_win

import (
	"fmt"
	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
	"strings"
	"unsafe"
)

type VolumeInfoImpl struct {
	Path        string
	Filesystem  string
	MountPoints []string
	DiskExtents []DISK_EXTENT
}

type EnumVolumeContextImpl struct {
	factory common.DriveFactory
	volumes []*VolumeInfoImpl
}

func (item *VolumeInfoImpl) ToVolumeInfo() common.VolumeInfo {
	return common.VolumeInfo{
		Path:        item.Path,
		Filesystem:  item.Filesystem,
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
		if int(diskExtent.DiskNumber) == driveInfo.WindowsDevNum {
			results = append(results, volume.ToVolumeInfo())
		}
	}
	return results
}

func (ctx *EnumVolumeContextImpl) OpenDriveByVolumePath(volumePath string) (common.DriveHandle, error) {
	volumePath = strings.TrimSuffix(volumePath, "\\")
	for _, volume := range ctx.volumes {
		if strings.TrimSuffix(volume.Path, "\\") == volumePath {
			if len(volume.DiskExtents) > 0 {
				return ctx.factory.OpenByPath(fmt.Sprintf("\\\\.\\PhysicalDrive%d", volume.DiskExtents[0].DiskNumber))
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

	volumeNameBuf := [320]uint16{}
	volumePathBuf := [4096]uint16{}
	fsNameBuf := [128]uint16{}
	dataBuffer := [4096]byte{}

	fvHandle, err := windows.FindFirstVolume(&volumeNameBuf[0], uint32(len(volumeNameBuf)))
	if err != nil {
		return nil, err
	}

	defer windows.FindVolumeClose(fvHandle)

	for {
		var volumePathLen uint32
		zerofill(volumePathBuf[:])

		item := &VolumeInfoImpl{}

		err = windows.GetVolumePathNamesForVolumeName(&volumeNameBuf[0], &volumePathBuf[0], uint32(len(volumePathBuf)), &volumePathLen)
		if err == nil && volumePathLen > 0 {
			for pos := 0; volumePathBuf[pos] != 0; {
				textLen := wcslen(volumePathBuf[pos:])
				mountPoint := windows.UTF16ToString(volumePathBuf[pos : pos+textLen])
				item.MountPoints = append(item.MountPoints, mountPoint)
				pos += textLen
			}
		}

		err = windows.GetVolumeInformation(&volumeNameBuf[0], nil, 0, nil, nil, nil, &fsNameBuf[0], uint32(len(fsNameBuf)))
		if err == nil {
			textLen := wcslen(fsNameBuf[:])
			item.Filesystem = windows.UTF16ToString(fsNameBuf[:textLen])
		}

		textLen := wcslen(volumeNameBuf[:])
		for ; (textLen > 0) && (volumeNameBuf[textLen-1] == '\\'); textLen-- {
			volumeNameBuf[textLen-1] = 0
		}

		item.Path = windows.UTF16ToString(volumeNameBuf[:textLen])
		handle, err := windows.CreateFile(&volumeNameBuf[0], 0, 0, nil, windows.OPEN_EXISTING, 0, 0)
		if err == nil {
			defer windows.CloseHandle(handle)

			zerofill(fsNameBuf[:])
			var bytesReturned uint32
			err = windows.DeviceIoControl(
				handle,
				IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS,
				nil,
				0,
				&dataBuffer[0],
				uint32(len(dataBuffer)),
				&bytesReturned,
				nil,
			)
			if err == nil {
				volumeDiskExtentHeader := (*VOLUME_DISK_EXTENTS)(unsafe.Pointer(&dataBuffer[0]))
				offset := unsafe.Offsetof(volumeDiskExtentHeader.Extents)
				for i := 0; i < int(volumeDiskExtentHeader.NumberOfDiskExtents); i++ {
					diskExtent := (*DISK_EXTENT)(unsafe.Pointer(&dataBuffer[offset]))
					item.DiskExtents = append(item.DiskExtents, *diskExtent)
					offset += unsafe.Sizeof(*diskExtent)
				}
			}
		}

		impl.volumes = append(impl.volumes, item)

		zerofill(volumeNameBuf[:])
		err = windows.FindNextVolume(fvHandle, &volumeNameBuf[0], uint32(len(volumeNameBuf)))
		if err != nil {
			break
		}
	}

	return impl, nil
}

