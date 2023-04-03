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
	fd int
}

type LinuxNvmeDriverHandle struct {
	common.NvmeDriverHandle
	fd int
	ns_id int
	identity [4096]byte
}

func NewLinuxNvmeDriver() *LinuxNvmeDriver {
	return &LinuxNvmeDriver{}
}

func (d *LinuxNvmeDriver) OpenByFd(fd int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(fd)
	return driverHandle, err
}

func (d *LinuxNvmeDriver) ReadIdentify(fd int) ([]byte, error) {
	identifyBuf := make([]byte, 4096)
	identifyCmd := NvmeAdminCmd{}
	identifyCmd.Opcode = NVME_ADMIN_OP_IDENTIFY
	identifyCmd.Nsid = 0
	identifyCmd.Addr = uintptr(unsafe.Pointer(&identifyBuf))
	identifyCmd.DataLen = 4096
	identifyCmd.Cdw10 = 1
	identifyCmd.Cdw11 = 0
	result := d.doNvmeAdminPassthru(&identifyCmd)
	if result != nil {
		return nil, result
	}

	return identifyBuf[:], nil
}

func (d *LinuxNvmeDriver) doNvmeAdminPassthru(cmd *NvmeAdminCmd) error {
	data := NvmeIoctlAdminCmd{}
	data.Opcode = cmd.Opcode
	data.Flags = cmd.Flags
	data.Rsvd1 = cmd.Rsvd1
	data.Nsid = cmd.Nsid
	data.Cdw2 = cmd.Cdw2
	data.Metadata = cmd.Metadata
	data.Addr = uint64(cmd.Addr)
	data.MetadataLen = cmd.MetadataLen
	data.DataLen = cmd.DataLen
	data.Cdw10 = cmd.Cdw10
	data.Cdw11 = cmd.Cdw11
	data.Cdw12 = cmd.Cdw12
	data.Cdw13 = cmd.Cdw13
	data.Cdw14 = cmd.Cdw14
	data.Cdw15 = cmd.Cdw15
	data.TimeoutMs = cmd.TimeoutMs
	data.Result = cmd.Result
	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(d.fd),
		NVME_IOCTL_ADMIN_CMD,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}
	return nil
}

func (d *LinuxNvmeDriver) doNvmeIoPassthru(cmd *NvmePassthruCmd) error {
	data := NvmeIoctlPassthruCmd{}
	data.Opcode = cmd.Opcode
	data.Flags = cmd.Flags
	data.Rsvd1 = cmd.Rsvd1
	data.Nsid = cmd.Nsid
	data.Cdw2 = cmd.Cdw2
	data.Cdw3 = cmd.Cdw3
	data.Metadata = cmd.Metadata
	data.Addr = uint64(cmd.Addr)
	data.MetadataLen = cmd.MetadataLen
	data.Cdw10 = cmd.Cdw10
	data.Cdw11 = cmd.Cdw11
	data.Cdw12 = cmd.Cdw12
	data.Cdw13 = cmd.Cdw13
	data.Cdw14 = cmd.Cdw14
	data.Cdw15 = cmd.Cdw15
	data.TimeoutMs = cmd.TimeoutMs
	data.Result = cmd.Result
	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(d.fd),
		NVME_IOCTL_IO_CMD,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}

	return nil
}

func (d *LinuxNvmeDriver) DoNvmeIo(io *NvmeUserIo) error {
	data := NvmeIoctlUserIo{}
	data.Opcode = io.Opcode
	data.Flags = io.Flags
	data.Control = io.Control
	data.Nblocks = io.Nblocks
	data.Rsvd = io.Rsvd
	data.Metadata = io.Metadata
	data.Addr = io.Addr
	data.Slba = io.Slba
	data.Dsmgmt = io.Dsmgmt
	data.Reftag = io.Reftag
	data.Apptag = io.Apptag
	data.Appmask = io.Appmask
	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(d.fd),
		NVME_IOCTL_SUBMIT_IO,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}
	
	return nil
}

func (d *LinuxNvmeDriver) openImpl(fd int) (*LinuxNvmeDriverHandle, error) {
	d.fd = fd
	identity, err := d.ReadIdentify(fd)
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

func (s *LinuxNvmeDriverHandle) GetIdentity() []byte {
	return s.identity[:]
}

func (s *LinuxNvmeDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return scsiSecurityCommand(s.fd, rw, dma, protocol, comId, buffer, timeoutSecs)
}