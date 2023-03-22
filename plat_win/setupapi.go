//go:build windows
// +build windows

package plat_win

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	modsetupapi                          = windows.NewLazySystemDLL("setupapi.dll")
	procSetupDiEnumDeviceInterfaces      = modsetupapi.NewProc("SetupDiEnumDeviceInterfaces")
	procSetupDiGetDeviceInterfaceDetailW = modsetupapi.NewProc("SetupDiGetDeviceInterfaceDetailW")
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

func SetupDiEnumInterfaceDevice(deviceInfoSet windows.DevInfo, deviceInfoData *windows.DevInfoData, interfaceClassGuid *windows.GUID, memberIndex uint32) (*SP_DEVICE_INTERFACE_DATA, windows.Errno) {
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

func SetupDiGetDeviceInterfaceDetailW(deviceInfoSet windows.DevInfo, devInterfaceData *SP_DEVICE_INTERFACE_DATA, deviceInterfaceDetailData *byte, deviceInterfaceDetailDataSize uint32, requiredSize *uint32, deviceInfoData *windows.DevInfoData) windows.Errno {
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
