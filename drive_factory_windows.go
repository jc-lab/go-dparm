//go:build windows
// +build windows

package go_dparm

import (
	"errors"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/plat_win"
	"golang.org/x/sys/windows"
	"log"
	"syscall"
	"unsafe"
)

var GUID_DEVINTERFACE_DISK = windows.GUID{
	0x53f56307,
	0xb6bf,
	0x11d0,
	[8]byte{0x94, 0xf2, 0x00, 0xa0, 0xc9, 0x1e, 0xfb, 0x8b},
}

var (
	modsetupapi                          = windows.NewLazySystemDLL("setupapi.dll")
	procSetupDiEnumDeviceInterfaces      = modsetupapi.NewProc("SetupDiEnumDeviceInterfaces")
	procSetupDiGetDeviceInterfaceDetailW = modsetupapi.NewProc("SetupDiGetDeviceInterfaceDetailW")
)

const (
	IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS = 0x560000
)

type SP_DEVICE_INTERFACE_DATA struct {
	CbSize             uint32
	InterfaceClassGuid windows.GUID
	Flags              uint32
	Reserved           uintptr
}

type SP_DEVICE_INTERFACE_DETAIL_DATA_W struct {
	Size            uint32
	DevicePathFirst uint16
}

type DISK_EXTENT struct {
	DiskNumber     uint32
	StartingOffset uint64
	ExtentLength   uint64
}

type VOLUME_DISK_EXTENTS struct {
	NumberOfDiskExtents uint32
	Extents             [1]DISK_EXTENT
}

type WindowsDriveFactory struct {
	drivers []plat_win.WinDriver
}

func NewWindowsDriveFactory() *WindowsDriveFactory {
	factory := &WindowsDriveFactory{}
	factory.drivers = []plat_win.WinDriver{
		plat_win.NewNvmeWinDriver(),
		//windows.NewSamsungNvmeDriver(),
		plat_win.NewWindowsNvmeDriver(),
		plat_win.NewScsiDriver(),
		plat_win.NewAtaDriver(),
	}
	return factory
}

func NewSystemDriveFactory() DriveFactory {
	return NewWindowsDriveFactory()
}

func (f *WindowsDriveFactory) OpenByPath(path string) (DriveHandle, error) {
	handle, err := plat_win.OpenDevice(path)
	if err != nil {
		return nil, err
	}

	driveHandle, err := f.OpenByHandle(handle, path)
	if err == nil {
		return driveHandle, nil
	}

	_ = windows.CloseHandle(handle)

	return nil, err
}

func (f *WindowsDriveFactory) OpenByHandle(handle windows.Handle, path string) (DriveHandle, error) {
	impl := &DriveHandleImpl{}
	impl.Info.DevicePath = path
	basicInfo := plat_win.ReadBasicInfo(handle)
	if basicInfo.StorageDeviceNumber != nil {
		impl.Info.WindowsDevNum = int(basicInfo.StorageDeviceNumber.DeviceNumber)
	}
	if basicInfo.DiskGeometryEx != nil {
		impl.Info.TotalCapacity = int64(basicInfo.DiskGeometryEx.DiskSize)
	}

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByHandle(handle)
		if err == nil {
			impl.dh = driverHandle
			impl.init()
			return impl, nil
		}
	}

	return nil, errors.New("Not supported device")
}

func (f *WindowsDriveFactory) EnumDrives() ([]common.DriveInfo, error) {
	devInfoHandle, err := windows.SetupDiGetClassDevsEx(
		&GUID_DEVINTERFACE_DISK,
		"",
		0,
		windows.DIGCF_PRESENT|windows.DIGCF_DEVICEINTERFACE,
		0,
		"",
	)
	if err != nil {
		return nil, err
	}

	defer windows.SetupDiDestroyDeviceInfoList(devInfoHandle)

	var results []common.DriveInfo

	detailBuffer := make([]byte, 1024)
	for devIndex := 0; ; devIndex++ {
		devInfoData, err := windows.SetupDiEnumDeviceInfo(devInfoHandle, devIndex)
		if err != nil {
			break
		}
		_ = devInfoData

		devInterfaceData, errno := setupDiEnumInterfaceDevice(devInfoHandle, nil, &GUID_DEVINTERFACE_DISK, uint32(devIndex))
		if errno == windows.ERROR_NO_MORE_ITEMS {
			continue
		} else if errno != 0 {
			log.Println(errno)
			break
		}

		var detailSize uint32
		var dataLength uint32
		errno = setupDiGetDeviceInterfaceDetailW(devInfoHandle, devInterfaceData, &detailBuffer[0], uint32(len(detailBuffer)), &detailSize, nil)
		if errno == windows.ERROR_INSUFFICIENT_BUFFER {
			temp := detailSize % 1024
			if temp > 0 {
				detailSize += 1024 - temp
			}
			detailBuffer = make([]byte, detailSize)
			errno = setupDiGetDeviceInterfaceDetailW(devInfoHandle, devInterfaceData, &detailBuffer[0], uint32(len(detailBuffer)), &detailSize, nil)
		}
		dataLength = detailSize - 4
		if errno != 0 {
			continue
		}

		detailData := unsafe.Slice((*uint16)(unsafe.Pointer(&detailBuffer[4])), dataLength)
		devicePath := windows.UTF16ToString(detailData)
		driveHandle, err := f.OpenByPath(devicePath)
		if err != nil {
			log.Println(err)
			continue
		}

		defer driveHandle.Close()

		results = append(results, *driveHandle.GetDriveInfo())
	}

	return results, nil
}

type VolumeInfoImpl struct {
	Path        string
	Filesystem  string
	MountPoints []string
	DiskExtents []DISK_EXTENT
}

type EnumVolumeContextImpl struct {
	EnumVolumeContext
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

func (f *WindowsDriveFactory) EnumVolumes() (EnumVolumeContext, error) {
	impl := &EnumVolumeContextImpl{}

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

func zerofill(buf []uint16) {
	for i := range buf {
		buf[i] = 0
	}
}

func wcslen(buf []uint16) int {
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return i
		}
	}
	return len(buf)
}

func setupDiEnumInterfaceDevice(deviceInfoSet windows.DevInfo, deviceInfoData *windows.DevInfoData, interfaceClassGuid *windows.GUID, memberIndex uint32) (*SP_DEVICE_INTERFACE_DATA, windows.Errno) {
	var devInterfaceData SP_DEVICE_INTERFACE_DATA
	devInterfaceData.CbSize = uint32(unsafe.Sizeof(devInterfaceData))
	r1, _, e1 := syscall.SyscallN(
		procSetupDiEnumDeviceInterfaces.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
		uintptr(unsafe.Pointer(interfaceClassGuid)),
		uintptr(memberIndex),
		uintptr(unsafe.Pointer(&devInterfaceData)),
	)
	if r1 == 0 {
		return nil, e1
	}
	return &devInterfaceData, 0
}

func setupDiGetDeviceInterfaceDetailW(deviceInfoSet windows.DevInfo, devInterfaceData *SP_DEVICE_INTERFACE_DATA, deviceInterfaceDetailData *byte, deviceInterfaceDetailDataSize uint32, requiredSize *uint32, deviceInfoData *windows.DevInfoData) windows.Errno {
	if deviceInterfaceDetailData != nil {
		detailHeader := (*SP_DEVICE_INTERFACE_DETAIL_DATA_W)(unsafe.Pointer(deviceInterfaceDetailData))
		detailHeader.Size = 8
	}
	r1, _, e1 := syscall.SyscallN(
		procSetupDiGetDeviceInterfaceDetailW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(devInterfaceData)),
		uintptr(unsafe.Pointer(deviceInterfaceDetailData)),
		uintptr(deviceInterfaceDetailDataSize),
		uintptr(unsafe.Pointer(requiredSize)),
		uintptr(unsafe.Pointer(deviceInfoData)),
	)
	if r1 == 0 {
		return e1
	}
	return 0
}
