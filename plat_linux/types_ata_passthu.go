package plat_linux

const (
	ATA_PASSTHROUGH     = 0xa1
	ATA_TRUSTED_RCV     = 0x5c
	ATA_TRUSTED_SND     = 0x5e
	ATA_IDENTIFY_DEVICE = 0xec

	SCSI_INQUIRY          = 0x12
	SCSI_MODE_SENSE_6     = 0x1a
	SCSI_READ_CAPACITY_10 = 0x25
	SCSI_ATA_PASSTHRU_16  = 0x85
	SCSI_SECURITY_IN      = 0xa2
	SCSI_SECURITY_OUT     = 0xb5
)

const (
	HARD_RESET                  = 0
	SRST                        = 1
	NON_DATA                    = 3
	PIO_DATA_IN                 = 4
	PIO_DATA_OUT                = 5
	DMA                         = 6
	DMA_QUEUED                  = 7
	DEVICE_DIAGNOSTIC           = 8
	DEVICE_RESET                = 9
	UDMA_DATA_IN                = 10
	UDMA_DATA_OUT               = 11
	FPDMA                       = 12
	RETURN_RESPONSE_INFORMATION = 15
)

type ATA_PASSTHROUGH12 struct {
	OperationCode uint8
	B01           uint8
	B02           uint8
	Features      uint8
	SectorCount   uint8
	LbaLow        uint8
	LbaMid        uint8
	LbaHigh       uint8
	Device        uint8
	Command       uint8
	Reserved3     uint8
	Control       uint8
}

func (c *ATA_PASSTHROUGH12) GetProtocol() uint8 {
	return (c.B01 >> 1) & 0x0F
}

func (c *ATA_PASSTHROUGH12) SetProtocol(v uint8) {
	c.B01 = (c.B01 & 0xE1) | ((v & 0x0F) << 1)
}

func (c *ATA_PASSTHROUGH12) GetMultipleCount() uint8 {
	return (c.B01 >> 5) & 0x07
}

func (c *ATA_PASSTHROUGH12) SetMultipleCount(v uint8) {
	c.B02 = (c.B01 & 0x1F) | (v & 0x07)
}

func (c *ATA_PASSTHROUGH12) GetTLength() uint8 {
	return c.B02 & 0x03
}

func (c *ATA_PASSTHROUGH12) SetTLength(v uint8) {
	c.B02 = (c.B02 & 0xFC) | (v & 0x03)
}

func (c *ATA_PASSTHROUGH12) IsByteBlock() bool {
	return c.B02&0x04 != 0
}

func (c *ATA_PASSTHROUGH12) SetByteBlock(v bool) {
	if v {
		c.B02 |= 0x04
	} else {
		c.B02 &= ^byte(0x04)
	}
}

func (c *ATA_PASSTHROUGH12) IsTDir() bool {
	return c.B02&0x08 != 0
}

func (c *ATA_PASSTHROUGH12) SetTDir(v bool) {
	if v {
		c.B02 |= 0x08
	} else {
		c.B02 &= ^byte(0x08)
	}
}

func (c *ATA_PASSTHROUGH12) IsCkCond() bool {
	return c.B02&0x20 != 0
}

func (c *ATA_PASSTHROUGH12) SetCkCond(v bool) {
	if v {
		c.B02 |= 0x20
	} else {
		c.B02 &= ^byte(0x20)
	}
}

func (c *ATA_PASSTHROUGH12) GetOffline() uint8 {
	return (c.B02 >> 6) & 0x03
}

func (c *ATA_PASSTHROUGH12) SetOffline(v uint8) {
	c.B02 = (c.B02 & 0x3F) | ((v & 0x03) << 6)
}

type ATA_PASSTHROUGH16 struct {
	OperationCode   uint8 // 0x85 - SCSIOP_ATA_PASSTHROUGH16
	B01             uint8
	B02             uint8
	Features15_8    uint8
	Features7_0     uint8
	SectorCount15_8 uint8
	SectorCount7_0  uint8
	LbaLow15_8      uint8
	LbaLow7_0       uint8
	LbaMid15_8      uint8
	LbaMid7_0       uint8
	LbaHigh15_8     uint8
	LbaHigh7_0      uint8
	Device          uint8
	Command         uint8
	Control         uint8
}

func (c *ATA_PASSTHROUGH16) IsExtend() bool {
	return c.B02&0x01 != 0
}

func (c *ATA_PASSTHROUGH16) SetExtend(v bool) {
	if v {
		c.B02 |= 0x01
	} else {
		c.B02 &= ^byte(0x01)
	}
}

func (c *ATA_PASSTHROUGH16) GetProtocol() uint8 {
	return (c.B01 >> 1) & 0x0F
}

func (c *ATA_PASSTHROUGH16) SetProtocol(v uint8) {
	c.B01 = (c.B01 & 0xE1) | ((v & 0x0F) << 1)
}

func (c *ATA_PASSTHROUGH16) GetMultipleCount() uint8 {
	return (c.B01 >> 5) & 0x07
}

func (c *ATA_PASSTHROUGH16) SetMultipleCount(v uint8) {
	c.B02 = (c.B02 & 0x1F) | (v & 0x07)
}

func (c *ATA_PASSTHROUGH16) GetTLength() uint8 {
	return c.B02 & 0x03
}

func (c *ATA_PASSTHROUGH16) SetTLength(v uint8) {
	c.B02 = (c.B02 & 0xFC) | (v & 0x03)
}

func (c *ATA_PASSTHROUGH16) IsByteBlock() bool {
	return c.B02&0x04 != 0
}

func (c *ATA_PASSTHROUGH16) SetByteBlock(v bool) {
	if v {
		c.B02 |= 0x04
	} else {
		c.B02 &= ^byte(0x04)
	}
}

func (c *ATA_PASSTHROUGH16) IsTDir() bool {
	return c.B02&0x08 != 0
}

func (c *ATA_PASSTHROUGH16) SetTDir(v bool) {
	if v {
		c.B02 |= 0x08
	} else {
		c.B02 &= ^byte(0x08)
	}
}

func (c *ATA_PASSTHROUGH16) IsCkCond() bool {
	return c.B02&0x20 != 0
}

func (c *ATA_PASSTHROUGH16) SetCkCond(v bool) {
	if v {
		c.B02 |= 0x20
	} else {
		c.B02 &= ^byte(0x20)
	}
}

func (c *ATA_PASSTHROUGH16) GetOffline() uint8 {
	return (c.B02 >> 6) & 0x03
}

func (c *ATA_PASSTHROUGH16) SetOffline(v uint8) {
	c.B02 = (c.B02 & 0x3F) | ((v & 0x03) << 6)
}
