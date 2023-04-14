package plat_linux

const (
	MAX_PARTITIONS = 64
)

// device types

const (
	NO_DEV = 0xffff
	ATA_DEV = 0x0000
	ATAPI_DEV = 0x0001
)


// bit definitions within the words

const (
	VALID = 0xc000
	VALID_VAL = 0x4000
)

// word 0: gen config

const (
    NOT_ATA           = 0x8000
    NOT_ATAPI         = 0x4000 // (check only if bit 15 == 1)
    MEDIA_REMOVABLE   = 0x0080
    DRIVE_NOT_REMOVABLE = 0x0040 // bit obsoleted in ATA 6
    INCOMPLETE        = 0x0004
    DRQ_RESPONSE_TIME = 0x0060
    DRQ_3MS_VAL       = 0x0000
    DRQ_INTR_VAL      = 0x0020
    DRQ_50US_VAL      = 0x0040
    PKT_SIZE_SUPPORTED = 0x0003
    PKT_SIZE_12_VAL   = 0x0000
    PKT_SIZE_16_VAL   = 0x0001
    EQPT_TYPE         = 0x1f00
    SHIFT_EQPT        = 8
)

const (
	CDROM = 0x0005
)
