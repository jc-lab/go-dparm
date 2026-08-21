//go:build windows
// +build windows

package plat_win

import (
	"errors"
	"fmt"
	"log"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/nvme"
	"golang.org/x/sys/windows"
)

// IntelNvmDriver talks to NVMe devices that are hidden behind the Intel RST
// (Rapid Storage Technology) storage miniport, where the standard NVMe
// pass-through interfaces are not available:
//
//	\\.\PhysicalDriveN
//	  -> IOCTL_SCSI_GET_ADDRESS  => Port / Path / Target / Lun
//	  -> \\.\Scsi<Port>:
//	     -> IOCTL_SCSI_MINIPORT ("IntelNvm", 0xF0002808)
type IntelNvmDriver struct {
	WinDriver
}

type IntelNvmDriverHandle struct {
	common.NvmeDriverHandle
	handle      windows.Handle
	scsiHandle  windows.Handle
	scsiAddress SCSI_ADDRESS
	identity    []byte
}

func NewIntelNvmDriver() *IntelNvmDriver {
	return &IntelNvmDriver{}
}

func (d *IntelNvmDriver) OpenByHandle(handle windows.Handle) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *IntelNvmDriver) GetScsiAddress(handle windows.Handle) (*SCSI_ADDRESS, error) {
	sadr := SCSI_ADDRESS{}
	sadr.Length = uint32(unsafe.Sizeof(sadr))

	var bytesReturned uint32
	err := windows.DeviceIoControl(
		handle,
		IOCTL_SCSI_GET_ADDRESS,
		nil,
		0,
		(*byte)(unsafe.Pointer(&sadr)),
		uint32(unsafe.Sizeof(sadr)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &sadr, nil
}

// QueryNvmeIdentity reads the 4096 byte Identify Controller data structure
// (Opcode 0x06, CNS 0x01) through the Intel RST miniport.
func (d *IntelNvmDriver) QueryNvmeIdentity(handle windows.Handle, scsiAddress *SCSI_ADDRESS) ([]byte, error) {
	nptwb := newIntelNvmPassThrough(scsiAddress)

	nptwb.SetOpcode(NVME_ADMIN_OP_IDENTIFY)
	nptwb.SetNsid(0)
	// https://nvmexpress.org/wp-content/uploads/NVM_Express_Revision_1.3.pdf
	// Page 112
	// CDW10.CNS = 1: The Identify Controller data structure is returned to
	// the host for this controller.
	nptwb.SetCdw(10, NVME_IDENTIFY_CNS_CONTROLLER)

	if err := intelNvmIoctl(handle, nptwb); err != nil {
		return nil, err
	}

	identity := make([]byte, len(nptwb.DataBuffer))
	copy(identity, nptwb.DataBuffer[:])

	return identity, nil
}

func (d *IntelNvmDriver) openImpl(handle windows.Handle) (*IntelNvmDriverHandle, error) {
	scsiAddress, err := d.GetScsiAddress(handle)
	if err != nil {
		log.Printf("\t[IntelNvm] GetScsiAddress failed: %+v", err)
		return nil, err
	}

	scsiPath := fmt.Sprintf("\\\\.\\SCSI%d:", scsiAddress.PortNumber)
	log.Printf("\t[IntelNvm] scsi address: port=%d, path=%d, target=%d, lun=%d (%s)",
		scsiAddress.PortNumber, scsiAddress.PathId, scsiAddress.TargetId, scsiAddress.Lun, scsiPath)

	scsiHandle, err := OpenDevice(scsiPath)
	if err != nil {
		log.Printf("\t[IntelNvm] OpenDevice failed: %+v", err)
		return nil, err
	}

	identity, err := d.QueryNvmeIdentity(scsiHandle, scsiAddress)
	if err != nil {
		_ = windows.CloseHandle(scsiHandle)
		return nil, err
	}
	if len(identity) != 4096 {
		_ = windows.CloseHandle(scsiHandle)
		return nil, fmt.Errorf("invalid identity size: %d", len(identity))
	}
	if isAllZero(identity) {
		// The miniport accepted the request but did not return anything;
		// this is not an Intel RST managed NVMe device.
		_ = windows.CloseHandle(scsiHandle)
		return nil, errors.New("empty identify controller data")
	}

	driverHandle := &IntelNvmDriverHandle{
		handle:      handle,
		scsiHandle:  scsiHandle,
		scsiAddress: *scsiAddress,
		identity:    identity,
	}

	return driverHandle, nil
}

func (s *IntelNvmDriverHandle) GetDriverName() string {
	return "IntelNvmDriver"
}

func (s *IntelNvmDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingNvme
}

func (s *IntelNvmDriverHandle) ReopenWritable() error {
	return nil
}

func (s *IntelNvmDriverHandle) Close() {
	if s.scsiHandle != windows.InvalidHandle {
		_ = windows.CloseHandle(s.scsiHandle)
		s.scsiHandle = windows.InvalidHandle
	}
	_ = windows.CloseHandle(s.handle)
}

func (s *IntelNvmDriverHandle) GetIdentity() []byte {
	return s.identity
}

func (s *IntelNvmDriverHandle) DoNvmeAdminPassthru(cmd *nvme.NvmeAdminCmd) error {
	nptwb := newIntelNvmPassThrough(&s.scsiAddress)

	if cmd.DataLen > uint32(len(nptwb.DataBuffer)) {
		return errors.New("too long data")
	}
	dataRef := cmd.DataBuffer
	if cmd.DataAddr != 0 {
		dataRef = unsafe.Slice((*byte)(unsafe.Pointer(cmd.DataAddr)), cmd.DataLen)
	}
	copy(nptwb.DataBuffer[:], dataRef)

	nptwb.SetOpcode(cmd.Opcode)
	nptwb.SetNsid(cmd.Nsid)
	nptwb.SetCdw(10, cmd.Cdw10)
	nptwb.SetCdw(11, cmd.Cdw11)
	nptwb.SetCdw(12, cmd.Cdw12)
	nptwb.SetCdw(13, cmd.Cdw13)
	nptwb.SetCdw(14, cmd.Cdw14)
	nptwb.SetCdw(15, cmd.Cdw15)

	if err := intelNvmIoctl(s.scsiHandle, nptwb); err != nil {
		return err
	}

	copy(dataRef, nptwb.DataBuffer[:cmd.DataLen])
	cmd.Result = nptwb.CplEntry[0]

	return nil
}

func (s *IntelNvmDriverHandle) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {
	if dataSize <= 0 || dataSize > INTEL_NVME_DATA_BUFFER_LEN {
		return nil, fmt.Errorf("invalid data size: %d", dataSize)
	}
	if (dataSize % 4) != 0 {
		return nil, fmt.Errorf("data size must be dword aligned: %d", dataSize)
	}

	nptwb := newIntelNvmPassThrough(&s.scsiAddress)

	// CDW10: LID (bit 0~7), RAE (bit 15), NUMDL (bit 16~31).
	// NUMDL is the number of dwords to transfer, minus one.
	cdw10 := (logId & 0xFF) | (uint32(dataSize/4-1) << 16)
	if rae {
		cdw10 |= 1 << 15
	}

	nptwb.SetOpcode(NVME_ADMIN_OP_GET_LOG_PAGE)
	nptwb.SetNsid(nsid)
	nptwb.SetCdw(10, cdw10)

	if err := intelNvmIoctl(s.scsiHandle, nptwb); err != nil {
		return nil, err
	}

	result := make([]byte, dataSize)
	copy(result, nptwb.DataBuffer[:dataSize])

	return result, nil
}

func (s *IntelNvmDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return errors.New("not supported")
}

func newIntelNvmPassThrough(scsiAddress *SCSI_ADDRESS) *INTEL_NVME_PASS_THROUGH {
	nptwb := &INTEL_NVME_PASS_THROUGH{}

	nptwb.SrbIoCtrl.HeaderLength = uint32(unsafe.Sizeof(nptwb.SrbIoCtrl))
	copyFromAsciiToBuffer(nptwb.SrbIoCtrl.Signature[:], INTEL_NVME_SIG_STR)
	nptwb.SrbIoCtrl.Timeout = INTEL_NVME_PT_TIMEOUT
	nptwb.SrbIoCtrl.ControlCode = IOCTL_INTEL_NVME_PASS_THROUGH
	nptwb.SrbIoCtrl.Length = uint32(unsafe.Sizeof(*nptwb) - unsafe.Sizeof(nptwb.SrbIoCtrl))

	nptwb.Version = INTEL_NVME_PASS_THROUGH_VERSION
	if scsiAddress != nil {
		// PathId alone selects the device. TargetId and Lun are reserved in
		// this interface, so they are left as zero; Seagate's opensea-transport
		// does the same on the non-CSMI path, and CrystalDiskInfo enumerates
		// devices by PathId only.
		nptwb.PathId = scsiAddress.PathId
	}
	nptwb.IsIOCommandSet = 0 // admin queue
	nptwb.DataBufferOffset = uint32(unsafe.Offsetof(nptwb.DataBuffer))
	nptwb.DataBufferLength = uint32(unsafe.Sizeof(nptwb.DataBuffer))

	return nptwb
}

func intelNvmIoctl(handle windows.Handle, nptwb *INTEL_NVME_PASS_THROUGH) error {
	var bytesReturned uint32
	err := windows.DeviceIoControl(
		handle,
		IOCTL_SCSI_MINIPORT,
		(*byte)(unsafe.Pointer(nptwb)),
		uint32(unsafe.Sizeof(*nptwb)),
		(*byte)(unsafe.Pointer(nptwb)),
		uint32(unsafe.Sizeof(*nptwb)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return err
	}

	// The miniport reports the result of the request in ReturnCode, which the
	// WDK documents as the SrbStatus member. SRB_STATUS_SUCCESS (0x01) is what
	// the Intel RST miniport returns on success; PENDING (0x00) is the value we
	// sent in, meaning the miniport left it untouched.
	srbStatus := nptwb.SrbIoCtrl.ReturnCode & SRB_STATUS_MASK
	if srbStatus != SRB_STATUS_SUCCESS && srbStatus != SRB_STATUS_PENDING {
		return fmt.Errorf("intel nvme passthrough failed: srb status=0x%x", nptwb.SrbIoCtrl.ReturnCode)
	}

	// The NVMe level result lives in the completion queue entry, not in the
	// SRB status.
	if status := nptwb.GetCompletionStatus(); status != 0 {
		return fmt.Errorf("nvme command failed: status=0x%x", status)
	}

	return nil
}

func isAllZero(buf []byte) bool {
	for _, b := range buf {
		if b != 0 {
			return false
		}
	}
	return true
}
