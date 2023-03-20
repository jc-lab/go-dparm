package plat_win

const (
	NVME_STORPORT_DRIVER               = 0xE000
	NVME_PASS_THROUGH_SRB_IO_CODE      = 0xe0002000
	NVME_SIG_STR                       = "NvmeMini"
	NVME_SIG_STR_LEN                   = 8
	NVME_FROM_DEV_TO_HOST              = 2
	NVME_IOCTL_VENDOR_SPECIFIC_DW_SIZE = 6
	NVME_IOCTL_CMD_DW_SIZE             = 16
	NVME_IOCTL_COMPLETE_DW_SIZE        = 4
	NVME_PT_TIMEOUT                    = 40
)

type NvmeAdminOpCode = byte

const (
	NVME_ADMIN_OP_DELETE_SQ      NvmeAdminOpCode = 0x00
	NVME_ADMIN_OP_CREATE_SQ      NvmeAdminOpCode = 0x01
	NVME_ADMIN_OP_GET_LOG_PAGE   NvmeAdminOpCode = 0x02
	NVME_ADMIN_OP_DELETE_CQ      NvmeAdminOpCode = 0x04
	NVME_ADMIN_OP_CREATE_CQ      NvmeAdminOpCode = 0x05
	NVME_ADMIN_OP_IDENTIFY       NvmeAdminOpCode = 0x06
	NVME_ADMIN_OP_ABORT_CMD      NvmeAdminOpCode = 0x08
	NVME_ADMIN_OP_SET_FEATURES   NvmeAdminOpCode = 0x09
	NVME_ADMIN_OP_GET_FEATURES   NvmeAdminOpCode = 0x0A
	NVME_ADMIN_OP_ASYNC_EVENT    NvmeAdminOpCode = 0x0C
	NVME_ADMIN_OP_NS_MGMT        NvmeAdminOpCode = 0x0D
	NVME_ADMIN_OP_ACTIVATE_FW    NvmeAdminOpCode = 0x10
	NVME_ADMIN_OP_DOWNLOAD_FW    NvmeAdminOpCode = 0x11
	NVME_ADMIN_OP_DEV_SELF_TEST  NvmeAdminOpCode = 0x14
	NVME_ADMIN_OP_NS_ATTACH      NvmeAdminOpCode = 0x15
	NVME_ADMIN_OP_KEEP_ALIVE     NvmeAdminOpCode = 0x18
	NVME_ADMIN_OP_DIRECTIVE_SEND NvmeAdminOpCode = 0x19
	NVME_ADMIN_OP_DIRECTIVE_RECV NvmeAdminOpCode = 0x1A
	NVME_ADMIN_OP_VIRTUAL_MGMT   NvmeAdminOpCode = 0x1C
	NVME_ADMIN_OP_NVME_MI_SEND   NvmeAdminOpCode = 0x1D
	NVME_ADMIN_OP_NVME_MI_RECV   NvmeAdminOpCode = 0x1E
	NVME_ADMIN_OP_DBBUF          NvmeAdminOpCode = 0x7C
	NVME_ADMIN_OP_FORMAT_NVM     NvmeAdminOpCode = 0x80
	NVME_ADMIN_OP_SECURITY_SEND  NvmeAdminOpCode = 0x81
	NVME_ADMIN_OP_SECURITY_RECV  NvmeAdminOpCode = 0x82
	NVME_ADMIN_OP_SANITIZE_NVM   NvmeAdminOpCode = 0x84
	NVME_ADMIN_OP_GET_LBA_STATUS NvmeAdminOpCode = 0x86
)

type SRB_IO_CONTROL struct {
	HeaderLength uint32
	Signature    [8]byte
	Timeout      uint32
	ControlCode  uint32
	ReturnCode   uint32
	Length       uint32
}

type NVME_PASS_THROUGH_IOCTL struct {
	SrbIoCtrl       SRB_IO_CONTROL
	VendorSpecific  [NVME_IOCTL_VENDOR_SPECIFIC_DW_SIZE]uint32
	NVMeCmd         [NVME_IOCTL_CMD_DW_SIZE]uint32
	CplEntry        [NVME_IOCTL_COMPLETE_DW_SIZE]uint32
	Direction       uint32
	QueueId         uint32
	DataBufferLen   uint32
	MetaDataLen     uint32
	ReturnBufferLen uint32
	DataBuffer      [4096]byte
}

type NVMe_COMMAND_DWORD_0 struct {
	OPC byte
	B01 byte
	CID uint16
}

func (p *NVMe_COMMAND_DWORD_0) SetFuse(value byte) {
	p.B01 = (p.B01 & 0xFC) | (value & 0x03)
}

func (p *NVMe_COMMAND_DWORD_0) GetFuse() byte {
	return p.B01 & 0x03
}

type NVMe_COMMAND struct {
	/*
	 * [Command Dword 0] This field is common to all commands and is defined
	 * in Figure 6.
	 */
	CDW0 NVMe_COMMAND_DWORD_0

	/*
	 * [Namespace Identifier] This field indicates the namespace that this
	 * command applies to. If the namespace is not used for the command then
	 * this field shall be cleared to 0h. If a command shall be applied to all
	 * namespaces on the device then this value shall be set to FFFFFFFFh.
	 */
	NSID uint32

	/* DWORD 2 3 */
	DW02 uint32
	DW03 uint32

	/*
	 * [Metadata Pointer] This field contains the address of a contiguous
	 * physical buffer of metadata. This field is only used if metadata is not
	 * interleaved with the LBA data as specified in the Format NVM command.
	 * This field shall be Dword aligned.
	 */
	MPTR uint64

	/* [PRP Entry 1] This field contains the first PRP entry for the command. */
	PRP1 uint64

	/*
	 * [PRP Entry 2] This field contains the second PRP entry for the command.
	 * If the data transfer spans more than two memory pages then this field is
	 * a PRP List pointer.
	 */
	PRP2 uint64

	/* [Command Dword 10] This field is command specific Dword 10. */
	/*
	 * Defined in Admin and NVM Vendor Specific Command format.
	 * Number of DWORDs in PRP data transfer (in Figure 8).
	 */
	CDW10_OR_NDP uint32

	/* [Command Dword 11] This field is command specific Dword 11. */
	/*
	 * Defined in Admin and NVM Vendor Specific Command format.
	 * Number of DWORDs in MPTR Metadata transfer (in Figure 8).
	 */
	CDW11_OR_NDM uint32

	/* [Command Dword 12] This field is command specific Dword 12. */
	CDW12 uint32

	/* [Command Dword 13] This field is command specific Dword 13. */
	CDW13 uint32

	/* [Command Dword 14] This field is command specific Dword 14. */
	CDW14 uint32

	/* [Command Dword 15] This field is command specific Dword 15. */
	CDW15 uint32
}
