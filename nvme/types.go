package nvme

type StatusCode uint16
type AdminOpCode uint8
type GetLogPageCommand uint8
type GetLogPageIdentifier uint8
type SanitizeCommand uint32

/**
 * NVM_Express_Revision_1.3.pdf
 * Figure 31 : Status Code - Generic Command Status values
 */
const (
	NVME_SC_SUCCESS                    = StatusCode(0x0)
	NVME_SC_INVALID_OPCODE             = StatusCode(0x1)
	NVME_SC_INVALID_FIELD              = StatusCode(0x2)
	NVME_SC_CMDID_CONFLICT             = StatusCode(0x3)
	NVME_SC_DATA_XFER_ERROR            = StatusCode(0x4)
	NVME_SC_POWER_LOSS                 = StatusCode(0x5)
	NVME_SC_INTERNAL                   = StatusCode(0x6)
	NVME_SC_ABORT_REQ                  = StatusCode(0x7)
	NVME_SC_ABORT_QUEUE                = StatusCode(0x8)
	NVME_SC_FUSED_FAIL                 = StatusCode(0x9)
	NVME_SC_FUSED_MISSING              = StatusCode(0xa)
	NVME_SC_INVALID_NS                 = StatusCode(0xb)
	NVME_SC_CMD_SEQ_ERROR              = StatusCode(0xc)
	NVME_SC_SGL_INVALID_LAST           = StatusCode(0xd)
	NVME_SC_SGL_INVALID_COUNT          = StatusCode(0xe)
	NVME_SC_SGL_INVALID_DATA           = StatusCode(0xf)
	NVME_SC_SGL_INVALID_METADATA       = StatusCode(0x10)
	NVME_SC_SGL_INVALIDYPE             = StatusCode(0x11)
	NVME_SC_CMB_INVALID_USE            = StatusCode(0x12)
	NVME_SC_PRP_INVALID_OFFSET         = StatusCode(0x13)
	NVME_SC_ATOMIC_WRITE_UNIT_EXCEEDED = StatusCode(0x14)
	NVME_SC_OPERATION_DENIED           = StatusCode(0x15)
	NVME_SC_SGL_INVALID_OFFSET         = StatusCode(0x16)

	NVME_SC_INCONSISTENT_HOST_ID = StatusCode(0x18)
	NVME_SC_KEEP_ALIVE_EXPIRED   = StatusCode(0x19)
	NVME_SC_KEEP_ALIVE_INVALID   = StatusCode(0x1A)
	NVME_SC_PREEMPT_ABORT        = StatusCode(0x1B)
	NVME_SC_SANITIZE_FAILED      = StatusCode(0x1C)
	NVME_SC_SANITIZE_IN_PROGRESS = StatusCode(0x1D)

	NVME_SC_NS_WRITE_PROTECTED = StatusCode(0x20)
	NVME_SC_CMD_INTERRUPTED    = StatusCode(0x21)
	NVME_SCRANSIENTRANSPORT    = StatusCode(0x22)

	NVME_SC_LBA_RANGE            = StatusCode(0x80)
	NVME_SC_CAP_EXCEEDED         = StatusCode(0x81)
	NVME_SC_NS_NOT_READY         = StatusCode(0x82)
	NVME_SC_RESERVATION_CONFLICT = StatusCode(0x83)
	NVME_SC_FORMAT_IN_PROGRESS   = StatusCode(0x84)

	/*
	 * Command Specific Status:
	 */
	NVME_SC_CQ_INVALID                   = StatusCode(0x100)
	NVME_SC_QID_INVALID                  = StatusCode(0x101)
	NVME_SC_QUEUE_SIZE                   = StatusCode(0x102)
	NVME_SC_ABORT_LIMIT                  = StatusCode(0x103)
	NVME_SC_ABORT_MISSING                = StatusCode(0x104)
	NVME_SC_ASYNC_LIMIT                  = StatusCode(0x105)
	NVME_SC_FIRMWARE_SLOT                = StatusCode(0x106)
	NVME_SC_FIRMWARE_IMAGE               = StatusCode(0x107)
	NVME_SC_INVALID_VECTOR               = StatusCode(0x108)
	NVME_SC_INVALID_LOG_PAGE             = StatusCode(0x109)
	NVME_SC_INVALID_FORMAT               = StatusCode(0x10a)
	NVME_SC_FW_NEEDS_CONV_RESET          = StatusCode(0x10b)
	NVME_SC_INVALID_QUEUE                = StatusCode(0x10c)
	NVME_SC_FEATURE_NOT_SAVEABLE         = StatusCode(0x10d)
	NVME_SC_FEATURE_NOT_CHANGEABLE       = StatusCode(0x10e)
	NVME_SC_FEATURE_NOT_PER_NS           = StatusCode(0x10f)
	NVME_SC_FW_NEEDS_SUBSYS_RESET        = StatusCode(0x110)
	NVME_SC_FW_NEEDS_RESET               = StatusCode(0x111)
	NVME_SC_FW_NEEDS_MAXIME              = StatusCode(0x112)
	NVME_SC_FW_ACTIVATE_PROHIBITED       = StatusCode(0x113)
	NVME_SC_OVERLAPPING_RANGE            = StatusCode(0x114)
	NVME_SC_NS_INSUFFICIENT_CAP          = StatusCode(0x115)
	NVME_SC_NS_ID_UNAVAILABLE            = StatusCode(0x116)
	NVME_SC_NS_ALREADY_ATTACHED          = StatusCode(0x118)
	NVME_SC_NS_IS_PRIVATE                = StatusCode(0x119)
	NVME_SC_NS_NOT_ATTACHED              = StatusCode(0x11a)
	NVME_SCHIN_PROV_NOT_SUPP             = StatusCode(0x11b)
	NVME_SC_CTRL_LIST_INVALID            = StatusCode(0x11c)
	NVME_SC_DEVICE_SELFEST_IN_PROGRESS   = StatusCode(0x11d)
	NVME_SC_BP_WRITE_PROHIBITED          = StatusCode(0x11e)
	NVME_SC_INVALID_CTRL_ID              = StatusCode(0x11f)
	NVME_SC_INVALID_SECONDARY_CTRL_STATE = StatusCode(0x120)
	NVME_SC_INVALID_NUM_CTRL_RESOURCE    = StatusCode(0x121)
	NVME_SC_INVALID_RESOURCE_ID          = StatusCode(0x122)
	NVME_SC_PMR_SAN_PROHIBITED           = StatusCode(0x123)
	NVME_SC_ANA_INVALID_GROUP_ID         = StatusCode(0x124)
	NVME_SC_ANA_ATTACH_FAIL              = StatusCode(0x125)

	/*
	 * Command Set Specific - Namespace Types commands:
	 */
	NVME_SC_IOCS_NOT_SUPPORTED        = StatusCode(0x129)
	NVME_SC_IOCS_NOT_ENABLED          = StatusCode(0x12A)
	NVME_SC_IOCS_COMBINATION_REJECTED = StatusCode(0x12B)
	NVME_SC_INVALID_IOCS              = StatusCode(0x12C)

	/*
	 * I/O Command Set Specific - NVM commands:
	 */
	NVME_SC_BAD_ATTRIBUTES          = StatusCode(0x180)
	NVME_SC_INVALID_PI              = StatusCode(0x181)
	NVME_SC_READ_ONLY               = StatusCode(0x182)
	NVME_SC_CMD_SIZE_LIMIT_EXCEEDED = StatusCode(0x183)

	/*
	 * I/O Command Set Specific - Fabrics commands:
	 */
	NVME_SC_CONNECT_FORMAT        = StatusCode(0x180)
	NVME_SC_CONNECT_CTRL_BUSY     = StatusCode(0x181)
	NVME_SC_CONNECT_INVALID_PARAM = StatusCode(0x182)
	NVME_SC_CONNECT_RESTART_DISC  = StatusCode(0x183)
	NVME_SC_CONNECT_INVALID_HOST  = StatusCode(0x184)

	NVME_SC_DISCOVERY_RESTART = StatusCode(0x190)
	NVME_SC_AUTH_REQUIRED     = StatusCode(0x191)

	/*
	 * I/O Command Set Specific - Zoned Namespace commands:
	 */
	NVME_SC_ZONE_BOUNDARY_ERROR         = StatusCode(0x1B8)
	NVME_SC_ZONE_IS_FULL                = StatusCode(0x1B9)
	NVME_SC_ZONE_IS_READ_ONLY           = StatusCode(0x1BA)
	NVME_SC_ZONE_IS_OFFLINE             = StatusCode(0x1BB)
	NVME_SC_ZONE_INVALID_WRITE          = StatusCode(0x1BC)
	NVME_SCOO_MANY_ACTIVE_ZONES         = StatusCode(0x1BD)
	NVME_SCOO_MANY_OPEN_ZONES           = StatusCode(0x1BE)
	NVME_SC_ZONE_INVALID_STATERANSITION = StatusCode(0x1BF)

	/*
	 * Media and Data Integrity Errors:
	 */
	NVME_SC_WRITE_FAULT     = StatusCode(0x280)
	NVME_SC_READ_ERROR      = StatusCode(0x281)
	NVME_SC_GUARD_CHECK     = StatusCode(0x282)
	NVME_SC_APPTAG_CHECK    = StatusCode(0x283)
	NVME_SC_REFTAG_CHECK    = StatusCode(0x284)
	NVME_SC_COMPARE_FAILED  = StatusCode(0x285)
	NVME_SC_ACCESS_DENIED   = StatusCode(0x286)
	NVME_SC_UNWRITTEN_BLOCK = StatusCode(0x287)

	/*
	 * Path-related Errors:
	 */
	NVME_SC_ANA_PERSISTENT_LOSS = StatusCode(0x301)
	NVME_SC_ANA_INACCESSIBLE    = StatusCode(0x302)
	NVME_SC_ANARANSITION        = StatusCode(0x303)

	NVME_SC_CRD = StatusCode(0x1800)
	NVME_SC_DNR = StatusCode(0x4000)
)

const (
	NVME_ADMIN_OP_DELETE_SQ      = AdminOpCode(0x00)
	NVME_ADMIN_OP_CREATE_SQ      = AdminOpCode(0x01)
	NVME_ADMIN_OP_GET_LOG_PAGE   = AdminOpCode(0x02)
	NVME_ADMIN_OP_DELETE_CQ      = AdminOpCode(0x04)
	NVME_ADMIN_OP_CREATE_CQ      = AdminOpCode(0x05)
	NVME_ADMIN_OP_IDENTIFY       = AdminOpCode(0x06)
	NVME_ADMIN_OP_ABORT_CMD      = AdminOpCode(0x08)
	NVME_ADMIN_OP_SET_FEATURES   = AdminOpCode(0x09)
	NVME_ADMIN_OP_GET_FEATURES   = AdminOpCode(0x0A)
	NVME_ADMIN_OP_ASYNC_EVENT    = AdminOpCode(0x0C)
	NVME_ADMIN_OP_NS_MGMT        = AdminOpCode(0x0D)
	NVME_ADMIN_OP_ACTIVATE_FW    = AdminOpCode(0x10)
	NVME_ADMIN_OP_DOWNLOAD_FW    = AdminOpCode(0x11)
	NVME_ADMIN_OP_DEV_SELFEST    = AdminOpCode(0x14)
	NVME_ADMIN_OP_NS_ATTACH      = AdminOpCode(0x15)
	NVME_ADMIN_OP_KEEP_ALIVE     = AdminOpCode(0x18)
	NVME_ADMIN_OP_DIRECTIVE_SEND = AdminOpCode(0x19)
	NVME_ADMIN_OP_DIRECTIVE_RECV = AdminOpCode(0x1A)
	NVME_ADMIN_OP_VIRTUAL_MGMT   = AdminOpCode(0x1C)
	NVME_ADMIN_OP_NVME_MI_SEND   = AdminOpCode(0x1D)
	NVME_ADMIN_OP_NVME_MI_RECV   = AdminOpCode(0x1E)
	NVME_ADMIN_OP_DBBUF          = AdminOpCode(0x7C)
	NVME_ADMIN_OP_FORMAT_NVM     = AdminOpCode(0x80)
	NVME_ADMIN_OP_SECURITY_SEND  = AdminOpCode(0x81)
	NVME_ADMIN_OP_SECURITY_RECV  = AdminOpCode(0x82)
	NVME_ADMIN_OP_SANITIZE_NVM   = AdminOpCode(0x84)
	NVME_ADMIN_OP_GET_LBA_STATUS = AdminOpCode(0x86)
)

/**
 * NVM_Express_Revision_1.3.pdf
 * Figure 86 : Get Log Page - Command Dword 10
 */
const (
	NVME_NO_LOG_LSP      = GetLogPageCommand(0x0)
	NVME_NO_LOG_LPO      = GetLogPageCommand(0x0)
	NVME_LOG_ANA_LSP_RGO = GetLogPageCommand(0x1)
	NVMEELEM_LSP_CREATE  = GetLogPageCommand(0x1)
)

/**
 * NVM_Express_Revision_1.3.pdf
 * Figure 90 : Get Log Page - Log Page Identifiers
 */
const (
	NVME_GET_LOG_PAGE_RESERVED           = GetLogPageIdentifier(0x00)
	NVME_GET_LOG_PAGE_ERROR_INFO         = GetLogPageIdentifier(0x01)
	NVME_GET_LOG_PAGE_SMART              = GetLogPageIdentifier(0x02)
	NVME_GET_LOG_PAGE_FIRMWARE_SLOT_INFO = GetLogPageIdentifier(0x03)
)

/**
 * NVM_Express_Revision_1.3.pdf
 * Figure 178 : Sanitize - Command Dword 10
 */
const (
	NVME_SANITIZE_NO_DEALLOC       = SanitizeCommand(0x00000200)
	NVME_SANITIZE_OIPBP            = SanitizeCommand(0x00000100)
	NVME_SANITIZE_OWPASS_SHIFT     = SanitizeCommand(0x00000004) /* 07:04 */
	NVME_SANITIZE_AUSE             = SanitizeCommand(0x00000008)
	NVME_SANITIZE_ACT_CRYPTO_ERASE = SanitizeCommand(0x00000004)
	NVME_SANITIZE_ACT_OVERWRITE    = SanitizeCommand(0x00000003)
	NVME_SANITIZE_ACT_BLOCK_ERASE  = SanitizeCommand(0x00000002)
	NVME_SANITIZE_ACT_EXIT         = SanitizeCommand(0x00000001)

	/* Sanitize Monitor/Log */
	NVME_SANITIZE_LOG_DATA_LEN              = SanitizeCommand(0x0014)
	NVME_SANITIZE_LOG_GLOBAL_DATA_ERASED    = SanitizeCommand(0x0100)
	NVME_SANITIZE_LOG_NUM_CMPLTED_PASS_MASK = SanitizeCommand(0x00F8)
	NVME_SANITIZE_LOG_STATUS_MASK           = SanitizeCommand(0x0007)
	NVME_SANITIZE_LOG_NEVER_SANITIZED       = SanitizeCommand(0x0000)
	NVME_SANITIZE_LOG_COMPLETED_SUCCESS     = SanitizeCommand(0x0001)
	NVME_SANITIZE_LOG_IN_PROGESS            = SanitizeCommand(0x0002)
	NVME_SANITIZE_LOG_COMPLETED_FAILED      = SanitizeCommand(0x0003)
	NVME_SANITIZE_LOG_ND_COMPLETED_SUCCESS  = SanitizeCommand(0x0004)
)

type UserIo struct {
	Opcode   uint8  `struc:"uint8"`
	Flags    uint8  `struc:"uint8"`
	Control  uint16 `struc:"uint16"`
	Nblocks  uint16 `struc:"uint16"`
	Rsvd     uint16 `struc:"uint16"`
	Metadata uint64 `struc:"uint64"`
	Addr     uint64 `struc:"uint64"`
	Slba     uint64 `struc:"uint64"`
	Dsmgmt   uint32 `struc:"uint32"`
	Reftag   uint32 `struc:"uint32"`
	Apptag   uint16 `struc:"uint16"`
	Appmask  uint16 `struc:"uint16"`
}

type PassthruCmd struct {
	Opcode      uint8  `struc:"uint8"`
	Flags       uint8  `struc:"uint8"`
	Rsvd1       uint16 `struc:"uint16"`
	Nsid        uint32 `struc:"uint32"`
	Cdw2        uint32 `struc:"uint32"`
	Cdw3        uint32 `struc:"uint32"`
	Metadata    uint64 `struc:"uint64"`
	Addr        uint64 `struc:"uint64"`
	MetadataLen uint32 `struc:"uint32"`
	DataLen     uint32 `struc:"uint32"`
	Cdw10       uint32 `struc:"uint32"`
	Cdw11       uint32 `struc:"uint32"`
	Cdw12       uint32 `struc:"uint32"`
	Cdw13       uint32 `struc:"uint32"`
	Cdw14       uint32 `struc:"uint32"`
	Cdw15       uint32 `struc:"uint32"`
	TimeoutMs   uint32 `struc:"uint32"`
	Result      uint32 `struc:"uint32"`
}
type PassthruCmd64 struct {
	Opcode      uint8  `struc:"uint8"`
	Flags       uint8  `struc:"uint8"`
	Rsvd1       uint16 `struc:"uint16"`
	Nsid        uint32 `struc:"uint32"`
	Cdw2        uint32 `struc:"uint32"`
	Cdw3        uint32 `struc:"uint32"`
	Metadata    uint64 `struc:"uint64"`
	Addr        uint64 `struc:"uint64"`
	MetadataLen uint32 `struc:"uint32"`
	DataLen     uint32 `struc:"uint32"`
	Cdw10       uint32 `struc:"uint32"`
	Cdw11       uint32 `struc:"uint32"`
	Cdw12       uint32 `struc:"uint32"`
	Cdw13       uint32 `struc:"uint32"`
	Cdw14       uint32 `struc:"uint32"`
	Cdw15       uint32 `struc:"uint32"`
	TimeoutMs   uint32 `struc:"uint32"`
	Rsvd2       uint32 `struc:"uint32"`
	Result      uint64 `struc:"uint64"`
}

type IdentifyPowerState struct {
	MaxPower/* centiwatts */ uint16          `struc:"uint16"`
	Rsvd2                             uint8  `struc:"uint8"`
	Flags                             uint8  `struc:"uint8"`
	EntryLat/* microseconds */ uint32        `struc:"uint32"`
	ExitLat/* microseconds */ uint32         `struc:"uint32"`
	Readput                           uint8  `struc:"uint8"`
	ReadLat                           uint8  `struc:"uint8"`
	Writeput                          uint8  `struc:"uint8"`
	WriteLat                          uint8  `struc:"uint8"`
	IdlePower                         uint16 `struc:"uint16"`
	IdleScale                         uint8  `struc:"uint8"`
	Rsvd19                            uint8  `struc:"uint8"`
	ActivePower                       uint16 `struc:"uint16"`
	ActiveWorkScale                   uint8  `struc:"uint8"`
	Rsvd23                            uint8  `struc:"[9]uint8"`
}

type IdentifyControllerVersion struct {
	ver uint32 `struc:"uint32"`
	//
	//TertiaryVersion uint8 `struc:"uint8"`
	//MinorVersion uint8 `struc:"uint8"`
	//MajorVersion uint16 `struc:"uint16"`
}

type IdentifyController struct {
	Vid    uint16 `struc:"uint16"`
	Ssvid  uint16 `struc:"uint16"`
	Sn     uint8  `struc:"[20]uint8"`
	Mn     uint8  `struc:"[40]uint8"`
	Fr     uint8  `struc:"[8]uint8"`
	Rab    uint8  `struc:"uint8"`
	Ieee   uint8  `struc:"[3]uint8"`
	Cmic   uint8  `struc:"uint8"`
	Mdts   uint8  `struc:"uint8"`
	Cntlid uint16 `struc:"uint16"`
	IdentifyControllerVersion
	Rtd3r     uint32 `struc:"uint32"`
	Rtd3e     uint32 `struc:"uint32"`
	Oaes      uint32 `struc:"uint32"`
	Ctratt    uint32 `struc:"uint32"`
	Rrls      uint16 `struc:"uint16"`
	Rsvd102   uint8  `struc:"[9]uint8"`
	Cntrltype uint8  `struc:"uint8"`
	Fguid     uint8  `struc:"[16]uint8"`
	Crdt1     uint16 `struc:"uint16"`
	Crdt2     uint16 `struc:"uint16"`
	Crdt3     uint16 `struc:"uint16"`
	Rsvd134   uint8  `struc:"[122]uint8"`
	Oacs      uint16 `struc:"uint16"`
	Acl       uint8  `struc:"uint8"`
	Aerl      uint8  `struc:"uint8"`
	Frmw      uint8  `struc:"uint8"`
	Lpa       uint8  `struc:"uint8"`
	Elpe      uint8  `struc:"uint8"`
	Npss      uint8  `struc:"uint8"`
	Avscc     uint8  `struc:"uint8"`
	Apsta     uint8  `struc:"uint8"`
	Wctemp    uint16 `struc:"uint16"`
	Cctemp    uint16 `struc:"uint16"`
	Mtfa      uint16 `struc:"uint16"`
	Hmpre     uint32 `struc:"uint32"`
	Hmmin     uint32 `struc:"uint32"`
	Tnvmcap   uint8  `struc:"[16]uint8"`
	Unvmcap   uint8  `struc:"[16]uint8"`
	Rpmbs     uint32 `struc:"uint32"`
	Edstt     uint16 `struc:"uint16"`
	Dsto      uint8  `struc:"uint8"`
	Fwug      uint8  `struc:"uint8"`
	Kas       uint16 `struc:"uint16"`
	Hctma     uint16 `struc:"uint16"`
	Mntmt     uint16 `struc:"uint16"`
	Mxtmt     uint16 `struc:"uint16"`
	Sanicap   uint32 `struc:"uint32"`
	Hmminds   uint32 `struc:"uint32"`
	Hmmaxd    uint16 `struc:"uint16"`
	Nsetidmax uint16 `struc:"uint16"`
	Endgidmax uint16 `struc:"uint16"`
	Anatt     uint8  `struc:"uint8"`
	Anacap    uint8  `struc:"uint8"`
	Anagrpmax uint32 `struc:"uint32"`
	Nanagrpid uint32 `struc:"uint32"`
	Pels      uint32 `struc:"uint32"`
	Rsvd356   uint8  `struc:"[156]uint8"`
	Sqes      uint8  `struc:"uint8"`
	Cqes      uint8  `struc:"uint8"`
	Maxcmd    uint16 `struc:"uint16"`
	Nn        uint32 `struc:"uint32"`
	Oncs      uint16 `struc:"uint16"`
	Fuses     uint16 `struc:"uint16"`
	Fna       uint8  `struc:"uint8"`
	Vwc       uint8  `struc:"uint8"`
	Awun      uint16 `struc:"uint16"`
	Awupf     uint16 `struc:"uint16"`
	Icsvscc   uint8  `struc:"uint8"`
	Nwpc      uint8  `struc:"uint8"`
	Acwu      uint16 `struc:"uint16"`
	Ocfs      uint16 `struc:"uint16"`
	Sgls      uint32 `struc:"uint32"`
	Mnan      uint32 `struc:"uint32"`
	Rsvd544   uint8  `struc:"[224]uint8"`
	Subnqn    uint8  `struc:"[256]uint8"`
	Rsvd1024  uint8  `struc:"[768]uint8"`
	Ioccsz    uint32 `struc:"uint32"`
	Iorcsz    uint32 `struc:"uint32"`
	Icdoff    uint16 `struc:"uint16"`
	Ctrattr   uint8  `struc:"uint8"`
	Msdbd     uint8  `struc:"uint8"`
	Rsvd1804  uint8  `struc:"[244]uint8"`
	psd       [32]IdentifyPowerState
	Vs        uint8 `struc:"[1024]uint8"`
}

/**
 * NVM_Express_Revision_1.3.pdf
 * 5.14.1.9.2 Sanitize Status (Log Identifier 81h)
 */
type SanitizeLogPage struct {
	Progress          uint16 `struc:"uint16"`
	Status            uint16 `struc:"uint16"`
	Cdw10Info         uint32 `struc:"uint32"`
	EstOverwriteime   uint32 `struc:"uint32"`
	EstBlockEraseime  uint32 `struc:"uint32"`
	EstCryptoEraseime uint32 `struc:"uint32"`
	// extended
	EstOverwriteimeWithNoDeallocate   uint32 `struc:"uint32"`
	EstBlockEraseimeWithNoDeallocate  uint32 `struc:"uint32"`
	EstCryptoEraseimeWithNoDeallocate uint32 `struc:"uint32"`
}

/**
 * NVM_Express_Revision_1.3.pdf
 * 5.14.1.2 SMART/Health Information (Log Identifier 02h)
 */
type SmartLogPage struct {
	CriticalWarning                           uint8  `struc:"uint8"`
	Compositeemperature                       uint16 `struc:"uint16"`
	AvailableSpare                            uint8  `struc:"uint8"`
	AvailableSparehreshold                    uint8  `struc:"uint8"`
	PercentageUsed                            uint8  `struc:"uint8"`
	Rev01                                     uint8  `struc:"[26]uint8"`
	DataUnitsRead                             uint8  `struc:"[16]uint8"`
	DataUnitsWritten                          uint8  `struc:"[16]uint8"`
	HostReadCommands                          uint8  `struc:"[16]uint8"`
	HostWriteCommands                         uint8  `struc:"[16]uint8"`
	ControllerBusyime                         uint8  `struc:"[16]uint8"`
	PowerCycles                               uint8  `struc:"[16]uint8"`
	PowerOnHours                              uint8  `struc:"[16]uint8"`
	UnsafeShutdowns                           uint8  `struc:"[16]uint8"`
	MediaAndDataIntegrityErrors               uint8  `struc:"[16]uint8"`
	NumberOfErrorInformationLogEntries        uint8  `struc:"[16]uint8"`
	WarningCompositeemperatureime             uint8  `struc:"[4]uint8"`
	CriticalCompositeemperatureime            uint8  `struc:"[4]uint8"`
	TemperatureSensor                         uint16 `struc:"[8]uint16"`
	ThermalManagementemperatureransitionCount uint32 `struc:"[2]uint32"`
	TotalimeForhermalManagementemperature     uint32 `struc:"[2]uint32"`
	RevRemaining                              uint8  `struc:"[280]uint8"`
}
