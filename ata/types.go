package ata

type OpCode uint8

const (
	ATA_OP_DSM                      = OpCode(0x06) // Data Set Management (TRIM)
	ATA_OP_READ_PIO                 = OpCode(0x20)
	ATA_OP_READ_PIO_ONCE            = OpCode(0x21)
	ATA_OP_READ_LONG                = OpCode(0x22)
	ATA_OP_READ_LONG_ONCE           = OpCode(0x23)
	ATA_OP_READ_PIO_EXT             = OpCode(0x24)
	ATA_OP_READ_DMA_EXT             = OpCode(0x25)
	ATA_OP_READ_LOG_EXT             = OpCode(0x2f)
	ATA_OP_READ_FPDMA               = OpCode(0x60) // NCQ
	ATA_OP_WRITE_PIO                = OpCode(0x30)
	ATA_OP_WRITE_LONG               = OpCode(0x32)
	ATA_OP_WRITE_LONG_ONCE          = OpCode(0x33)
	ATA_OP_WRITE_PIO_EXT            = OpCode(0x34)
	ATA_OP_WRITE_DMA_EXT            = OpCode(0x35)
	ATA_OP_WRITE_FPDMA              = OpCode(0x61) // NCQ
	ATA_OP_READ_VERIFY              = OpCode(0x40)
	ATA_OP_READ_VERIFY_ONCE         = OpCode(0x41)
	ATA_OP_READ_VERIFY_EXT          = OpCode(0x42)
	ATA_OP_WRITE_UNC_EXT            = OpCode(0x45) // lba48, no data, uses feat reg
	ATA_OP_FORMAT_TRACK             = OpCode(0x50)
	ATA_OP_TRUSTED_RECV             = OpCode(0x5c)
	ATA_OP_TRUSTED_RECV_DMA         = OpCode(0x5d)
	ATA_OP_TRUSTED_SEND             = OpCode(0x5e)
	ATA_OP_TRUSTED_SEND_DMA         = OpCode(0x5f)
	ATA_OP_DOWNLOAD_MICROCODE       = OpCode(0x92)
	ATA_OP_STANDBYNOW2              = OpCode(0x94)
	ATA_OP_CHECKPOWERMODE2          = OpCode(0x98)
	ATA_OP_SLEEPNOW2                = OpCode(0x99)
	ATA_OP_PIDENTIFY                = OpCode(0xa1)
	ATA_OP_READ_NATIVE_MAX          = OpCode(0xf8)
	ATA_OP_READ_NATIVE_MAX_EXT      = OpCode(0x27)
	ATA_OP_GET_NATIVE_MAX_EXT       = OpCode(0x78)
	ATA_OP_SMART                    = OpCode(0xb0)
	ATA_OP_DCO                      = OpCode(0xb1)
	ATA_OP_SET_SECTOR_CONFIGURATION = OpCode(0xb2)
	ATA_OP_SANITIZE                 = OpCode(0xb4)
	ATA_OP_ERASE_SECTORS            = OpCode(0xc0)
	ATA_OP_READ_DMA                 = OpCode(0xc8)
	ATA_OP_WRITE_DMA                = OpCode(0xca)
	ATA_OP_DOORLOCK                 = OpCode(0xde)
	ATA_OP_DOORUNLOCK               = OpCode(0xdf)
	ATA_OP_STANDBYNOW1              = OpCode(0xe0)
	ATA_OP_IDLEIMMEDIATE            = OpCode(0xe1)
	ATA_OP_SETIDLE                  = OpCode(0xe3)
	ATA_OP_SET_MAX                  = OpCode(0xf9)
	ATA_OP_SET_MAX_EXT              = OpCode(0x37)
	ATA_OP_SET_MULTIPLE             = OpCode(0xc6)
	ATA_OP_CHECKPOWERMODE1          = OpCode(0xe5)
	ATA_OP_SLEEPNOW1                = OpCode(0xe6)
	ATA_OP_FLUSHCACHE               = OpCode(0xe7)
	ATA_OP_FLUSHCACHE_EXT           = OpCode(0xea)
	ATA_OP_IDENTIFY                 = OpCode(0xec)
	ATA_OP_SETFEATURES              = OpCode(0xef)
	ATA_OP_SECURITY_SET_PASS        = OpCode(0xf1)
	ATA_OP_SECURITY_UNLOCK          = OpCode(0xf2)
	ATA_OP_SECURITY_ERASE_PREPARE   = OpCode(0xf3)
	ATA_OP_SECURITY_ERASE_UNIT      = OpCode(0xf4)
	ATA_OP_SECURITY_FREEZE_LOCK     = OpCode(0xf5)
	ATA_OP_SECURITY_DISABLE         = OpCode(0xf6)
	ATA_OP_VENDOR_SPECIFIC_0x80     = OpCode(0x80)
)

/*
 * Some useful ATA register bits
 */
const (
	ATA_USING_LBA = (1 << 6)
	ATA_STAT_DRQ  = (1 << 3)
	ATA_STAT_ERR  = (1 << 0)
)

/*
 * Useful parameters for initHdioTaskfile():
 */
const (
	RW_READ     = 0
	RW_WRITE    = 1
	LBA28_OK    = 0
	LBA48_FORCE = 1
)

/**
 * Working Draft ATA Command Set - 4 (ACS-4)
 * 7.44 SMART Table 133 - FEATURE field values
 *
 * SFF-8035R2
 */
const (
	SMART_FEAT_READ_ATTRIBUTE_VALUES     = 0xd0
	SMART_FEAT_READ_ATTRIBUTE_THRESHOLDS = 0xd1
	SMART_FEAT_EXECUTE_OFFLINE_IMMEDIATE = 0xd4
	SMART_FEAT_READ_LOG                  = 0xd5
	SMART_FEAT_WRITE_LOG                 = 0xd6
	SMART_FEAT_RETURN_STATUS             = 0xda
)

const (
	SMART_LBA_HIGH = 0xc2
	SMART_LBA_LOW  = 0x4f
)

const (
	SMART_RETURN_STATUS_HI_EXCEEDED  = 0x2c
	SMART_RETURN_STATUS_MID_EXCEEDED = 0xf4
)

/*
 * Definitions and structures for use with SGIO + ATA16:
 */
type LbaRegs struct {
	Feat  uint8 `struc:"uint8"`
	Nsect uint8 `struc:"uint8"`
	Lbal  uint8 `struc:"uint8"` // 16bit able with hob
	Lbam  uint8 `struc:"uint8"` // 16bit able with hob
	Lbah  uint8 `struc:"uint8"` // 16bit able with hob
}

type Tf struct {
	Dev     uint8  `struc:"uint8"`
	Command OpCode `struc:"uint8"`
	Error   uint8  `struc:"uint8"`
	Status  uint8  `struc:"uint8"`
	IsLba48 uint8  `struc:"uint8"`
	Lob     LbaRegs
	Hob     LbaRegs
}

type IdentityGeneralConfiguration struct {
	//    uint16T reserved1 : 1;
	//    uint16T retired3 : 1;
	//    uint16T responseIncomplete : 1;
	//    uint16T retired2 : 3;
	//    uint16T fixedDevice : 1;   // obsolete
	//    uint16T removableMedia : 1;// obsolete
	//    uint16T retired1 : 7;
	//    uint16T deviceType : 1;
	A uint8 `struc:"uint8"`
	B uint8 `struc:"uint8"`
}

type IdentityTrustedComputing struct {
	A uint16 `struc:"uint16"`
}

type IdentityCapabilities struct {
	A uint16 `struc:"uint16"`
	B uint16 `struc:"uint16"`
}

type IdentityWord53 struct {
	A                          uint8 `struc:"uint8"`
	FreeFallControlSensitivity uint8 `struc:"uint8"`
}

type IdentityWord59 struct {
	CurrentMultiSectorSetting uint8 `struc:"uint8"`
	B                         uint8 `struc:"uint8"`

	//
	//MultiSectorSettingValid: 1 uint8 `struc:"uint8"`
	//ReservedByte59: 3 uint8 `struc:"uint8"`
	//SanitizeFeatureSupported: 1 uint8 `struc:"uint8"`
	//CryptoScrambleExtCommandSupported: 1 uint8 `struc:"uint8"`
	//OverwriteExtCommandSupported: 1 uint8 `struc:"uint8"`
	//BlockEraseExtCommandSupported: 1 uint8 `struc:"uint8"`
}

type IdentityAdditionalSupported struct {
	A uint16 `struc:"uint16"`
}

type IdentityWord75 struct {
	A uint16 `struc:"uint16"`
	//QueueDepth: 5 uint16 `struc:"uint16"` //  Maximum queue depth - 1
	//ReservedWord75: 11 uint16 `struc:"uint16"`
}

type IdentitySerialAtaCapabilities struct {
	A uint16 `struc:"uint16"`
	B uint16 `struc:"uint16"`
}
type IdentitySerialAtaFeaturesSupported struct {
	A uint16 `struc:"uint16"`
}
type IdentitySerialAtaFeaturesEnabled struct {
	A uint16 `struc:"uint16"`
}
type IdentityCommandSetSupport struct {
	A uint16 `struc:"uint16"`
	B uint16 `struc:"uint16"`
	C uint16 `struc:"uint16"`
}

func (p *IdentityCommandSetSupport) GetSmartCommands() bool {
	return (p.A & 0x0001) != 0
}

type IdentityCommandSetActive struct {
	A uint16 `struc:"uint16"`
	B uint16 `struc:"uint16"`
	C uint16 `struc:"uint16"`
}

func (p *IdentityCommandSetActive) GetSmartCommands() bool {
	return (p.A & 0x0001) != 0
}

type IdentityNormalSecurityEraseUnit struct {
	A uint16 `struc:"uint16"`
}
type IdentityEnhancedSecurityEraseUnit struct {
	A uint16 `struc:"uint16"`
}
type IdentityPhysicalLogicalSectorSize struct {
	A uint16 `struc:"uint16"`
}
type IdentityCommandSupportActiveExt struct {
	A uint16 `struc:"uint16"`
}
type IdentityCommandSetSupportExt struct {
	A uint16 `struc:"uint16"`
}
type IdentityCommandSetActiveExt struct {
	A uint16 `struc:"uint16"`
}

type IdentityWord127 struct {
	A uint16 `struc:"uint16"`
}

type IdentitySecurityStatus struct {
	A uint16 `struc:"uint16"`
}
type IdentityCfgPowerMode1 struct {
	A uint16 `struc:"uint16"`
}

type IdentityWord168 struct {
	A uint16 `struc:"uint16"`
}

type IdentityDataSetManagementFeature struct {
	A uint16 `struc:"uint16"`
}

func (f *IdentityDataSetManagementFeature) GetTrim() bool {
	return (f.A & 0x01) != 0
}

type IdentitySctSommandTransport struct {
	A uint16 `struc:"uint16"`
}
type IdentityBlockAlignment struct {
	A uint16 `struc:"uint16"`
}
type IdentityNvCacheCapabilities struct {
	A uint16 `struc:"uint16"`
}
type IdentityNvCacheOptions struct {
	A uint16 `struc:"uint16"`
}
type IdentityTransportMajorVersion struct {
	A uint16 `struc:"uint16"`
}

type IdentityDeviceData struct {
	GeneralConfiguration  IdentityGeneralConfiguration
	NumCylinders          uint16    `struc:"uint16"` // word 1 obsolete
	SpecificConfiguration uint16    `struc:"uint16"` // word 2
	NumHeads              uint16    `struc:"uint16"` // word 3 obsolete
	Retired1              [2]uint16 `struc:"[2]uint16"`
	NumSectorsPerTrack    uint16    `struc:"uint16"` // word 6 obsolete
	VendorUnique1         [3]uint16 `struc:"[3]uint16"`
	SerialNumber          [20]uint8 `struc:"[20]uint8"` // word 10-19
	Retired2              [2]uint16 `struc:"[2]uint16"`
	Obsolete1             uint16    `struc:"uint16"`
	FirmwareRevision      [8]uint8  `struc:"[8]uint8"`  // word 23-26
	ModelNumber           [40]uint8 `struc:"[40]uint8"` // word 27-46
	MaximumBlockTransfer  uint8     `struc:"uint8"`     // word 47. 01h-10h = Maximum number of sectors that shall be transferred per interrupt on READ/WRITE MULTIPLE commands
	VendorUnique2         uint8     `struc:"uint8"`

	TrustComputing IdentityTrustedComputing // word 48

	Capabilities IdentityCapabilities // word 49-50

	ObsoleteWords51 [2]uint16 `struc:"[2]uint16"`

	Word53 IdentityWord53

	NumberOfCurrentCylinders uint16 `struc:"uint16"` // word 54 obsolete
	NumberOfCurrentHeads     uint16 `struc:"uint16"` // word 55 obsolete
	CurrentSectorsPerTrack   uint16 `struc:"uint16"` // word 56 obsolete
	CurrentSectorCapacity    uint32 `struc:"uint32"` // word 57 word 58 obsolete

	Word59 IdentityWord59

	UserAddressableSectors uint32 `struc:"uint32"` // word 60-61 for 28-bit commands

	ObsoleteWord62 uint16 `struc:"uint16"`

	MultiWordDmaSupport uint8 `struc:"uint8"` // word 63
	MultiWordDmaActive  uint8 `struc:"uint8"`

	AdvancedPioModes uint8 `struc:"uint8"` // word 64. bit 0:1 - PIO mode supported
	ReservedByte64   uint8 `struc:"uint8"`

	MinimumMwxferCycleTime     uint16                      `struc:"uint16"` // word 65
	RecommendedMwxferCycleTime uint16                      `struc:"uint16"` // word 66
	MinimumPioCycleTime        uint16                      `struc:"uint16"` // word 67
	MinimumPioCycleTimeIordy   uint16                      `struc:"uint16"` // word 68
	AdditionalSupported        IdentityAdditionalSupported // word 69

	ReservedWords70 [5]uint16 `struc:"[5]uint16"` // word 70 - reserved
	// word 71:74 - Reserved for the IDENTIFY PACKET DEVICE command

	Word75 IdentityWord75

	SerialAtaCapabilities IdentitySerialAtaCapabilities

	// Word 78
	SerialAtaFeaturesSupported IdentitySerialAtaFeaturesSupported

	// Word 79
	SerialAtaFeaturesEnabled IdentitySerialAtaFeaturesEnabled
	MajorRevision            uint16 `struc:"uint16"` // word 80. bit 5 - supports ATA5; bit 6 - supports ATA6; bit 7 - supports ATA7; bit 8 - supports ATA8-ACS; bit 9 - supports ACS-2;
	MinorRevision            uint16 `struc:"uint16"` // word 81. T13 minior version number

	CommandSetSupport IdentityCommandSetSupport
	CommandSetActive  IdentityCommandSetActive

	UltraDmaSupport uint8 `struc:"uint8"` // word 88. bit 0 - UDMA mode 0 is supported ... bit 6 - UDMA mode 6 and below are supported
	UltraDmaActive  uint8 `struc:"uint8"` // word 88. bit 8 - UDMA mode 0 is selected ... bit 14 - UDMA mode 6 is selected

	NormalSecurityEraseUnit IdentityNormalSecurityEraseUnit // word 89

	EnhancedSecurityEraseUnit IdentityEnhancedSecurityEraseUnit // word 90

	CurrentApmLevel uint8 `struc:"uint8"` // word 91
	ReservedWord91  uint8 `struc:"uint8"`

	MasterPasswordId uint16 `struc:"uint16"` // word 92. Master Password Identifier

	HardwareResetResult uint16 `struc:"uint16"` // word 93

	CurrentAcousticValue     uint8 `struc:"uint8"` // word 94. obsolete
	RecommendedAcousticValue uint8 `struc:"uint8"`

	StreamMinRequestSize         uint16 `struc:"uint16"` // word 95
	StreamingTransferTimeDma     uint16 `struc:"uint16"` // word 96
	StreamingAccessLatencyDmaPio uint16 `struc:"uint16"` // word 97
	StreamingPerfGranularity     uint32 `struc:"uint32"` // word 98 99

	Max48bitLba [2]uint32 `struc:"[2]uint32"` // word 100-103

	StreamingTransferTime uint16 `struc:"uint16"` // word 104. Streaming Transfer Time - PIO

	DsmCap uint16 `struc:"uint16"` // word 105

	PhysicalLogicalSectorSize IdentityPhysicalLogicalSectorSize // word 106

	InterSeekDelay                uint16    `struc:"uint16"`    //word 107.     Inter-seek delay for ISO 7779 standard acoustic testing
	WorldWideName                 [4]uint16 `struc:"[4]uint16"` //words 108-111
	ReservedForWorldWideName128   [4]uint16 `struc:"[4]uint16"` //words 112-115
	ReservedForTlcTechnicalReport uint16    `struc:"uint16"`    //word 116
	WordsPerLogicalSector         [2]uint16 `struc:"[2]uint16"` //words 117-118 Logical sector size (DWord)

	CommandSetSupportExt IdentityCommandSetSupportExt //word 119
	CommandSetActiveExt  IdentityCommandSetActiveExt  //word 120

	ReservedForExpandedSupportandActive [6]uint16 `struc:"[6]uint16"`

	Word127 IdentityWord127

	SecurityStatus IdentitySecurityStatus

	ReservedWord129 [31]uint16 `struc:"[31]uint16"` //word 129...159. Vendor specific

	CfaPowerMode1         IdentityCfgPowerMode1
	ReservedForCfaWord161 [7]uint16 `struc:"[7]uint16"` //Words 161-167

	Word168                  IdentityWord168 //Word 168
	DataSetManagementFeature IdentityDataSetManagementFeature

	AdditionalProductID [4]uint16 `struc:"[4]uint16"` //Words 170-173

	ReservedForCfaWord174 [2]uint16 `struc:"[2]uint16"` //Words 174-175

	CurrentMediaSerialNumber [30]uint16 `struc:"[30]uint16"` //Words 176-205

	SctCommandTransport IdentitySctSommandTransport //Words 206
	ReservedWord207     [2]uint16                   `struc:"[2]uint16"` //Words 207-208

	BlockAlignment IdentityBlockAlignment //Word 209

	WriteReadVerifySectorCountMode3Only [2]uint16 `struc:"[2]uint16"` //Words 210-211
	WriteReadVerifySectorCountMode2Only [2]uint16 `struc:"[2]uint16"` //Words 212-213

	NvCacheCapabilities IdentityNvCacheCapabilities //Word 214. obsolete
	NvCacheSizeLSW      uint16                      `struc:"uint16"` //Word 215. obsolete
	NvCacheSizeMSW      uint16                      `struc:"uint16"` //Word 216. obsolete

	NominalMediaRotationRate uint16 `struc:"uint16"` //Word 217; value 0001h means non-rotating media.

	ReservedWord218 uint16 `struc:"uint16"` //Word 218

	NvCacheOptions IdentityNvCacheOptions //Word 219. obsolete

	WriteReadVerifySectorCountMode uint8 `struc:"uint8"` //Word 220. Write-Read-Verify feature set current mode
	ReservedWord220                uint8 `struc:"uint8"`

	ReservedWord221 uint16 `struc:"uint16"` //Word 221

	TransportMajorVersion IdentityTransportMajorVersion //Word 222 Transport major version number
	TransportMinorVersion uint16                        `struc:"uint16"` // Word 223

	ReservedWord224 [6]uint16 `struc:"[6]uint16"` // Word 224...229

	ExtendedNumberOfUserAddressableSectors [2]uint32 `struc:"[2]uint32"` // Words 230...233 Extended Number of User Addressable Sectors

	MinBlocksPerDownloadMicrocodeMode03 uint16 `struc:"uint16"` // Word 234 Minimum number of 512-byte data blocks per Download Microcode mode 03h operation
	MaxBlocksPerDownloadMicrocodeMode03 uint16 `struc:"uint16"` // Word 235 Maximum number of 512-byte data blocks per Download Microcode mode 03h operation

	ReservedWord236 [19]uint16 `struc:"[19]uint16"` // Word 236...254

	Signature uint8 `struc:"uint8"` //Word 255
	CheckSum  uint8 `struc:"uint8"`
}

const (
	SMART_ATTRIBUTES_NUMBER = 30
)

// 12
type SmartAttribute struct {
	Id uint8 `struc:"uint8"`
	/**
	 * Status flag
	 * Bit 0 (pre-failure/advisory bit)
	 * Bit 1 (on-line data collection bit)
	 * Bits 2-5 (vendor specific)
	 * Bits 6-15 (Reserved)
	 */
	Flags      uint16   `struc:"uint16"`
	Current    uint8    `struc:"uint8"`
	Worst      uint8    `struc:"uint8"`
	Raw        [6]uint8 `struc:"[6]uint8"`
	Reserved01 uint8    `struc:"uint8"`
}

type SmartAttributeThreshold struct {
	Id         uint8     `struc:"uint8"`
	Threshold  uint8     `struc:"uint8"`
	Reserved01 [10]uint8 `struc:"[10]uint8"`
}

/**
 * ATA8-ACS
 * Table 49 — Device SMART data structure
 *
 * SFF-8035R2
 */
type SmartAttributeValues struct {
	RevNumber                                        uint16 `struc:"uint16"`
	Attributes                                       [SMART_ATTRIBUTES_NUMBER]SmartAttribute
	OfflineDataCollectionStatus                      uint8  `struc:"uint8"`
	SelfTestExecStatus                               uint8  `struc:"uint8"`
	TotalTimeToCompleteOfflineDataCollectionActivity uint16 `struc:"uint16"`
	VendorSpecific366                                uint8  `struc:"uint8"`
	OfflineDataCollectionCapability                  uint8  `struc:"uint8"`
	SmartCapability                                  uint16 `struc:"uint16"`
	/**
	 * 7-1: Reserved
	 *   0: device error logging supported
	 */
	ErrorlogCapability uint8 `struc:"uint8"`
	VendorSpecific371  uint8 `struc:"uint8"`
	/**
	 * Unit: minutes
	 */
	ShortSelftestRoutineRecommandedPollingTime uint8 `struc:"uint8"`
	/**
	 * Extended self-test routine recommended polling time (7:0) in minutes.
	 * If FFh use bytes 375 and 376 for the polling time.
	 */
	ExtendedSelftestRoutineRecommandedPollingTimeA  uint8      `struc:"uint8"`
	ConveyanceSelftestRoutineRecommandedPollingTime uint8      `struc:"uint8"`
	ExtendedSelftestRoutineRecommandedPollingTimeB  [2]uint8   `struc:"[2]uint8"`
	Reserved01                                      [9]uint8   `struc:"[9]uint8"`
	VendorSpecific386                               [125]uint8 `struc:"[125]uint8"`
	Checksum                                        uint8      `struc:"uint8"`
}

/**
 * ATA8-ACS
 * Table 49 — Device SMART data structure
 *
 * SFF-8035R2 Table 3
 */
type SmartAttributeThresholds struct {
	RevNumber         uint16 `struc:"uint16"`
	Attributes        [SMART_ATTRIBUTES_NUMBER]SmartAttribute
	Reserved01        [18]uint8  `struc:"[18]uint8"`
	VendorSpecific380 [131]uint8 `struc:"[131]uint8"`
	Checksum          uint8      `struc:"uint8"`
}

/*
 * Sanitize Device FEATURE field values
 */
const (
	SANITIZE_STATUS_EXT          uint16 = 0x0000
	SANITIZE_CRYPTO_SCRAMBLE_EXT uint16 = 0x0011
	SANITIZE_BLOCK_ERASE_EXT     uint16 = 0x0012
	SANITIZE_OVERWRITE_EXT       uint16 = 0x0014
	SANITIZE_FREEZE_LOCK_EXT     uint16 = 0x0020
	SANITIZE_ANTIFREEZE_LOCK_EXT uint16 = 0x0040
)

/*
 * Sanitize commands keys
 */
const (
	SANITIZE_FREEZE_LOCK_KEY     uint32 = 0x46724C6B /* "FrLk" */
	SANITIZE_ANTIFREEZE_LOCK_KEY uint32 = 0x416E7469 /* "Anti" */
	SANITIZE_CRYPTO_SCRAMBLE_KEY uint32 = 0x43727970 /* "Cryp" */
	SANITIZE_BLOCK_ERASE_KEY     uint32 = 0x426B4572 /* "BkEr" */
	SANITIZE_OVERWRITE_KEY       uint32 = 0x00004F57 /* "OW"   */
)

/*
 * Sanitize outputs flags
 */
const (
	SANITIZE_FLAG_OPERATION_SUCCEEDED   = (1 << 7)
	SANITIZE_FLAG_OPERATION_IN_PROGRESS = (1 << 6)
	SANITIZE_FLAG_DEVICE_IN_FROZEN      = (1 << 5)
	SANITIZE_FLAG_ANTIFREEZE_BIT        = (1 << 4)
)
