//go:build linux
// +build linux

package plat_linux

import (
	"bufio"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/diskfs/go-diskfs/partition/mbr"
	"os"
	"regexp"
	"strconv"
	"strings"
	_ "unsafe"

	"github.com/jc-lab/go-dparm/common"
)

var partitionDeviceRegex = regexp.MustCompile("^(/dev/.+\\d)p(\\d+)$")

type PartitionImpl struct {
	devicePath string
	parentPath string
	start      uint64
	size       uint64
	style      common.PartitionStyle
	mbrInfo    *common.MbrPartitionInfo
	gptInfo    *common.GptPartitionInfo
}

type VolumeInfoImpl struct {
	Path        string
	Filesystem  string
	MountPoints []string
	Partition   *PartitionImpl
}

type EnumVolumeContextImpl struct {
	factory common.DriveFactory
	volumes map[string]*VolumeInfoImpl
}

func (item *VolumeInfoImpl) ToVolumeInfo() common.VolumeInfo {
	out := common.VolumeInfo{
		Path:        item.Path,
		Filesystem:  item.Filesystem,
		MountPoints: item.MountPoints,
	}
	if item.Partition != nil {
		out.Partitions = append(out.Partitions, item.Partition)
	}
	return out
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

func (ctx *EnumVolumeContextImpl) OpenDriveByPartition(partition common.Partition) (common.DriveHandle, error) {
	partitionImpl := partition.(*PartitionImpl)
	return ctx.factory.OpenByPath(partitionImpl.parentPath)
}

func EnumVolumes(factory common.DriveFactory) (*EnumVolumeContextImpl, error) {
	cachedBasics := map[string]*LinuxBasicInfo{}

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
			volumeInfo = &VolumeInfoImpl{}
			impl.volumes[devicePath] = volumeInfo
		}

		if volumeInfo.Path == "" {
			volumeInfo.Path = devicePath
			volumeInfo.Filesystem = filesystem
		}
		volumeInfo.MountPoints = append(volumeInfo.MountPoints, mountPath)

		matches := partitionDeviceRegex.FindStringSubmatch(devicePath)
		if len(matches) == 3 {
			parentPath := matches[1]
			partNum, _ := strconv.Atoi(matches[2])
			partNum -= 1
			basicInfo, exists := cachedBasics[parentPath]
			if !exists {
				fd, err := OpenDevice(parentPath)
				if err == nil {
					basicInfo, err = ReadBasicInfo(fd, parentPath)
				}
			}
			if basicInfo != nil && basicInfo.PartitionTable != nil {
				if basicInfo.PartitionStyle == common.PartitionStyleGpt {
					table := basicInfo.PartitionTable.(*gpt.Table)
					if partNum < len(table.Partitions) {
						partition := table.Partitions[partNum]
						volumeInfo.Partition = &PartitionImpl{
							devicePath: devicePath,
							parentPath: parentPath,
							style:      basicInfo.PartitionStyle,
							start:      uint64(partition.GetStart()),
							size:       uint64(partition.GetStart()),
							gptInfo: &common.GptPartitionInfo{
								PartitionType: strings.ToUpper("{" + string(partition.Type) + "}"),
								PartitionId:   strings.ToUpper("{" + partition.GUID + "}"),
							},
						}
					}
				} else if basicInfo.PartitionStyle == common.PartitionStyleMbr {
					table := basicInfo.PartitionTable.(*mbr.Table)
					if partNum < len(table.Partitions) {
						partition := table.Partitions[partNum]
						volumeInfo.Partition = &PartitionImpl{
							devicePath: devicePath,
							parentPath: parentPath,
							style:      basicInfo.PartitionStyle,
							start:      uint64(partition.GetStart()),
							size:       uint64(partition.GetStart()),
							mbrInfo: &common.MbrPartitionInfo{
								PartitionType: byte(partition.Type),
								BootIndicator: partition.Bootable,
							},
						}
					}
				}
			}
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}
	}

	return impl, nil
}

func (p *PartitionImpl) GetStart() uint64 {
	return p.start
}

func (p *PartitionImpl) GetSize() uint64 {
	return p.size
}

func (p *PartitionImpl) GetPartitionStyle() common.PartitionStyle {
	return p.style
}

func (p *PartitionImpl) GetMbrInfo() *common.MbrPartitionInfo {
	return p.mbrInfo
}

func (p *PartitionImpl) GetGptInfo() *common.GptPartitionInfo {
	return p.gptInfo
}
