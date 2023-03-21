//go:build windows
// +build windows

package plat_win

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	IOCTL_STORAGE_GET_DEVICE_NUMBER  = 0x2d1080
	IOCTL_DISK_GET_DRIVE_GEOMETRY_EX = 0x700a0
)

type WinBasicInfo struct {
	StorageDeviceNumber *STORAGE_DEVICE_NUMBER
	DiskGeometryEx      *DISK_GEOMETRY_EX
}

func OpenDevice(path string) (windows.Handle, error) {
	sPath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	return windows.CreateFile(
		sPath,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
}

func ReadBasicInfo(handle windows.Handle) *WinBasicInfo {
	result := &WinBasicInfo{}
	storageDeviceNumber := STORAGE_DEVICE_NUMBER{}
	diskGeometryEx := DISK_GEOMETRY_EX{}

	var bytesReturned uint32

	err := windows.DeviceIoControl(
		handle,
		IOCTL_STORAGE_GET_DEVICE_NUMBER,
		nil,
		0,
		(*byte)(unsafe.Pointer(&storageDeviceNumber)),
		uint32(unsafe.Sizeof(storageDeviceNumber)),
		&bytesReturned,
		nil,
	)
	if err == nil {
		result.StorageDeviceNumber = &storageDeviceNumber
	}

	err = windows.DeviceIoControl(
		handle,
		IOCTL_DISK_GET_DRIVE_GEOMETRY_EX,
		nil,
		0,
		(*byte)(unsafe.Pointer(&diskGeometryEx)),
		uint32(unsafe.Sizeof(diskGeometryEx)),
		&bytesReturned,
		nil,
	)
	if err == nil {
		result.DiskGeometryEx = &diskGeometryEx
	}

	return result
}

func copyToPointer(dest unsafe.Pointer, src []byte, len int) {
	destRef := unsafe.Slice((*byte)(dest), len)
	copy(destRef, src[:len])
}

func copyFromPointer(dest []byte, src unsafe.Pointer, len int) {
	srcRef := unsafe.Slice((*byte)(src), len)
	copy(dest, srcRef[:len])
}

func copyFromAsciiToBuffer(dest []byte, text string) {
	c := len(text)
	for i := 0; i < c; i++ {
		dest[i] = text[i]
	}
}
