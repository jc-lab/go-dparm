//go:build windows
// +build windows

package plat_win

import (
	"errors"
	"fmt"
	"github.com/jc-lab/go-dparm/nvme"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
)

type WindowsNvmeDriver struct {
	WinDriver
}

type WindowsNvmeDriverHandle struct {
	common.NvmeDriverHandle
	handle   windows.Handle
	identity [4096]byte
}

func NewWindowsNvmeDriver() *WindowsNvmeDriver {
	return &WindowsNvmeDriver{}
}

func (d *WindowsNvmeDriver) OpenByHandle(handle windows.Handle) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *WindowsNvmeDriver) QueryNvmeIdentity(handle windows.Handle) ([]byte, error) {
	nptwb := StorageQueryWithBuffer{}

	nptwb.Query.PropertyId = StorageAdapterProtocolSpecificProperty
	nptwb.Query.QueryType = PropertyStandardQuery
	nptwb.ProtocolSpecific.ProtocolType = ProtocolTypeNvme
	nptwb.ProtocolSpecific.DataType = NVMeDataTypeIdentify
	nptwb.ProtocolSpecific.ProtocolDataRequestValue = NVME_IDENTIFY_CNS_CONTROLLER
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue = 0
	nptwb.ProtocolSpecific.ProtocolDataOffset = uint32(unsafe.Offsetof(nptwb.Buffer) - unsafe.Sizeof(nptwb.Query))
	nptwb.ProtocolSpecific.ProtocolDataLength = uint32(unsafe.Sizeof(nptwb.Buffer))

	var bytesReturned uint32
	err := windows.DeviceIoControl(
		handle,
		IOCTL_STORAGE_QUERY_PROPERTY,
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return nptwb.Buffer[:], nil
}

func (d *WindowsNvmeDriver) openImpl(handle windows.Handle) (*WindowsNvmeDriverHandle, error) {
	identity, err := d.QueryNvmeIdentity(handle)
	if err != nil {
		return nil, err
	}
	if len(identity) != 4096 {
		return nil, errors.New(fmt.Sprintf("invalid identity size: %d", len(identity)))
	}

	driverHandle := &WindowsNvmeDriverHandle{
		handle: handle,
	}
	copy(driverHandle.identity[:], identity)

	return driverHandle, nil
}

func (s *WindowsNvmeDriverHandle) GetDriverName() string {
	return "WindowsNvmeDriver"
}

func (s *WindowsNvmeDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingNvme
}

func (s *WindowsNvmeDriverHandle) ReopenWritable() error {
	return nil
}

func (s *WindowsNvmeDriverHandle) Close() {
	_ = windows.CloseHandle(s.handle)
}

func (s *WindowsNvmeDriverHandle) GetIdentity() []byte {
	return s.identity[:]
}

func (s *WindowsNvmeDriverHandle) DoNvmeAdminPassthru(cmd *nvme.NvmeAdminCmd) error {
	return errors.New("not supported")
}

func (s *WindowsNvmeDriverHandle) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {
	headerSize := int(unsafe.Sizeof(StorageQueryWithoutBuffer{}))
	buffer := make([]byte, headerSize+dataSize)
	nptwb := (*StorageQueryWithoutBuffer)(unsafe.Pointer(&buffer[0]))

	nptwb.Query.PropertyId = StorageDeviceProtocolSpecificProperty
	nptwb.Query.QueryType = PropertyStandardQuery
	nptwb.ProtocolSpecific.ProtocolType = ProtocolTypeNvme
	nptwb.ProtocolSpecific.DataType = NVMeDataTypeLogPage
	nptwb.ProtocolSpecific.ProtocolDataRequestValue = logId

	// logId == NVME_LOG_PAGE_HEALTH_INFO
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue = 0  // This will be passed as the lower 32 bit of log page offset if controller supports extended data for the Get Log Page.
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue2 = 0 // This will be passed as the higher 32 bit of log page offset if controller supports extended data for the Get Log Page.
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue3 = 0 // This will be passed as Log Specific Identifier in CDW11.
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue4 = 0 // This will map to STORAGE_PROTOCOL_DATA_SUBVALUE_GET_LOG_PAGE definition, then user can pass Retain Asynchronous Event, Log Specific Field.

	nptwb.ProtocolSpecific.ProtocolDataOffset = uint32(uintptr(headerSize) - unsafe.Sizeof(nptwb.Query))
	nptwb.ProtocolSpecific.ProtocolDataLength = uint32(dataSize)

	var bytesReturned uint32
	err := windows.DeviceIoControl(
		s.handle,
		IOCTL_STORAGE_QUERY_PROPERTY,
		&buffer[0],
		uint32(len(buffer)),
		&buffer[0],
		uint32(len(buffer)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, err
	}

	sizeOfSTORAGE_PROTOCOL_DATA_DESCRIPTOR := uint32(unsafe.Sizeof(STORAGE_PROTOCOL_DATA_DESCRIPTOR{}))
	resultProtocolDataDesc := (*STORAGE_PROTOCOL_DATA_DESCRIPTOR)(unsafe.Pointer(&buffer[0]))

	if resultProtocolDataDesc.Version != sizeOfSTORAGE_PROTOCOL_DATA_DESCRIPTOR || resultProtocolDataDesc.Size != sizeOfSTORAGE_PROTOCOL_DATA_DESCRIPTOR {
		return nil, fmt.Errorf("invalid response header: 0x%x, 0x%x", resultProtocolDataDesc.Version, resultProtocolDataDesc.Size)
	}

	return buffer[headerSize:], nil
}

func (s *WindowsNvmeDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return scsiSecurityCommand(s.handle, rw, dma, protocol, comId, buffer, timeoutSecs)
}
