//go:build linux
// +build linux

package plat_linux

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/nvme"
)

type LinuxNvmeDriver struct {
	LinuxDriver
}

type LinuxNvmeDriverHandle struct {
	common.NvmeDriverHandle
	fd       int
	ns_id    int
	identity [4096]byte
}

func NewLinuxNvmeDriver() *LinuxNvmeDriver {
	return &LinuxNvmeDriver{}
}

func (d *LinuxNvmeDriver) OpenByFd(fd int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(fd)
	return driverHandle, err
}

func (d *LinuxNvmeDriver) openImpl(fd int) (*LinuxNvmeDriverHandle, error) {
	driverHandle := &LinuxNvmeDriverHandle{
		fd: fd,
	}
	identity, err := driverHandle.ReadIdentify(fd)
	if err != nil {
		return nil, err
	}
	if len(identity) != 4096 {
		return nil, fmt.Errorf("invalid identity size: %d", len(identity))
	}

	copy(driverHandle.identity[:], identity)

	return driverHandle, nil
}

func (s *LinuxNvmeDriverHandle) ReadIdentify(fd int) ([]byte, error) {
	// Set fd if not set
	if s.fd == 0 {
		s.fd = fd
	}

	identifyBuf := make([]byte, 4096)
	identifyCmd := nvme.NvmeAdminCmd{}
	identifyCmd.Opcode = uint8(nvme.NVME_ADMIN_OP_IDENTIFY)
	identifyCmd.Nsid = 0
	identifyCmd.Addr = *(*uint64)(unsafe.Pointer(&identifyBuf))
	identifyCmd.DataLen = 4096
	identifyCmd.Cdw10 = 1
	identifyCmd.Cdw11 = 0

	result := s.DoNvmeAdminPassthru(&identifyCmd)
	if result != nil {
		return nil, result
	}

	return identifyBuf, nil
}

func (s *LinuxNvmeDriverHandle) DoNvmeAdminPassthru(cmd *nvme.NvmeAdminCmd) error {
	data := nvme.NvmeAdminCmd{}
	data.Opcode = cmd.Opcode
	data.Flags = cmd.Flags
	data.Rsvd1 = cmd.Rsvd1
	data.Nsid = cmd.Nsid
	data.Cdw2 = cmd.Cdw2
	data.Cdw3 = cmd.Cdw3
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
		uintptr(s.fd),
		NVME_IOCTL_ADMIN_CMD,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}
	return nil
}

func (s *LinuxNvmeDriverHandle) DoNvmeIoPassthru(cmd *nvme.PassthruCmd) error {
	data := nvme.PassthruCmd{}
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
		uintptr(s.fd),
		NVME_IOCTL_IO_CMD,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}

	return nil
}

func (s *LinuxNvmeDriverHandle) DoNvmeIo(io *nvme.UserIo) error {
	data := nvme.UserIo{}
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
		uintptr(s.fd),
		NVME_IOCTL_SUBMIT_IO,
		uintptr(unsafe.Pointer(&data)),
	)
	if err != unix.Errno(0) {
		return err
	}

	return nil
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

func (s *LinuxNvmeDriverHandle) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {
	var rootError error

	offset, xferLen := 0, dataSize
	lsp, lpo, lsi := nvme.NVME_NO_LOG_LSP, offset, 0

	dataBuffer := make([]byte, dataSize)

	for {
		if offset >= dataSize {
			return dataBuffer, rootError
		}

		xferLen = dataSize - offset
		if xferLen > 4096 {
			xferLen = 4096
		}

		numd := uint32((dataSize >> 2) - 1)
		numdh := uint32((numd >> 16) & 0xffff)
		numdl := uint32(numd & 0xffff)
		cdw10 := logId | (numdl << 16) | uint32(internal.Ternary(rae, (1<<15), 0)) | (uint32(lsp) << 8)

		cmd := &nvme.NvmeAdminCmd{}
		cmd.Opcode = uint8(nvme.NVME_ADMIN_OP_GET_LOG_PAGE)
		cmd.Nsid = nsid
		cmd.Addr = *(*uint64)(unsafe.Pointer(&dataBuffer))
		cmd.DataLen = uint32(dataSize)
		cmd.Cdw10 = cdw10
		cmd.Cdw11 = numdh | uint32(lsi<<16)
		cmd.Cdw12 = uint32(lpo)
		cmd.Cdw13 = uint32(lpo >> 32)
		cmd.Cdw14 = 0

		rootError = s.DoNvmeAdminPassthru(cmd)
		if rootError == nil {
			return dataBuffer, nil
		}

		offset += xferLen
	}
}

func (s *LinuxNvmeDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return scsiSecurityCommand(s.fd, rw, dma, protocol, comId, buffer, timeoutSecs)
}
