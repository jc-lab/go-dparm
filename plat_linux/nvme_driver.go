//go:build linux
// +build linux

package plat_linux

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/common"
)

type LinuxNvmeDriver struct {
	LinuxDriver
}

type LinuxNvmeDriverHandle struct {
	common.NvmeDriverHandle
	fd int
	scsiFd int
	identity [4096]byte
}

func NewLinuxNvmeDriver() *LinuxNvmeDriver {
	return &LinuxNvmeDriver{}
}

func (d *LinuxNvmeDriver) OpenByHandle(fd int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(fd)
	return driverHandle, err
}

func (d *LinuxNvmeDriver) QueryNvmeIdentity(fd int) ([]byte, error) {
	nptwb := StorageQueryWithBuffer{}

	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(SG_IO),
		uintptr(unsafe.Pointer(&nptwb)),
	)
	if err != 0 {
		return nil, err
	}

	return nptwb.Buffer[:], nil
}

func (d *LinuxNvmeDriver) openImpl(fd int) (*LinuxNvmeDriverHandle, error) {
	identity, err := d.QueryNvmeIdentity(fd)
	if err != nil {
		return nil, err
	}
	if len(identity) != 4096 {
		return nil, fmt.Errorf("invalid identity size: %d", len(identity))
	}

	driverHandle := &LinuxNvmeDriverHandle{
		fd: fd,
	}
	copy(driverHandle.identity[:], identity)

	return driverHandle, nil
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
	_ = unix.Close(s.fd)
}

func (s *LinuxNvmeDriverHandle) GetNvmeIdentity() []byte {
	return s.identity[:]
}

func (s *LinuxNvmeDriverHandle) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {

	_, _, err := unix.Syscall(
		unix.SYS_FCNTL,
		uintptr(s.fd),
		uintptr(SG_IO),
		uintptr(0),
	)
	if err != unix.Errno(0) {
		return nil, err
	}

	return nil, nil
}

func (S *LinuxNvmeDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return nil
}
