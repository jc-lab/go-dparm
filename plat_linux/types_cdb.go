package plat_linux

type CDB10 struct {
	OperationCode     uint8
	B01               uint8
	LogicalBlockByte  [4]uint8
	Reserved2         uint8
	TransferBlocksMsb uint8
	TransferBlocksLsb uint8
	Control           uint8
}

func (c *CDB10) IsRelativeAddress() bool {
	return c.B01&0x01 != 0
}

func (c *CDB10) SetRelativeAddress(v bool) {
	if v {
		c.B01 |= 0x01
	} else {
		c.B01 &= ^byte(0x01)
	}
}

func (c *CDB10) IsForceUnitAccess() bool {
	return c.B01&0x08 != 0
}

func (c *CDB10) SetForceUnitAccess(v bool) {
	if v {
		c.B01 |= 0x08
	} else {
		c.B01 &= ^byte(0x08)
	}
}

func (c *CDB10) IsDisablePageOut() bool {
	return c.B01&0x10 != 0
}

func (c *CDB10) SetDisablePageOut(v bool) {
	if v {
		c.B01 |= 0x10
	} else {
		c.B01 &= ^byte(0x10)
	}
}

func (c *CDB10) GetLogicalUnitNumber() uint8 {
	return c.B01 >> 5
}

func (c *CDB10) SetLogicalUnitNumber(v uint8) {
	c.B01 = (c.B01 & 0x1F) | (v << 5)
}

type CDB16 struct {
	OperationCode  uint8
	B01            uint8
	LogicalBlock   [8]uint8
	TransferLength [4]uint8
	Reserved2      uint8
	Control        uint8
}

func (c *CDB16) IsForceUnitAccess() bool {
	return c.B01&0x08 != 0
}

func (c *CDB16) SetForceUnitAccess(v bool) {
	if v {
		c.B01 |= 0x08
	} else {
		c.B01 &= ^byte(0x08)
	}
}

func (c *CDB16) IsDisablePageOut() bool {
	return c.B01&0x10 != 0
}

func (c *CDB16) SetDisablePageOut(v bool) {
	if v {
		c.B01 |= 0x10
	} else {
		c.B01 &= ^byte(0x10)
	}
}

func (c *CDB16) GetProtection() uint8 {
	return c.B01 >> 5
}

func (c *CDB16) SetProtection(v uint8) {
	c.B01 = (c.B01 & 0x1F) | (v << 5)
}
