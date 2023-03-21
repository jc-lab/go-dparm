package plat_linux

type SCSI_PASS_THROUGH_DIRECT struct {
	Length             uint16   `struc:"uint16"`
	ScsiStatus         byte     `struc:"uint8"`
	PathId             byte     `struc:"uint8"`
	TargetId           byte     `struc:"uint8"`
	Lun                byte     `struc:"uint8"`
	CdbLength          byte     `struc:"uint8"`
	SenseInfoLength    byte     `struc:"uint8"`
	DataIn             byte     `struc:"uint8"`
	DataTransferLength uint32   `struc:"uint32"`
	TimeOutValue       uint32   `struc:"uint32"`
	DataBuffer         uintptr  `struc:"off_t"`
	SenseInfoOffset    uint32   `struc:"uint32,offsetof=SenseInfo"`
	Cdb                [16]byte `struc:"[16]uint8"`
}

type SCSI_PASS_THROUGH_DIRECT_WITH_SENSE_BUF struct {
	SCSI_PASS_THROUGH_DIRECT
	Filter    uint32 // realign buffers to double word boundary
	SenseData [32]byte
}

type SCSI_SECURITY_PROTOCOL struct {
	OperationCode uint8
	Protocol      uint8
	ProtocolSp    uint16
	B04           uint8
	B05           uint8
	Length        uint32
	B10           uint8
	Control       uint8
}

func (p *SCSI_SECURITY_PROTOCOL) IsInc512() bool {
	return p.B04&0x80 != 0
}

func (p *SCSI_SECURITY_PROTOCOL) SetInc512(v bool) {
	if v {
		p.B04 |= 0x80
	} else {
		p.B04 &= ^byte(0x80)
	}
}

type SCSI_ADDRESS struct {
	Length     uint32
	PortNumber uint8
	PathId     uint8
	TargetId   uint8
	Lun        uint8
}

const (
	SCSI_IOCTL_DATA_OUT         = byte(0)
	SCSI_IOCTL_DATA_IN          = byte(1)
	SCSI_IOCTL_DATA_UNSPECIFIED = byte(2)
)

//
// Command Descriptor Block constants.
//

const (
	CDB6GENERIC_LENGTH  = 6
	CDB10GENERIC_LENGTH = 10
	CDB12GENERIC_LENGTH = 12
)

//
// Mode Sense/Select page constants.
//

const (
	MODE_PAGE_VENDOR_SPECIFIC               = 0x00
	MODE_PAGE_ERROR_RECOVERY                = 0x01
	MODE_PAGE_DISCONNECT                    = 0x02
	MODE_PAGE_FORMAT_DEVICE                 = 0x03 // disk
	MODE_PAGE_MRW                           = 0x03 // cdrom
	MODE_PAGE_RIGID_GEOMETRY                = 0x04
	MODE_PAGE_FLEXIBILE                     = 0x05 // disk
	MODE_PAGE_WRITE_PARAMETERS              = 0x05 // cdrom
	MODE_PAGE_VERIFY_ERROR                  = 0x07
	MODE_PAGE_CACHING                       = 0x08
	MODE_PAGE_PERIPHERAL                    = 0x09
	MODE_PAGE_CONTROL                       = 0x0A
	MODE_PAGE_MEDIUM_TYPES                  = 0x0B
	MODE_PAGE_NOTCH_PARTITION               = 0x0C
	MODE_PAGE_CD_AUDIO_CONTROL              = 0x0E
	MODE_PAGE_DATA_COMPRESS                 = 0x0F
	MODE_PAGE_DEVICE_CONFIG                 = 0x10
	MODE_PAGE_XOR_CONTROL                   = 0x10 // disk
	MODE_PAGE_MEDIUM_PARTITION              = 0x11
	MODE_PAGE_ENCLOSURE_SERVICES_MANAGEMENT = 0x14
	MODE_PAGE_EXTENDED                      = 0x15
	MODE_PAGE_EXTENDED_DEVICE_SPECIFIC      = 0x16
	MODE_PAGE_CDVD_FEATURE_SET              = 0x18
	MODE_PAGE_PROTOCOL_SPECIFIC_LUN         = 0x18
	MODE_PAGE_PROTOCOL_SPECIFIC_PORT        = 0x19
	MODE_PAGE_POWER_CONDITION               = 0x1A
	MODE_PAGE_LUN_MAPPING                   = 0x1B
	MODE_PAGE_FAULT_REPORTING               = 0x1C
	MODE_PAGE_CDVD_INACTIVITY               = 0x1D // cdrom
	MODE_PAGE_ELEMENT_ADDRESS               = 0x1D
	MODE_PAGE_TRANSPORT_GEOMETRY            = 0x1E
	MODE_PAGE_DEVICE_CAPABILITIES           = 0x1F
	MODE_PAGE_CAPABILITIES                  = 0x2A // cdrom

	MODE_SENSE_RETURN_ALL = 0x3f

	MODE_SENSE_CURRENT_VALUES    = 0x00
	MODE_SENSE_CHANGEABLE_VALUES = 0x40
	MODE_SENSE_DEFAULT_VAULES    = 0x80
	MODE_SENSE_SAVED_VALUES      = 0xc0
)

//
// SCSI CDB operation codes
//

const (
	// 6-byte commands:
	SCSIOP_TEST_UNIT_READY     = 0x00
	SCSIOP_REZERO_UNIT         = 0x01
	SCSIOP_REWIND              = 0x01
	SCSIOP_REQUEST_BLOCK_ADDR  = 0x02
	SCSIOP_REQUEST_SENSE       = 0x03
	SCSIOP_FORMAT_UNIT         = 0x04
	SCSIOP_READ_BLOCK_LIMITS   = 0x05
	SCSIOP_REASSIGN_BLOCKS     = 0x07
	SCSIOP_INIT_ELEMENT_STATUS = 0x07
	SCSIOP_READ6               = 0x08
	SCSIOP_RECEIVE             = 0x08
	SCSIOP_WRITE6              = 0x0A
	SCSIOP_PRINT               = 0x0A
	SCSIOP_SEND                = 0x0A
	SCSIOP_SEEK6               = 0x0B
	SCSIOP_TRACK_SELECT        = 0x0B
	SCSIOP_SLEW_PRINT          = 0x0B
	SCSIOP_SET_CAPACITY        = 0x0B // tape
	SCSIOP_SEEK_BLOCK          = 0x0C
	SCSIOP_PARTITION           = 0x0D
	SCSIOP_READ_REVERSE        = 0x0F
	SCSIOP_WRITE_FILEMARKS     = 0x10
	SCSIOP_FLUSH_BUFFER        = 0x10
	SCSIOP_SPACE               = 0x11
	SCSIOP_INQUIRY             = 0x12
	SCSIOP_VERIFY6             = 0x13
	SCSIOP_RECOVER_BUF_DATA    = 0x14
	SCSIOP_MODE_SELECT         = 0x15
	SCSIOP_RESERVE_UNIT        = 0x16
	SCSIOP_RELEASE_UNIT        = 0x17
	SCSIOP_COPY                = 0x18
	SCSIOP_ERASE               = 0x19
	SCSIOP_MODE_SENSE          = 0x1A
	SCSIOP_START_STOP_UNIT     = 0x1B
	SCSIOP_STOP_PRINT          = 0x1B
	SCSIOP_LOAD_UNLOAD         = 0x1B
	SCSIOP_RECEIVE_DIAGNOSTIC  = 0x1C
	SCSIOP_SEND_DIAGNOSTIC     = 0x1D
	SCSIOP_MEDIUM_REMOVAL      = 0x1E

	// 10-byte commands
	SCSIOP_READ_FORMATTED_CAPACITY = 0x23
	SCSIOP_READ_CAPACITY           = 0x25
	SCSIOP_READ                    = 0x28
	SCSIOP_WRITE                   = 0x2A
	SCSIOP_SEEK                    = 0x2B
	SCSIOP_LOCATE                  = 0x2B
	SCSIOP_POSITION_TO_ELEMENT     = 0x2B
	SCSIOP_WRITE_VERIFY            = 0x2E
	SCSIOP_VERIFY                  = 0x2F
	SCSIOP_SEARCH_DATA_HIGH        = 0x30
	SCSIOP_SEARCH_DATA_EQUAL       = 0x31
	SCSIOP_SEARCH_DATA_LOW         = 0x32
	SCSIOP_SET_LIMITS              = 0x33
	SCSIOP_READ_POSITION           = 0x34
	SCSIOP_SYNCHRONIZE_CACHE       = 0x35
	SCSIOP_COMPARE                 = 0x39
	SCSIOP_COPY_COMPARE            = 0x3A
	SCSIOP_WRITE_DATA_BUFF         = 0x3B
	SCSIOP_READ_DATA_BUFF          = 0x3C
	SCSIOP_WRITE_LONG              = 0x3F
	SCSIOP_CHANGE_DEFINITION       = 0x40
	SCSIOP_WRITE_SAME              = 0x41
	SCSIOP_READ_SUB_CHANNEL        = 0x42
	SCSIOP_UNMAP                   = 0x42 // block device
	SCSIOP_READ_TOC                = 0x43
	SCSIOP_READ_HEADER             = 0x44
	SCSIOP_REPORT_DENSITY_SUPPORT  = 0x44 // tape
	SCSIOP_PLAY_AUDIO              = 0x45
	SCSIOP_GET_CONFIGURATION       = 0x46
	SCSIOP_PLAY_AUDIO_MSF          = 0x47
	SCSIOP_PLAY_TRACK_INDEX        = 0x48
	SCSIOP_SANITIZE                = 0x48 // block device
	SCSIOP_PLAY_TRACK_RELATIVE     = 0x49
	SCSIOP_GET_EVENT_STATUS        = 0x4A
	SCSIOP_PAUSE_RESUME            = 0x4B
	SCSIOP_LOG_SELECT              = 0x4C
	SCSIOP_LOG_SENSE               = 0x4D
	SCSIOP_STOP_PLAY_SCAN          = 0x4E
	SCSIOP_XDWRITE                 = 0x50
	SCSIOP_XPWRITE                 = 0x51
	SCSIOP_READ_DISK_INFORMATION   = 0x51
	SCSIOP_READ_DISC_INFORMATION   = 0x51 // proper use of disc over disk
	SCSIOP_READ_TRACK_INFORMATION  = 0x52
	SCSIOP_XDWRITE_READ            = 0x53
	SCSIOP_RESERVE_TRACK_RZONE     = 0x53
	SCSIOP_SEND_OPC_INFORMATION    = 0x54 // optimum power calibration
	SCSIOP_MODE_SELECT10           = 0x55
	SCSIOP_RESERVE_UNIT10          = 0x56
	SCSIOP_RESERVE_ELEMENT         = 0x56
	SCSIOP_RELEASE_UNIT10          = 0x57
	SCSIOP_RELEASE_ELEMENT         = 0x57
	SCSIOP_REPAIR_TRACK            = 0x58
	SCSIOP_MODE_SENSE10            = 0x5A
	SCSIOP_CLOSE_TRACK_SESSION     = 0x5B
	SCSIOP_READ_BUFFER_CAPACITY    = 0x5C
	SCSIOP_SEND_CUE_SHEET          = 0x5D
	SCSIOP_PERSISTENT_RESERVE_IN   = 0x5E
	SCSIOP_PERSISTENT_RESERVE_OUT  = 0x5F

	// 12-byte commands
	SCSIOP_REPORT_LUNS                  = 0xA0
	SCSIOP_BLANK                        = 0xA1
	SCSIOP_ATA_PASSTHROUGH12            = 0xA1
	SCSIOP_SEND_EVENT                   = 0xA2
	SCSIOP_SECURITY_PROTOCOL_IN         = 0xA2
	SCSIOP_SEND_KEY                     = 0xA3
	SCSIOP_MAINTENANCE_IN               = 0xA3
	SCSIOP_REPORT_KEY                   = 0xA4
	SCSIOP_MAINTENANCE_OUT              = 0xA4
	SCSIOP_MOVE_MEDIUM                  = 0xA5
	SCSIOP_LOAD_UNLOAD_SLOT             = 0xA6
	SCSIOP_EXCHANGE_MEDIUM              = 0xA6
	SCSIOP_SET_READ_AHEAD               = 0xA7
	SCSIOP_MOVE_MEDIUM_ATTACHED         = 0xA7
	SCSIOP_READ12                       = 0xA8
	SCSIOP_GET_MESSAGE                  = 0xA8
	SCSIOP_SERVICE_ACTION_OUT12         = 0xA9
	SCSIOP_WRITE12                      = 0xAA
	SCSIOP_SEND_MESSAGE                 = 0xAB
	SCSIOP_SERVICE_ACTION_IN12          = 0xAB
	SCSIOP_GET_PERFORMANCE              = 0xAC
	SCSIOP_READ_DVD_STRUCTURE           = 0xAD
	SCSIOP_WRITE_VERIFY12               = 0xAE
	SCSIOP_VERIFY12                     = 0xAF
	SCSIOP_SEARCH_DATA_HIGH12           = 0xB0
	SCSIOP_SEARCH_DATA_EQUAL12          = 0xB1
	SCSIOP_SEARCH_DATA_LOW12            = 0xB2
	SCSIOP_SET_LIMITS12                 = 0xB3
	SCSIOP_READ_ELEMENT_STATUS_ATTACHED = 0xB4
	SCSIOP_REQUEST_VOL_ELEMENT          = 0xB5
	SCSIOP_SECURITY_PROTOCOL_OUT        = 0xB5
	SCSIOP_SEND_VOLUME_TAG              = 0xB6
	SCSIOP_SET_STREAMING                = 0xB6 // C/DVD
	SCSIOP_READ_DEFECT_DATA             = 0xB7
	SCSIOP_READ_ELEMENT_STATUS          = 0xB8
	SCSIOP_READ_CD_MSF                  = 0xB9
	SCSIOP_SCAN_CD                      = 0xBA
	SCSIOP_REDUNDANCY_GROUP_IN          = 0xBA
	SCSIOP_SET_CD_SPEED                 = 0xBB
	SCSIOP_REDUNDANCY_GROUP_OUT         = 0xBB
	SCSIOP_PLAY_CD                      = 0xBC
	SCSIOP_SPARE_IN                     = 0xBC
	SCSIOP_MECHANISM_STATUS             = 0xBD
	SCSIOP_SPARE_OUT                    = 0xBD
	SCSIOP_READ_CD                      = 0xBE
	SCSIOP_VOLUME_SET_IN                = 0xBE
	SCSIOP_SEND_DVD_STRUCTURE           = 0xBF
	SCSIOP_VOLUME_SET_OUT               = 0xBF
	SCSIOP_INIT_ELEMENT_RANGE           = 0xE7

	// 16-byte commands
	SCSIOP_XDWRITE_EXTENDED16            = 0x80 // disk
	SCSIOP_WRITE_FILEMARKS16             = 0x80 // tape
	SCSIOP_REBUILD16                     = 0x81 // disk
	SCSIOP_READ_REVERSE16                = 0x81 // tape
	SCSIOP_REGENERATE16                  = 0x82 // disk
	SCSIOP_EXTENDED_COPY                 = 0x83
	SCSIOP_POPULATE_TOKEN                = 0x83 // disk
	SCSIOP_WRITE_USING_TOKEN             = 0x83 // disk
	SCSIOP_RECEIVE_COPY_RESULTS          = 0x84
	SCSIOP_RECEIVE_ROD_TOKEN_INFORMATION = 0x84 //disk
	SCSIOP_ATA_PASSTHROUGH16             = 0x85
	SCSIOP_ACCESS_CONTROL_IN             = 0x86
	SCSIOP_ACCESS_CONTROL_OUT            = 0x87
	SCSIOP_READ16                        = 0x88
	SCSIOP_COMPARE_AND_WRITE             = 0x89
	SCSIOP_WRITE16                       = 0x8A
	SCSIOP_READ_ATTRIBUTES               = 0x8C
	SCSIOP_WRITE_ATTRIBUTES              = 0x8D
	SCSIOP_WRITE_VERIFY16                = 0x8E
	SCSIOP_VERIFY16                      = 0x8F
	SCSIOP_PREFETCH16                    = 0x90
	SCSIOP_SYNCHRONIZE_CACHE16           = 0x91
	SCSIOP_SPACE16                       = 0x91 // tape
	SCSIOP_LOCK_UNLOCK_CACHE16           = 0x92
	SCSIOP_LOCATE16                      = 0x92 // tape
	SCSIOP_WRITE_SAME16                  = 0x93
	SCSIOP_ERASE16                       = 0x93 // tape
	SCSIOP_READ_CAPACITY16               = 0x9E
	SCSIOP_GET_LBA_STATUS                = 0x9E
	SCSIOP_SERVICE_ACTION_IN16           = 0x9E
	SCSIOP_SERVICE_ACTION_OUT16          = 0x9F

	// 32-byte commands
	SCSIOP_OPERATION32 = 0x7F
)

const (
	// Service Action for 32 bit write commands
	SERVICE_ACTION_XDWRITE      = 0x0004
	SERVICE_ACTION_XPWRITE      = 0x0006
	SERVICE_ACTION_XDWRITEREAD  = 0x0007
	SERVICE_ACTION_WRITE        = 0x000B
	SERVICE_ACTION_WRITE_VERIFY = 0x000C
	SERVICE_ACTION_WRITE_SAME   = 0x000D
	SERVICE_ACTION_ORWRITE      = 0x000E

	// Service actions for 0x48
	SERVICE_ACTION_OVERWRITE    = 0x01
	SERVICE_ACTION_BLOCK_ERASE  = 0x02
	SERVICE_ACTION_CRYPTO_ERASE = 0x03
	SERVICE_ACTION_EXIT_FAILURE = 0x1f

	// Service actions for 0x83
	SERVICE_ACTION_POPULATE_TOKEN    = 0x10
	SERVICE_ACTION_WRITE_USING_TOKEN = 0x11

	// Service actions for 0x84
	SERVICE_ACTION_RECEIVE_TOKEN_INFORMATION = 0x07

	// Service actions for 0x9E
	SERVICE_ACTION_READ_CAPACITY16 = 0x10
	SERVICE_ACTION_GET_LBA_STATUS  = 0x12
)

//
// If the IMMED bit is 1, status is returned as soon
// as the operation is initiated. If the IMMED bit
// is 0, status is not returned until the operation
// is completed.
//

const (
	CDB_RETURN_ON_COMPLETION = 0
	CDB_RETURN_IMMEDIATE     = 1
)

// end_ntminitape

//
// CDB Force media access used in extended read and write commands.
//

const (
	CDB_FORCE_MEDIA_ACCESS = 0x08
)

//
// Denon CD ROM operation codes
//
const (
	SCSIOP_DENON_EJECT_DISC   = 0xE6
	SCSIOP_DENON_STOP_AUDIO   = 0xE7
	SCSIOP_DENON_PLAY_AUDIO   = 0xE8
	SCSIOP_DENON_READ_TOC     = 0xE9
	SCSIOP_DENON_READ_SUBCODE = 0xEB
)

//
// SCSI Bus Messages
//
const (
	SCSIMESS_ABORT                = 0x06
	SCSIMESS_ABORT_WITH_TAG       = 0x0D
	SCSIMESS_BUS_DEVICE_RESET     = 0x0C
	SCSIMESS_CLEAR_QUEUE          = 0x0E
	SCSIMESS_COMMAND_COMPLETE     = 0x00
	SCSIMESS_DISCONNECT           = 0x04
	SCSIMESS_EXTENDED_MESSAGE     = 0x01
	SCSIMESS_IDENTIFY             = 0x80
	SCSIMESS_IDENTIFY_WITH_DISCON = 0xC0
	SCSIMESS_IGNORE_WIDE_RESIDUE  = 0x23
	SCSIMESS_INITIATE_RECOVERY    = 0x0F
	SCSIMESS_INIT_DETECTED_ERROR  = 0x05
	SCSIMESS_LINK_CMD_COMP        = 0x0A
	SCSIMESS_LINK_CMD_COMP_W_FLAG = 0x0B
	SCSIMESS_MESS_PARITY_ERROR    = 0x09
	SCSIMESS_MESSAGE_REJECT       = 0x07
	SCSIMESS_NO_OPERATION         = 0x08
	SCSIMESS_HEAD_OF_QUEUE_TAG    = 0x21
	SCSIMESS_ORDERED_QUEUE_TAG    = 0x22
	SCSIMESS_SIMPLE_QUEUE_TAG     = 0x20
	SCSIMESS_RELEASE_RECOVERY     = 0x10
	SCSIMESS_RESTORE_POINTERS     = 0x03
	SCSIMESS_SAVE_DATA_POINTER    = 0x02
	SCSIMESS_TERMINATE_IO_PROCESS = 0x11
)

//
// SCSI Extended Message operation codes
//

const (
	SCSIMESS_MODIFY_DATA_POINTER  = 0x00
	SCSIMESS_SYNCHRONOUS_DATA_REQ = 0x01
	SCSIMESS_WIDE_DATA_REQUEST    = 0x03
)

//
// SCSI Extended Message Lengths
//

const (
	SCSIMESS_MODIFY_DATA_LENGTH = 5
	SCSIMESS_SYNCH_DATA_LENGTH  = 3
	SCSIMESS_WIDE_DATA_LENGTH   = 2
)
