package plat_linux

type u8 uint8
type le16 uint16
type le32 uint32

//
// Nvme status codes
//
const (
	NVME_SC_SUCCESS = 0x0
	NVME_SC_INVALID_OPCODE = 0x1
	NVME_SC_ONVALID_FIELD = 0x2
	NVME_SC_CMDID_CONFLICT = 0x3
	NVME_SC_DARA_XFER_ERROR = 0x4
	NVME_SC_POWER_LOSS = 0x5
	NVME_SC_INTERNAL = 0x6
	NVME_SC_ABORT_REQ = 0x7
	NVME_SC_ABORT_QUEUE = 0x8
	NVME_SC_FUSED_FAIL = 0x9
	NVME_SC_FUSE_MISSING = 0xa
	NVME_SC_INVALID_NS = 0xb
	NVME_SC_CMD_SEQ_ERROR = 0xc
	NVME_SC_SGL_INVALID_LAST = 0xd
	NVME_SC_SGL_INVALID_COUNT = 0xe
	NVME_SC_SGL_INVALID_DATA = 0xf
	NVME_SC_SGL_INVALID_METADATA = 0x10
	NVME_SC_SGL_INVALID_TYPE = 0x11
	NVME_SC_CMB_INVALID_USE = 0x12
	NVME_SC_PRP_INVALID_OFFSET = 0x13
	NVME_SC_ATOMIC_WRITE_UNIT_EXCEEDED = 0x14
	NVME_SC_OPERATION_DENIED = 0x15
	NVME_SC_SGL_INVALID_OFFSET = 0x16

	NVME_SC_INCONSISTENT_HOST_ID = 0x18
	NVME_SC_KEEP_ALIVE_EXPIRED = 0x19
	NVME_SC_KEEP_ALIVE_INVALID = 0x1A
	NVME_SC_PREEMPT_ABORT = 0x1B
	NVME_SC_SANITIZE_FAILED = 0x1C
	NVME_SC_SANITIZE_IN_PROGRESS = 0x1D

	NVME_SC_NS_WRITE_PROTECTED = 0x20
	NVME_SC_CMD_INTERRUPTED = 0x21
	NVME_SC_TRANSIENT_TRANSPORT = 0x22

	NVME_SC_LBA_RANGE = 0x80
	NVME_SC_CAP_EXCEEDED = 0x81
	NVME_SC_NS_NOT_READY = 0x82
	NVME_SC_RESERVATION_CONFLICT = 0x83
	NVME_SC_FORMAT_IN_PROGRESS = 0x84

	//
	// Command Specitic Status
	//
	NVME_SC_CQ_INVALID = 0x100
	NVME_SC_QID_INVALID = 0x101
	NVME_SC_QUEUE_SIZE = 0x102
	NVME_SC_ABORT_LIMIT = 0x103
	NVME_SC_ABORT_MISSING = 0x104
	NVME_SC_ASYNC_LIMIT = 0x105
	NVME_SC_FIRMWARE_SLOT = 0x106
	NVME_SC_FIRMWARE_IMAGE = 0x107
	NVME_SC_INVALID_VECTOR = 0x108
	NVME_SC_INVALID_LOG_PAGE = 0x109
	NVME_SC_INVALID_FORMAT = 0x10a
	NVME_SC_FW_NEEDS_CONV_RESET = 0x10b
	NVME_SC_INVALID_QUEUE = 0x10c
	NVME_SC_FEATURE_NOT_SAVEABLE = 0x10d
	NVME_SC_FEATURE_NOT_CHANGEABLE = 0x10e
	NVME_SC_FEATURE_NOT_PRE_NS = 0x10f
	NVME_SC_FW_NEEDS_SUBSYS_RESET = 0x110
	NVME_SC_FW_NEEDS_RESET = 0x111
	NVME_SC_FW_NEEDS_MAX_TIME = 0x112
	NVME_SC_FW_ACTIVATE_PROHIBITED = 0x113
	NVME_SC_OVERLAPPING_RANGE = 0x114
	NVME_SC_NS_INSUFFICIENT_CAP = 0x115
	NVME_SC_NS_ID_UNAVAILABLE = 0x116
	NVME_SC_NS_ALREADY_ATTACHED	= 0x118
	NVME_SC_NS_IS_PRIVATE = 0x119
	NVME_SC_NS_NOT_ATTACHED	= 0x11a
	NVME_SC_THIN_PROV_NOT_SUPP = 0x11b
	NVME_SC_CTRL_LIST_INVALID = 0x11c
	NVME_SC_DEVICE_SELF_TEST_IN_PROGRESS= 0x11d
	NVME_SC_BP_WRITE_PROHIBITED	= 0x11e
	NVME_SC_INVALID_CTRL_ID	= 0x11f
	NVME_SC_INVALID_SECONDARY_CTRL_STATE = 0x120
	NVME_SC_INVALID_NUM_CTRL_RESOURCE = 0x121
	NVME_SC_INVALID_RESOURCE_ID	= 0x122
	NVME_SC_PMR_SAN_PROHIBITED = 0x123
	NVME_SC_ANA_INVALID_GROUP_ID= 0x124
	NVME_SC_ANA_ATTACH_FAIL = 0x125

	//
	// Command Set Specific - Namespace Types commands:
	//
	NVME_SC_IOCS_NOT_SUPPORTED = 0x129
	NVME_SC_IOCS_NOT_ENABLED = 0x12A
	NVME_SC_IOCS_COMBINATION_REJECTED = 0x12B
	NVME_SC_INVALID_IOCS = 0x12C

	//
    // I/O Command Set Specific - NVM commands:
	//
	NVME_SC_BAD_ATTRIBUTES = 0x180
	NVME_SC_INVALID_PI = 0x181
	NVME_SC_READ_ONLY = 0x182
	NVME_SC_CMD_SIZE_LIMIT_EXCEEDED = 0x183

	//
	// I/O Command Set Specific - Fabrics commands:
	//
	NVME_SC_CONNECT_FORMAT = 0x180
	NVME_SC_CONNECT_CTRL_BUSY = 0x181
	NVME_SC_CONNECT_INVALID_PARAM = 0x182
	NVME_SC_CONNECT_RESTART_DISC = 0x183
	NVME_SC_CONNECT_INVALID_HOST = 0x184

	NVME_SC_DISCOVERY_RESTART = 0x190
	NVME_SC_AUTH_REQUIRED = 0x191

	//
	// I.O Command Set Specific - Zoned Namespace commands:
	//
	NVME_SC_ZONE_BOUNDARY_ERROR	= 0x1B8
	NVME_SC_ZONE_IS_FULL = 0x1B9
	NVME_SC_ZONE_IS_READ_ONLY = 0x1BA
	NVME_SC_ZONE_IS_OFFLINE	= 0x1BB
	NVME_SC_ZONE_INVALID_WRITE = 0x1BC
	NVME_SC_TOO_MANY_ACTIVE_ZONES = 0x1BD
	NVME_SC_TOO_MANY_OPEN_ZONES	= 0x1BE
	NVME_SC_ZONE_INVALID_STATE_TRANSITION = 0x1BF

	//
	// Media and Data Integrity Errors:
	//
	NVME_SC_WRITE_FAULT	= 0x280
	NVME_SC_READ_ERROR = 0x281
	NVME_SC_GUARD_CHECK	= 0x282
	NVME_SC_APPTAG_CHECK = 0x283
	NVME_SC_REFTAG_CHECK = 0x284
	NVME_SC_COMPARE_FAILED = 0x285
	NVME_SC_ACCESS_DENIED = 0x286
	NVME_SC_UNWRITTEN_BLOCK = 0x287

	//
	// Path-related Errors:
	//
	NVME_SC_ANA_PERSISTENT_LOSS	= 0x301
	NVME_SC_ANA_INACCESSIBLE = 0x302
	NVME_SC_ANA_TRANSITION = 0x303

	NVME_SC_CRD = 0x1800
	NVME_SC_DNR = 0x4000
)

// Nvme Admin Op Code
const (
	NVME_ADMIN_OP_DELETE_SQ	= 0x00
	NVME_ADMIN_OP_CREATE_SQ	= 0x01
	NVME_ADMIN_OP_GET_LOG_PAGE = 0x02
	NVME_ADMIN_OP_DELETE_CQ = 0x04
	NVME_ADMIN_OP_CREATE_CQ	= 0x05
	NVME_ADMIN_OP_IDENTIFY = 0x06
	NVME_ADMIN_OP_ABORT_CMD	= 0x08
	NVME_ADMIN_OP_SET_FEATURES = 0x09
	NVME_ADMIN_OP_GET_FEATURES = 0x0A
	NVME_ADMIN_OP_ASYNC_EVENT = 0x0C
	NVME_ADMIN_OP_NS_MGMT = 0x0D
	NVME_ADMIN_OP_ACTIVATE_FW = 0x10
	NVME_ADMIN_OP_DOWNLOAD_FW = 0x11
	NVME_ADMIN_OP_DEV_SELF_TEST	= 0x14
	NVME_ADMIN_OP_NS_ATTACH	= 0x15
	NVME_ADMIN_OP_KEEP_ALIVE = 0x18
	NVME_ADMIN_OP_DIRECTIVE_SEND = 0x19
	NVME_ADMIN_OP_DIRECTIVE_RECV = 0x1A
	NVME_ADMIN_OP_VIRTUAL_MGMT = 0x1C
	NVME_ADMIN_OP_NVME_MI_SEND = 0x1D
	NVME_ADMIN_OP_NVME_MI_RECV = 0x1E
	NVME_ADMIN_OP_DBBUF = 0x7C
	NVME_ADMIN_OP_FORMAT_NVM = 0x80
	NVME_ADMIN_OP_SECURITY_SEND	= 0x81
	NVME_ADMIN_OP_SECURITY_RECV	= 0x82
	NVME_ADMIN_OP_SANITIZE_NVM = 0x84
	NVME_ADMIN_OP_GET_LBA_STATUS = 0x86
)

//
// NVM_Express_Revision_1.3.pdf
// Figure 86 : Get Log Page - Command Dword 10
const (
	NVME_NO_LOG_LSP = 0x0
	NVME_NO_LOG_LPO = 0x0
	NVME_LOG_ANA_LSP_RGO = 0x1
	NVME_TELEM_LSP_CREATE = 0x1
)

//
// NVM_Express_Revision_1.3.pdf
// Figure 90 : Get Log Page - Log Page Identifiers
//
const (
	NVME_GET_LOG_PAGE_RESERVED = 0x00
	NVME_GET_LOG_PAGE_ERROR_INFO = 0x01
	NVME_GET_LOG_PAGE_SMART = 0x02
	NVME_GET_LOG_PAGE_FIRMWARE_SLOT_INFO = 0x03
)

//
// NVM_Express_Revision_1.3.pdf
// Figure 178 : Sanitize - Command Dword 10
//
const (
	NVME_SANITIZE_NO_DEALLOC = 0x00000200
	NVME_SANITIZE_OIPBP	 = 0x00000100
	NVME_SANITIZE_OWPASS_SHIFT	= 0x00000004 /* 07:04 */
	NVME_SANITIZE_AUSE = 0x00000008
	NVME_SANITIZE_ACT_CRYPTO_ERASE = 0x00000004
	NVME_SANITIZE_ACT_OVERWRITE	= 0x00000003
	NVME_SANITIZE_ACT_BLOCK_ERASE = 0x00000002
	NVME_SANITIZE_ACT_EXIT = 0x00000001

	/* Sanitize Monitor/Log */
	NVME_SANITIZE_LOG_DATA_LEN = 0x0014
	NVME_SANITIZE_LOG_GLOBAL_DATA_ERASED = 0x0100
	NVME_SANITIZE_LOG_NUM_CMPLTED_PASS_MASK	= 0x00F8
	NVME_SANITIZE_LOG_STATUS_MASK = 0x0007
	NVME_SANITIZE_LOG_NEVER_SANITIZED = 0x0000
	NVME_SANITIZE_LOG_COMPLETED_SUCCESS	= 0x0001
	NVME_SANITIZE_LOG_IN_PROGESS = 0x0002
	NVME_SANITIZE_LOG_COMPLETED_FAILED = 0x0003
	NVME_SANITIZE_LOG_ND_COMPLETED_SUCCESS = 0x0004
)

type NvmeUserIo struct {
	Opcode uint8
	Flags uint8
	Control uint16
	Nblocks uint16
	Rsvd uint16
	Metadata uint64
	Addr uint64
	Slba uint64
	Dsmgmt uint32
	Reftag uint32
	Apptag uint16
	Appmask uint16
}

type NvmePassthruCmd struct {
	Opcode uint8
	Flags uint8
	Rsvd1 uint16
	Nsid uint32
	Cdw2 uint32
	Cdw3 uint32
	Metadata uint64
	Addr uintptr
	MetadataLen uint32
	DataLen uint32
	Cdw10 uint32
	Cdw11 uint32
	Cdw12 uint32
	Cdw13 uint32
	Cdw14 uint32
	Cdw15 uint32
	TimeoutMs uint32
	Result uint32
}

type NvmeAdminCmd NvmePassthruCmd

type NvmePassthruCmd64 struct {
	Opcode uint8
	Flags uint8
	Rsvd1 uint16
	Nsid uint32
	Cdw2 uint32
	Cdw3 uint32
	Metadata uint64
	Addr uint64
	MetadataLen uint32
	DataLen uint32
	Cdw10 uint32
	Cdw11 uint32
	Cdw12 uint32
	Cdw13 uint32
	Cdw14 uint32
	Cdw15 uint32
	TimeoutMs uint32
	Rsvd2 uint32
	Result uint64
}

type NvmeIdentifyPowerState struct {
	MaxPower le16
	Rsvd2 u8
	Flags u8
	EntryLat le32
	ExitLat le32
	ReadTput u8
	ReadLat u8
	WriteTput u8
	WriteLat u8
	IdlePower le16
	IdleScale u8
	Rsvd19 u8
	ActivePower le16
	ActiveWorkScale u8
	Rsvd23 [9]u8
}

type NvmeIdentifycontroller struct {
	Vid le16
	Ssvid le16
	Sn [20]uint8
	Mn [40]uint8
	Fr [8]uint8
	Rab u8
	Ieee [3]u8
	Cmic u8
	Mdts u8
	Cntlid le16
	Ver le32
	Rtd3r le32
	Rtd3e le32
	Oaes le32
	Rrls le16
	Rsvd102 [9]u8
	Cntrltype u8
	Fguid [16]uint8
	Crdt1 le16
	Crdt2 le16
	Crdt3 le16
	Rsvd134 [122]u8
	Oacs le16
	Acl u8
	Aerl u8
	Frmw u8
	Lpa u8
	Elpe u8
	Npss u8
	Avcss u8
	Apsta u8
	Wctemp le16
	Cctemp le16
	Mtfa le16
	Hmpre le32
	Hmmin le32
	Tnvmvap [16]u8
	Unvmcap [16]u8
	Rpmbs le32
	Edstt le16
	Dsto u8
	Fwug u8
	Kas le16
	Hctma le16
	Mntmt le16
	Mxtmt le16
	Sanicap le32
	Hmminds le32
	Hmmaxd le16
	Nsetidmax le16
	Endgidmax le16
	Anatt u8
	Anacap u8
	Anagrpmax le32
	Nanagrpid le32
	Pels le32
	Rsvd356 [156]u8
	Sqes u8
	Cqes u8
	Maccmd le16
	Nn le32
	Oncs le16
	Fuses le16
	Fna u8
	Vwc u8
	Awun le16
	Awupf le16
	Icsvscc u8
	Nwpc u8
	Acwu le16
	Ocfs le16
	Sgls le32
	Mnan le32
	Rsvd544 [224]u8
	Subnqn [256]uint8
	Rsvd1024 [768]u8
	Ioccsz le32
	Iorcsz le32
	Icdoff le16
	Ctrattr u8
	Msdbd u8
	Rsvd1804 [244]u8
	Psd [32]NvmeIdentifyPowerState
	Vs [1024]u8
}

//
// NVM_Express_Revision_1.3.pdf
// 5.14.1.9.2 Sanitize Status (Log Identifier 81h)
//
type NvmeSanitizeLogPage struct {
	Progress le16
	Status le16
	Cdw10Info le32
	EstOverwriteTime le32
	EstBlockEraseTime le32
	EstCryptoEraseTime le32
	// extended
	EstOverwriteTimeWithNoDeallocate le32
	EstBlockEraseTimeWithNoDeallocate le32
	EstCryptoEraseTimeWithNoDeallocate le32
}

//
// NVM_Express_Revision_1.3.pdf
// 5.14.1.2 SMART/Health Information (Log Identifier 02h)
//
type NvmeSmartLogPage struct {
	CriticalWarning uint8
	CompositeTemperature uint16
	AvailableSpare uint8
	AvailableSpareThreshold uint8
	PercentageUsed uint8
	Rev01 [26]uint8
	DataUnitsRead [16]uint8
	DataUnitsWritten [16]uint8
	ControllerBusyTime [16]uint8
	PowerCycles [16]uint8
	PowerOnHours [16]uint8
	UnsafeShutdowns [16]uint8
	MediaAndDataIntegrityErrors [16]uint8
	NumberOfErrorInformationLogEntries [16]uint8
	WarningCompositeTemperatureTime [4]uint8
	CriticalCompositeTemperatureTime [4]uint8
	RevRemaining [280]uint8
}

type NvmeAdminCmdWithBuffer struct {
	Cmd NvmeAdminCmd
	Buffer [4096]byte
}