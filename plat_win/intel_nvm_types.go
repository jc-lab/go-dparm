package plat_win

const (
	// INTEL_NVME_SIG_STR is the SRB_IO_CONTROL signature used by the
	// Intel RST (Rapid Storage Technology) NVMe miniport.
	INTEL_NVME_SIG_STR = "IntelNvm"

	// IOCTL_INTEL_NVME_PASS_THROUGH == CTL_CODE(0xF000, 0xA02, METHOD_BUFFERED, FILE_ANY_ACCESS)
	IOCTL_INTEL_NVME_PASS_THROUGH = 0xF0002808

	INTEL_NVME_PASS_THROUGH_VERSION = 1
	INTEL_NVME_PT_TIMEOUT           = 10

	INTEL_NVME_CMD_DW_SIZE      = 16
	INTEL_NVME_COMPLETE_DW_SIZE = 4
	INTEL_NVME_DATA_BUFFER_LEN  = 0x1000
)

// SRB status values (srb.h) reported back through SRB_IO_CONTROL.ReturnCode,
// which the WDK documents as the SrbStatus of the miniport request.
// The upper two bits are flags, so the status has to be masked out.
const (
	SRB_STATUS_MASK    = 0x3F
	SRB_STATUS_PENDING = 0x00
	SRB_STATUS_SUCCESS = 0x01
)

// Indices into INTEL_NVME_PASS_THROUGH.Cmd, which holds the 64 byte NVMe
// command as a DWORD array. It is kept as an array (instead of a struct with
// 64bit MPTR/PRP fields) so that the whole request structure stays 4 byte
// aligned, exactly matching the packed layout the miniport expects on both
// 32bit and 64bit Windows.
const (
	INTEL_NVME_CMD_DW_CDW0    = 0
	INTEL_NVME_CMD_DW_NSID    = 1
	INTEL_NVME_CMD_DW_MPTR_LO = 4
	INTEL_NVME_CMD_DW_MPTR_HI = 5
	INTEL_NVME_CMD_DW_PRP1_LO = 6
	INTEL_NVME_CMD_DW_PRP1_HI = 7
	INTEL_NVME_CMD_DW_PRP2_LO = 8
	INTEL_NVME_CMD_DW_PRP2_HI = 9
	INTEL_NVME_CMD_DW_CDW10   = 10
)

// INTEL_NVME_PASS_THROUGH is the payload of an IOCTL_SCSI_MINIPORT request
// against \\.\Scsi<Port>: handled by the Intel RST miniport.
//
//	0x0000 ~ 0x001B  SRB_IO_CONTROL
//	0x001C ~ 0x001F  Version / PathId / TargetId / Lun
//	0x0020 ~ 0x00A3  NVME_PASS_THROUGH_PARAMETERS
//	0x00A4 ~ 0x10A3  DataBuffer
type INTEL_NVME_PASS_THROUGH struct {
	SrbIoCtrl SRB_IO_CONTROL

	Version  byte /* 0x001C */
	PathId   byte /* 0x001D */ // port number from Windows (non-RAID), or the CSMI port number
	TargetId byte /* 0x001E */ // reserved
	Lun      byte /* 0x001F */ // reserved

	/* NVME_PASS_THROUGH_PARAMETERS */
	Cmd [INTEL_NVME_CMD_DW_SIZE]uint32 /* 0x0020 ~ 0x005F */

	// IsIOCommandSet is 0 for the admin queue, 1 for the I/O command set.
	// It is followed by padding, because the field must be DWORD aligned.
	IsIOCommandSet        byte    /* 0x0060 */
	IsIOCommandSetPadding [3]byte /* 0x0061 ~ 0x0063 */

	// CplEntry is the COMPLETION_QUEUE_ENTRY (completion DW0 ~ DW3).
	// Note that CrystalDiskInfo names this field CplEntry[4] at 0x0060 and
	// QueueId at 0x0070, which is off by one DWORD; the values it writes are
	// the same, so only the completion read back is affected.
	CplEntry [INTEL_NVME_COMPLETE_DW_SIZE]uint32 /* 0x0064 ~ 0x0073 */

	// DataBufferOffset is a DWORD aligned offset from the beginning of the
	// SRB_IO_CONTROL, i.e. offsetof(DataBuffer).
	DataBufferOffset uint32   /* 0x0074 */
	DataBufferLength uint32   /* 0x0078 */
	Rsvd2            [40]byte /* 0x007C ~ 0x00A3 */

	DataBuffer [INTEL_NVME_DATA_BUFFER_LEN]byte /* 0x00A4 ~ */
}

// GetCompletionStatus returns the Status Field of the completion queue entry
// (DW3 bit 17~31: SC, SCT, CRD, M, DNR). Zero means the command succeeded.
func (p *INTEL_NVME_PASS_THROUGH) GetCompletionStatus() uint32 {
	return p.CplEntry[3] >> 17
}

func (p *INTEL_NVME_PASS_THROUGH) SetOpcode(opcode byte) {
	p.Cmd[INTEL_NVME_CMD_DW_CDW0] = (p.Cmd[INTEL_NVME_CMD_DW_CDW0] & 0xFFFFFF00) | uint32(opcode)
}

func (p *INTEL_NVME_PASS_THROUGH) SetFuse(fuse byte) {
	p.Cmd[INTEL_NVME_CMD_DW_CDW0] = (p.Cmd[INTEL_NVME_CMD_DW_CDW0] & ^uint32(0x03<<8)) | (uint32(fuse&0x03) << 8)
}

func (p *INTEL_NVME_PASS_THROUGH) SetNsid(nsid uint32) {
	p.Cmd[INTEL_NVME_CMD_DW_NSID] = nsid
}

func (p *INTEL_NVME_PASS_THROUGH) SetMptr(value uint64) {
	p.Cmd[INTEL_NVME_CMD_DW_MPTR_LO] = uint32(value)
	p.Cmd[INTEL_NVME_CMD_DW_MPTR_HI] = uint32(value >> 32)
}

func (p *INTEL_NVME_PASS_THROUGH) SetPrp1(value uint64) {
	p.Cmd[INTEL_NVME_CMD_DW_PRP1_LO] = uint32(value)
	p.Cmd[INTEL_NVME_CMD_DW_PRP1_HI] = uint32(value >> 32)
}

func (p *INTEL_NVME_PASS_THROUGH) SetPrp2(value uint64) {
	p.Cmd[INTEL_NVME_CMD_DW_PRP2_LO] = uint32(value)
	p.Cmd[INTEL_NVME_CMD_DW_PRP2_HI] = uint32(value >> 32)
}

// SetCdw sets a command specific DWORD, where n is in the range 10 ~ 15.
func (p *INTEL_NVME_PASS_THROUGH) SetCdw(n int, value uint32) {
	if n < 10 || n > 15 {
		panic("invalid cdw index")
	}
	p.Cmd[INTEL_NVME_CMD_DW_CDW10+(n-10)] = value
}

func (p *INTEL_NVME_PASS_THROUGH) GetCdw(n int) uint32 {
	if n < 10 || n > 15 {
		panic("invalid cdw index")
	}
	return p.Cmd[INTEL_NVME_CMD_DW_CDW10+(n-10)]
}
