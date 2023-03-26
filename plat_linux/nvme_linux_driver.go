//go:build linux
// +build linux

package plat_linux

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/common"
)

const (
	IOCTL_SCSI_MINIPORT    = 0x4d008
	IOCTL_SCSI_GET_ADDRESS = 0x41018
)

type NvmeLinuxDriver struct {
	LinuxDriver
}

type NvmeLinuxDriverHandle struct {
	common.NvmeDriverHandle
	handle     int
	scsiHandle int
	identity   []byte
}

func NewNvmeLinuxDriver() *NvmeLinuxDriver {
	return &NvmeLinuxDriver{}
}

func (d *NvmeLinuxDriver) OpenByHandle(handle int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *NvmeLinuxDriver) GetScsiPath(handle int) (string, error) {
	sadr := SCSI_ADDRESS{}

	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(handle),
		uintptr(IOCTL_SCSI_MINIPORT),
		uintptr(unsafe.Pointer(&sadr)),
	)
	if err != 0 {
		return "", err
	}

	return fmt.Sprintf("\\\\.\\SCSI%d:", sadr.PortNumber), nil
}

func (d *NvmeLinuxDriver) QueryNvmeIdentity(handle int) ([]byte, error) {
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

	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(handle),
		uintptr(IOCTL_SCSI_MINIPORT),
		uintptr(unsafe.Pointer(&nptwb)),
	)
	if err != 0 {
		return nil, err
	}

	return nptwb.DataBuffer[:], nil
}

func (d *NvmeLinuxDriver) openImpl(handle int) (*NvmeLinuxDriverHandle, error) {
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
		_ = unix.Close(scsiHandle)
		return nil, err
	}

	return &NvmeLinuxDriverHandle{
		handle:     handle,
		scsiHandle: scsiHandle,
		identity:   identity,
	}, nil
}

func (s *NvmeLinuxDriverHandle) GetDriverName() string {
	return "NvmeLinuxDriver"
}

func (s *NvmeLinuxDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingNvme
}

func (s *NvmeLinuxDriverHandle) ReopenWritable() error {
	return nil
}

func (s *NvmeLinuxDriverHandle) Close() {
	_ = unix.Close(s.handle)
}

func (s *NvmeLinuxDriverHandle) GetNvmeIdentity() []byte {
	return s.identity
}

func (s *NvmeLinuxDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return nil
}
