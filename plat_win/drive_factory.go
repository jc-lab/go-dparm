//go:build windows
// +build windows

package plat_win

import (
	"log"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
)

var GUID_DEVINTERFACE_DISK = windows.GUID{
	0x53f56307,
	0xb6bf,
	0x11d0,
	[8]byte{0x94, 0xf2, 0x00, 0xa0, 0xc9, 0x1e, 0xfb, 0x8b},
}

type WindowsDriveFactory struct {
	drivers []WinDriver
}

func NewWindowsDriveFactory() *WindowsDriveFactory {
	factory := &WindowsDriveFactory{}
	factory.drivers = []WinDriver{
		NewNvmeWinDriver(),
		//windows.NewSamsungNvmeDriver(),
		NewWindowsNvmeDriver(),
		NewScsiDriver(),
		NewAtaDriver(),
	}
	return factory
}

func (f *WindowsDriveFactory) OpenByPath(path string) (common.DriveHandle, error) {
	handle, err := OpenDevice(path)
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

func (f *WindowsDriveFactory) OpenByHandle(handle windows.Handle, path string) (common.DriveHandle, error) {
	impl := &common.DriveHandleImpl{}
	impl.Info.DrivingType = common.DrivingUnknown
	impl.Info.DevicePath = path
	impl.Info.Removable = -1

	basicInfo := ReadBasicInfo(handle)

	impl.Info.PartitionStyle = basicInfo.PartitionStyle
	impl.Info.GptDiskId = basicInfo.GptDiskId
	impl.Info.MbrDiskSignature = basicInfo.MbrSignature
	impl.Info.Partitions = basicInfo.Partitions

	if basicInfo.StorageDeviceNumber != nil {
		impl.Info.WindowsDevNum = int(basicInfo.StorageDeviceNumber.DeviceNumber)
	}

	if basicInfo.DiskGeometryEx != nil {
		impl.Info.TotalCapacity = int64(basicInfo.DiskGeometryEx.DiskSize)
	}

	for _, driver := range f.drivers {
		driverHandle, err := driver.OpenByHandle(handle)
		if err == nil {
			impl.Dh = driverHandle
			impl.Info.DrivingType = driverHandle.GetDrivingType()
			impl.Info.DriverName = driverHandle.GetDriverName()
			impl.Init()
			return impl, nil
		}
	}

	storageQueryResp, err := ReadStorageQuery(handle)
	if err == nil && storageQueryResp != nil {
		if impl.Info.Model == "" {
			impl.Info.Model = storageQueryResp.ProductId
		}
		if impl.Info.Serial == "" {
			impl.Info.Serial = storageQueryResp.SerialNumber
		}

		impl.Info.VendorId = storageQueryResp.VendorId
		impl.Info.ProductRevision = storageQueryResp.ProductRevision
		if storageQueryResp.RemovableMedia {
			impl.Info.Removable = 1
		} else {
			impl.Info.Removable = 0
		}
		switch storageQueryResp.DeviceType {
		case 0x00:
			impl.Info.DriveType = common.DriveTypeFixed
		case 0x01:
			impl.Info.DriveType = common.DriveTypeTape
		case 0x05:
			impl.Info.DriveType = common.DriveTypeCdrom
		case 0x09:
			impl.Info.DriveType = common.DriveTypeRemote
		}
	}

	return impl, nil
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

		devInterfaceData, errno := SetupDiEnumInterfaceDevice(devInfoHandle, nil, &GUID_DEVINTERFACE_DISK, uint32(devIndex))
		if errno == windows.ERROR_NO_MORE_ITEMS {
			continue
		} else if errno != 0 {
			log.Println(errno)
			break
		}

		var detailSize uint32
		var dataLength uint32
		errno = SetupDiGetDeviceInterfaceDetailW(devInfoHandle, devInterfaceData, &detailBuffer[0], uint32(len(detailBuffer)), &detailSize, nil)
		if errno == windows.ERROR_INSUFFICIENT_BUFFER {
			temp := detailSize % 1024
			if temp > 0 {
				detailSize += 1024 - temp
			}
			detailBuffer = make([]byte, detailSize)
			errno = SetupDiGetDeviceInterfaceDetailW(devInfoHandle, devInterfaceData, &detailBuffer[0], uint32(len(detailBuffer)), &detailSize, nil)
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

func (f *WindowsDriveFactory) EnumVolumes() (common.EnumVolumeContext, error) {
	return EnumVolumes(f)
}
