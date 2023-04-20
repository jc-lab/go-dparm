//go:build windows
// +build windows

package plat_win

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
)

const (
	IOCTL_SCSI_MINIPORT    = 0x4d008
	IOCTL_SCSI_GET_ADDRESS = 0x41018
)

type NvmeWinDriver struct {
	WinDriver
}

type NvmeWinDriverHandle struct {
	common.NvmeDriverHandle
	handle     windows.Handle
	scsiHandle windows.Handle
	identity   []byte
}

func NewNvmeWinDriver() *NvmeWinDriver {
	return &NvmeWinDriver{}
}

func (d *NvmeWinDriver) OpenByHandle(handle windows.Handle) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *NvmeWinDriver) GetScsiPath(handle windows.Handle) (string, error) {
	sadr := SCSI_ADDRESS{}

	var bytesReturned uint32
	err := windows.DeviceIoControl(
		handle,
		IOCTL_SCSI_GET_ADDRESS,
		(*byte)(unsafe.Pointer(&sadr)),
		uint32(unsafe.Sizeof(sadr)),
		(*byte)(unsafe.Pointer(&sadr)),
		uint32(unsafe.Sizeof(sadr)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\\\\.\\SCSI%d:", sadr.PortNumber), nil
}

func (d *NvmeWinDriver) QueryNvmeIdentity(handle windows.Handle) ([]byte, error) {
	nptwb := NVME_PASS_THROUGH_IOCTL{}

	nptwb.SrbIoCtrl.ControlCode = NVME_PASS_THROUGH_SRB_IO_CODE
	nptwb.SrbIoCtrl.HeaderLength = uint32(unsafe.Sizeof(nptwb.SrbIoCtrl))
	copyFromAsciiToBuffer(nptwb.SrbIoCtrl.Signature[:], NVME_SIG_STR)
	nptwb.SrbIoCtrl.Timeout = NVME_PT_TIMEOUT
	nptwb.SrbIoCtrl.Length = uint32(unsafe.Sizeof(nptwb) - unsafe.Sizeof(nptwb.SrbIoCtrl))
	nptwb.DataBufferLen = uint32(unsafe.Sizeof(nptwb.DataBuffer))
	nptwb.ReturnBufferLen = uint32(unsafe.Sizeof(nptwb))
	nptwb.Direction = NVME_FROM_DEV_TO_HOST

	pcommand := (*NVMe_COMMAND)(unsafe.Pointer(&nptwb.NVMeCmd))
	pcommand.CDW0.OPC = NVME_ADMIN_OP_IDENTIFY
	// https://nvmexpress.org/wp-content/uploads/NVM_Express_Revision_1.3.pdf
	// Page 112
	// The Identify Controller data structure is returned to the host for this controller.
	pcommand.CDW10_OR_NDP = 1

	var bytesReturned uint32
	err := windows.DeviceIoControl(
		handle,
		IOCTL_SCSI_MINIPORT,
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

	return nptwb.DataBuffer[:], nil
}

func (d *NvmeWinDriver) openImpl(handle windows.Handle) (*NvmeWinDriverHandle, error) {
	scsiPath, err := d.GetScsiPath(handle)
	if err != nil {
		return nil, err
	}

	scsiHandle, err := OpenDevice(scsiPath)
	if err != nil {
		return nil, err
	}

	identity, err := d.QueryNvmeIdentity(scsiHandle)
	if err != nil {
		_ = windows.CloseHandle(scsiHandle)
		return nil, err
	}
	if len(identity) != 4096 {
		_ = windows.CloseHandle(scsiHandle)
		return nil, errors.New(fmt.Sprintf("invalid identity size: %d", len(identity)))
	}

	driverHandle := &NvmeWinDriverHandle{
		handle:     handle,
		scsiHandle: scsiHandle,
	}
	copy(driverHandle.identity[:], identity)

	return driverHandle, nil
}

func (s *NvmeWinDriverHandle) GetDriverName() string {
	return "NvmeWinDriver"
}

func (s *NvmeWinDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingNvme
}

func (s *NvmeWinDriverHandle) ReopenWritable() error {
	return nil
}

func (s *NvmeWinDriverHandle) Close() {
	_ = windows.CloseHandle(s.handle)
}

func (s *NvmeWinDriverHandle) GetIdentity() []byte {
	return s.identity
}

func (s *NvmeWinDriverHandle) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {
	return nil, errors.New("Not supported")
}

func (s *NvmeWinDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return errors.New("Not supported")
}
