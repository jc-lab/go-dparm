//go:build linux
// +build linux

package plat_linux

import (
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/common"
)

const (
	IOCTL_STORAGE_QUERY_PROPERTY = 0x2d1400
)

type LinuxNvmeDriver struct {
	LinuxDriver
}

type LinuxNvmeDriverHandle struct {
	common.NvmeDriverHandle
	handle   int
	identity []byte
}

func NewLinuxNvmeDriver() *LinuxNvmeDriver {
	return &LinuxNvmeDriver{}
}

func (d *LinuxNvmeDriver) OpenByHandle(handle int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *LinuxNvmeDriver) QueryNvmeIdentity(handle int) ([]byte, error) {
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
	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(handle),
		uintptr(unsafe.Pointer(&nptwb)),
		uintptr(unsafe.Pointer(&bytesReturned)),
	)
	if err != 0 {
		return nil, err
	}

	return nptwb.Buffer[:], nil
}

func (d *LinuxNvmeDriver) openImpl(handle int) (*LinuxNvmeDriverHandle, error) {
	identity, err := d.QueryNvmeIdentity(handle)
	if err != nil {
		return nil, err
	}

	return &LinuxNvmeDriverHandle{
		handle:   handle,
		identity: identity,
	}, nil
}

func (s *LinuxNvmeDriverHandle) GetDriverName() string {
	return "LinuxNvmeDriver"
}

func (s *LinuxNvmeDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingNvme
}

func (s *LinuxNvmeDriverHandle) ReopenWritable() error {
	return nil
}

func (s *LinuxNvmeDriverHandle) Close() {
	_ = unix.Close(s.handle)
}

func (s *LinuxNvmeDriverHandle) GetNvmeIdentity() []byte {
	return s.identity
}

func (S *LinuxNvmeDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return nil
}
