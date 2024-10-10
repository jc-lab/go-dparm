package tcg

import (
	"github.com/jc-lab/go-dparm/internal"
	"github.com/lunixbochs/struc"
)

const (
	IO_BUFFER_ALIGNMENT = 1024
	MAX_BUFFER_LENGTH   = 61440
	MIN_BUFFER_LENGTH   = 2048
)

type FeatureCode uint16

const (
	FcTPer              FeatureCode = 0x0001
	FcLocking           FeatureCode = 0x0002
	FcGeometryReporting FeatureCode = 0x0003
	FcEnterprise        FeatureCode = 0x0100
	FcDataStore         FeatureCode = 0x0202
	FcSingleUser        FeatureCode = 0x0201
	FcOpalSscV100       FeatureCode = 0x0200
	FcOpalSscV200       FeatureCode = 0x0203
)

type VersionField struct {
	B02 uint8 `struc:"uint8,big"`
}

func (f *VersionField) GetVersion() uint8 {
	return (f.B02 >> 4) & 0x0f
}

func (f *VersionField) SetVersion(version uint8) {
	f.B02 = (f.B02 & 0x0f) | ((version & 0x0F) << 4)
}

// The Discovery 0 Header
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_SSC_v2.01_rev1.00.pdf
// 3.1.1.1
type Discovery0Header struct {
	Length     uint32   `struc:"uint32,big"` // the length of the header 48 in 2.00.100
	Revision   uint32   `struc:"uint32,big"` // revision of the header 1 in 2.00.100
	Reserved01 uint32   `struc:"uint32,big"`
	Reserved02 uint32   `struc:"uint32,big"`
	Reserved03 [16]byte `struc:"[16]byte"`
}

type Discovery0BasicFeature struct {
	FeatureCode FeatureCode `struc:"uint16,big"`
	VersionField
	//ReservedV   uint8  `struc:"uint8,big"`
	//Version     uint8  `struc:"uint8,big"`
	Length uint8 `struc:"uint8,big"`
}

// Discovery0TPerFeature The Discovery 0 - TPer Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_SSC_v2.01_rev1.00.pdf
// 3.1.1.2
type Discovery0TPerFeature struct {
	FeatureCode uint16 `struc:"uint16,big"`
	VersionField
	Length uint8 `struc:"uint8,big"`
	//Reserved01       uint8  `struc:"uint8,big"`
	//ComIDManagement  uint8  `struc:"uint8,big"`
	//Reserved02       uint8  `struc:"uint8,big"`
	//Streaming        uint8  `struc:"uint8,big"`
	//BufferManagement uint8  `struc:"uint8,big"`
	//Acknack          uint8  `struc:"uint8,big"`
	//Async            uint8  `struc:"uint8,big"`
	//Sync             uint8  `struc:"uint8,big"`
	B04        uint8  `struc:"uint8,big"`
	Reserved03 uint32 `struc:"uint32,big"`
	Reserved04 uint32 `struc:"uint32,big"`
	Reserved05 uint32 `struc:"uint32,big"`
}

// Discovery0LockingFeature The Discovery 0 - Locking Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_SSC_v2.01_rev1.00.pdf
// 3.1.1.3
type Discovery0LockingFeature struct {
	FeatureCode uint16 `struc:"uint16,big"`
	VersionField
	Length uint8 `struc:"uint8,big"`
	//Reserved01       uint8  `struc:"uint8,big"`
	//Reserved02       uint8  `struc:"uint8,big"`
	//MBRDone          uint8  `struc:"uint8,big"`
	//MBREnabled       uint8  `struc:"uint8,big"`
	//MediaEncryption  uint8  `struc:"uint8,big"`
	//Locked           uint8  `struc:"uint8,big"`
	//LockingEnabled   uint8  `struc:"uint8,big"`
	//LockingSupported uint8  `struc:"uint8,big"`
	B04        uint8  `struc:"uint8,big"`
	Reserved03 uint32 `struc:"uint32,big"`
	Reserved04 uint32 `struc:"uint32,big"`
	Reserved05 uint32 `struc:"uint32,big"`
}

// Discovery0GeometryReportingFeature The Discovery 0 - Geometry Reporting Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_SSC_v2.01_rev1.00.pdf
// 3.1.1.4
type Discovery0GeometryReportingFeature struct {
	FeatureCode uint16 `struc:"uint16,big"`
	VersionField
	Length uint8 `struc:"uint8,big"`
	//uint8_t align: 1
	//uint8_t reserved01: 7
	B04                  uint8    `struc:"uint8,big"`
	Reserved02           [7]uint8 `struc:"[7]uint8,big"`
	LogicalBlockSize     uint32   `struc:"uint32,big"`
	AlignmentGranularity uint64   `struc:"uint64,big"`
	LowestAlignedLba     uint64   `struc:"uint64,big"`
}

// Discovery0OpalSSCFeatureV100 is the Discovery 0 - Opal SSC Feature
// https://trustedcomputinggroup.org/wp-content/uploads/Opal_SSC_1.00_rev3.00-Final.pdf
// 3.1.1.4
type Discovery0OpalSSCFeatureV100 struct {
	FeatureCode uint16 `struc:"uint16,big"` // 0x0200
	VersionField
	Length    uint8  `struc:"uint8"`
	BaseComID uint16 `struc:"uint16,big"`
	NumComIDs uint16 `struc:"uint16,big"`
}

// Discovery0OpalSSCFeatureV200 is the Discovery 0 -
// Opal SSC V2.00 Feature
// https://trustedcomputinggroup.org/wp-content/uploads/Opal_SSC_1.00_rev3.00-Final.pdf
// 3.1.1.4
type Discovery0OpalSSCFeatureV200 struct {
	FeatureCode uint16 `struc:"uint16,big"` // 0x0203
	VersionField
	Length          uint8  `struc:"uint8"`
	BaseComID       uint16 `struc:"uint16,big"`
	NumComIDs       uint16 `struc:"uint16,big"`
	B08             uint8  `struc:"uint8"`
	NumLockingAdmin uint16 `struc:"uint16,big"`
	NumLockingUser  uint16 `struc:"uint16,big"`
	InitialPin      uint8  `struc:"uint8"`
	RevetedPin      uint8  `struc:"uint8"`
	Reserved02      uint8  `struc:"uint8"`
	Reserved03      uint32 `struc:"uint32,big"`
}

// Discovery0EnterpriseSSCFeature is the Discovery 0 - Enterprise SSC Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-SSC_Enterprise-v1.01_r1.00.pdf
// 3.6.2.7
type Discovery0EnterpriseSSCFeature struct {
	FeatureCode uint16 `struc:"uint16,big"` // 0x0100
	VersionField
	Length       uint8  `struc:"uint8"`
	BaseComID    uint16 `struc:"uint16,big"`
	NumberComIDs uint16 `struc:"uint16,big"`
	B08          uint8  `struc:"uint8"`
	Reserved02   uint8  `struc:"uint8"`
	Reserved03   uint16 `struc:"uint16,big"`
	Reserved04   uint32 `struc:"uint32,big"`
	Reserved05   uint32 `struc:"uint32,big"`
}

// Discovery0SingleUserModeFeature is the Discovery 0 - Single User Mode Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_Feature_Set_Single_User_Mode_v1-00_r1-00-Final.pdf
// 4.2.1
type Discovery0SingleUserModeFeature struct {
	FeatureCode uint16 `struc:"uint16,big"` // 0x0201
	VersionField
	Length           uint8  `struc:"uint8"`
	NumberLockingObj uint32 `struc:"uint32,big"`
	B08              uint8  `struc:"uint8"`
	Reserved02       uint8  `struc:"uint8"`
	Reserved03       uint16 `struc:"uint16,big"`
	Reserved04       uint32 `struc:"uint32,big"`
}

// Discovery0DataStoreTableFeature is the Discovery 0 - DataStore Table Feature
// https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage-Opal_Feature_Set-Additional_DataStore_Tables_v1_00_r1_00_Final.pdf
// 4.1.1
type Discovery0DataStoreTableFeature struct {
	FeatureCode uint16 `struc:"uint16,big"` // 0x0203
	VersionField
	Length             uint8  `struc:"uint8"`
	Reserved01         uint16 `struc:"uint16,big"`
	MaxTables          uint16 `struc:"uint16,big"`
	MaxSizeTables      uint32 `struc:"uint32,big"`
	TableSizeAlignment uint32 `struc:"uint32,big"`
}

type Discovery0FeatureUnion struct {
	Buffer [16]byte
}

func (u *Discovery0FeatureUnion) ToBasic() (*Discovery0BasicFeature, error) {
	result := &Discovery0BasicFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToTper() (*Discovery0TPerFeature, error) {
	result := &Discovery0TPerFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToLocking() (*Discovery0LockingFeature, error) {
	result := &Discovery0LockingFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToGeometry() (*Discovery0GeometryReportingFeature, error) {
	result := &Discovery0GeometryReportingFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToOpalSscV100() (*Discovery0OpalSSCFeatureV100, error) {
	result := &Discovery0OpalSSCFeatureV100{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToOpalSscV200() (*Discovery0OpalSSCFeatureV200, error) {
	result := &Discovery0OpalSSCFeatureV200{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToEnterpriseSSC() (*Discovery0EnterpriseSSCFeature, error) {
	result := &Discovery0EnterpriseSSCFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToSingleUserMode() (*Discovery0SingleUserModeFeature, error) {
	result := &Discovery0SingleUserModeFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Discovery0FeatureUnion) ToDataStoreTable() (*Discovery0DataStoreTableFeature, error) {
	result := &Discovery0DataStoreTableFeature{}
	if err := struc.Unpack(internal.NewWrappedBuffer(u.Buffer[:]), result); err != nil {
		return nil, err
	}
	return result, nil
}

// OpalComPacket is Reference: https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage_Opal_SSC_Application_Note_1-00_1-00-Final.pdf
type OpalComPacket struct {
	Reserved0     uint32 `struc:"uint32,big"`
	ExtendedComID [4]byte
	Outstanding   uint32 `struc:"uint32,big"`
	MinTransfer   uint32 `struc:"uint32,big"`
	Length        uint32 `struc:"uint32,big"`
}

type OpalPacket struct {
	Tsn        uint32 `struc:"uint32,big"`
	Hsn        uint32 `struc:"uint32,big"`
	SeqNumber  uint32 `struc:"uint32,big"`
	Reserved00 uint16 `struc:"uint16,big"`
	AckType    uint16 `struc:"uint16,big"`
	Ack        uint32 `struc:"uint32,big"`
	Length     uint32 `struc:"uint32,big"`
}

type OpalDataSubPacket struct {
	Reserved00 [6]byte `struc:"[6]byte,big"`
	Kind       uint16  `struc:"uint16,big"`
	Length     uint32  `struc:"uint32,big"`
}

type OpalHeader struct {
	Cp     OpalComPacket
	Pkt    OpalPacket
	Subpkt OpalDataSubPacket
}

type Buf []byte

type invokingUID interface {
	invokingUid() // dummy
}

type signAuthority interface {
	signAuthority() // dummy
}

type cmdMethod interface {
	cmdMethod() //dummy
}

type OpalUID [8]byte
type OpalMethod [8]byte

func (Buf) invokingUid()        {}
func (OpalUID) invokingUid()    {}
func (OpalMethod) invokingUid() {}

func (Buf) signAuthority()     {}
func (OpalUID) signAuthority() {}

func (Buf) cmdMethod()        {}
func (OpalMethod) cmdMethod() {}

var (
	SMUID_UID                  OpalUID = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
	THISSP_UID                 OpalUID = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	ADMINSP_UID                OpalUID = [8]byte{0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x01}
	LOCKINGSP_UID              OpalUID = [8]byte{0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x02}
	ENTERPRISE_LOCKINGSP_UID   OpalUID = [8]byte{0x00, 0x00, 0x02, 0x05, 0x00, 0x01, 0x00, 0x01}
	ANYBODY_UID                OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01}
	SID_UID                    OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x06}
	ADMIN1_UID                 OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x01, 0x00, 0x01}
	USER1_UID                  OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x03, 0x00, 0x01}
	USER2_UID                  OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x03, 0x00, 0x02}
	PSID_UID                   OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x01, 0xff, 0x01}
	ENTERPRISE_BANDMASTER0_UID OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x80, 0x01}
	ENTERPRISE_ERASEMASTER_UID OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x84, 0x01}

	/* tables */
	LOCKINGRANGE_GLOBAL           OpalUID = [8]byte{0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x00, 0x01}
	LOCKINGRANGE_ACE_RDLOCKED     OpalUID = [8]byte{0x00, 0x00, 0x00, 0x08, 0x00, 0x03, 0xE0, 0x01}
	LOCKINGRANGE_ACE_WRLOCKED     OpalUID = [8]byte{0x00, 0x00, 0x00, 0x08, 0x00, 0x03, 0xE8, 0x01}
	MBRCONTROL                    OpalUID = [8]byte{0x00, 0x00, 0x08, 0x03, 0x00, 0x00, 0x00, 0x01}
	MBR                           OpalUID = [8]byte{0x00, 0x00, 0x08, 0x04, 0x00, 0x00, 0x00, 0x00}
	AUTHORITY_TABLE               OpalUID = [8]byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00}
	C_PIN_TABLE                   OpalUID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x00}
	LOCKING_INFO_TABLE            OpalUID = [8]byte{0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01}
	ENTERPRISE_LOCKING_INFO_TABLE OpalUID = [8]byte{0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x00}

	/* C_PIN_TABLE object ID's */
	C_PIN_MSID   OpalUID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x84, 0x02}
	C_PIN_SID    OpalUID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x01}
	C_PIN_ADMIN1 OpalUID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x01, 0x00, 0x01}

	/* half UID's (only first 4 bytes used) */
	HALF_UID_AUTHORITY_OBJ_REF OpalUID = [8]byte{0x00, 0x00, 0x0C, 0x05, 0xff, 0xff, 0xff, 0xff}
	HALF_UID_BOOLEAN_ACE       OpalUID = [8]byte{0x00, 0x00, 0x04, 0x0E, 0xff, 0xff, 0xff, 0xff}

	/* special value for omitted optional parameter */
	UID_HEXFF OpalUID = [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

var (
	PROPERTIES    OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x01}
	STARTSESSION  OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x02}
	REVERT        OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x02, 0x02}
	ACTIVATE      OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x02, 0x03}
	EGET          OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x06}
	ESET          OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x07}
	NEXT          OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x08}
	EAUTHENTICATE OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x0c}
	GETACL        OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x0d}
	GENKEY        OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x10}
	REVERTSP      OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x11}
	GET           OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x16}
	SET           OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x17}
	AUTHENTICATE  OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x1c}
	RANDOM        OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x06, 0x01}
	ERASE         OpalMethod = [8]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x08, 0x03}
)

/*
 * Reference: https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage_Opal_SSC_Application_Note_1-00_1-00-Final.pdf
 */
type OpalToken int

const (
	// Boolean
	TRUE         OpalToken = 0x01
	FALSE        OpalToken = 0x00
	BOOLEAN_EXPR OpalToken = 0x03

	/**
	 * Cell Blocks
	 */
	TABLE       OpalToken = 0x00
	STARTROW    OpalToken = 0x01
	ENDROW      OpalToken = 0x02
	STARTCOLUMN OpalToken = 0x03
	ENDCOLUMN   OpalToken = 0x04
	VALUES      OpalToken = 0x01

	/*
	 * Credential Table Group
	 *
	 * Reference: https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage_Architecture_Core_Spec_v2.01_r1.00.pdf
	 * Table 181. C_PIN Table Description
	 * */
	CREDENTIAL_UID         OpalToken = 0x00
	CREDENTIAL_NAME        OpalToken = 0x01
	CREDENTIAL_COMMON_NAME OpalToken = 0x02
	CREDENTIAL_PIN         OpalToken = 0x03
	CREDENTIAL_CHAR_SET    OpalToken = 0x04
	CREDENTIAL_TRY_LIMIT   OpalToken = 0x05
	CREDENTIAL_TRIES       OpalToken = 0x06
	CREDENTIAL_PERSISTENCE OpalToken = 0x07

	/*
	 * Locking Table
	 *
	 * Reference: https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage_Architecture_Core_Spec_v2.01_r1.00.pdf
	 * Table 226. Locking Table Description
	 * */
	LOCKING_UID                OpalToken = 0x00
	LOCKING_NAME               OpalToken = 0x01
	LOCKING_COMMON_NAME        OpalToken = 0x02
	LOCKING_RANGE_START        OpalToken = 0x03
	LOCKING_RANGE_LENGTH       OpalToken = 0x04
	LOCKING_READ_LOCK_ENABLED  OpalToken = 0x05
	LOCKING_WRITE_LOCK_ENABLED OpalToken = 0x06
	LOCKING_READ_LOCKED        OpalToken = 0x07
	LOCKING_WRITE_LOCKED       OpalToken = 0x08
	LOCKING_LOCK_ON_RESET      OpalToken = 0x09
	LOCKING_ACTIVE_KEY         OpalToken = 0x0A
	LOCKING_NEXT_KEY           OpalToken = 0x0B
	LOCKING_GENERAL_STATUS     OpalToken = 0x13

	/*
	 * LockingInfo Table
	 *
	 * Reference: https://trustedcomputinggroup.org/wp-content/uploads/TCG_Storage_Architecture_Core_Spec_v2.01_r1.00.pdf
	 * Table 225. LockingInfo Table Description
	 * */
	LOCKINGINFO_UID             OpalToken = 0x00
	LOCKINGINFO_NAME            OpalToken = 0x01
	LOCKINGINFO_COMMON_NAME     OpalToken = 0x02
	LOCKINGINFO_ENCRYPT_SUPPORT OpalToken = 0x03
	LOCKINGINFO_MAXRANGES       OpalToken = 0x04

	/* mbr control */
	MBRENABLE OpalToken = 0x01
	MBRDONE   OpalToken = 0x02

	/* properties */
	HOSTPROPERTIES OpalToken = 0x00

	/* response tokenis() returned values */
	DTA_TOKENID_BYTESTRING OpalToken = 0xe0
	DTA_TOKENID_SINT       OpalToken = 0xe1
	DTA_TOKENID_UINT       OpalToken = 0xe2
	DTA_TOKENID_TOKEN      OpalToken = 0xe3 // actual token is returned

	STARTLIST       OpalToken = 0xf0
	ENDLIST         OpalToken = 0xf1
	STARTNAME       OpalToken = 0xf2
	ENDNAME         OpalToken = 0xf3
	CALL            OpalToken = 0xf8
	ENDOFDATA       OpalToken = 0xf9
	ENDOFSESSION    OpalToken = 0xfa
	STARTTRANSACTON OpalToken = 0xfb
	ENDTRANSACTON   OpalToken = 0xfc
	EMPTYATOM       OpalToken = 0xff
	WHERE           OpalToken = 0x00
)

type OpalTinyAtom int

const (
	UINT_00 OpalTinyAtom = 0x00
	UINT_01 OpalTinyAtom = 0x01
	UINT_02 OpalTinyAtom = 0x02
	UINT_03 OpalTinyAtom = 0x03
	UINT_04 OpalTinyAtom = 0x04
	UINT_05 OpalTinyAtom = 0x05
	UINT_06 OpalTinyAtom = 0x06
	UINT_07 OpalTinyAtom = 0x07
	UINT_08 OpalTinyAtom = 0x08
	UINT_09 OpalTinyAtom = 0x09
	UINT_10 OpalTinyAtom = 0x0a
	UINT_11 OpalTinyAtom = 0x0b
	UINT_12 OpalTinyAtom = 0x0c
	UINT_13 OpalTinyAtom = 0x0d
	UINT_14 OpalTinyAtom = 0x0e
	UINT_15 OpalTinyAtom = 0x0f
)

type OpalShortAtom int

const (
	UINT_3      OpalShortAtom = 0x83
	BYTESTRING4 OpalShortAtom = 0xa4
	BYTESTRING8 OpalShortAtom = 0xa8
)

type OpalLockingState int

const (
	READWRITE       OpalLockingState = 0x01
	READONLY        OpalLockingState = 0x02
	LOCKED          OpalLockingState = 0x03
	ARCHIVELOCKED   OpalLockingState = 0x04
	ARCHIVEUNLOCKED OpalLockingState = 0x05
)

type MethodStatus int

const (
	SUCCESS               MethodStatus = 0x00
	NOT_AUTHORIZED        MethodStatus = 0x01
	SP_BUSY               MethodStatus = 0x03
	SP_FAILED             MethodStatus = 0x04
	SP_DISABLED           MethodStatus = 0x05
	SP_FROZEN             MethodStatus = 0x06
	NO_SESSIONS_AVAILABLE MethodStatus = 0x07
	UNIQUENESS_CONFLICT   MethodStatus = 0x08
	INSUFFICIENT_SPACE    MethodStatus = 0x09
	INSUFFICIENT_ROWS     MethodStatus = 0x0A
	INVALID_FUNCTION      MethodStatus = 0x0B // defined in early specs, still used in some firmware
	INVALID_PARAMETER     MethodStatus = 0x0C
	INVALID_REFERENCE     MethodStatus = 0x0D // OBSOLETE
	TPER_MALFUNCTION      MethodStatus = 0x0F
	TRANSACTION_FAILURE   MethodStatus = 0x10
	RESPONSE_OVERFLOW     MethodStatus = 0x11
	AUTHORITY_LOCKED_OUT  MethodStatus = 0x12
	FAIL                  MethodStatus = 0x3F
)

type token interface {
	token() // dummy
}

func (OpalToken) token()        {}
func (OpalTinyAtom) token()     {}
func (OpalShortAtom) token()    {}
func (OpalLockingState) token() {}
func (OpalUID) token()          {}
func (OpalMethod) token()       {}
