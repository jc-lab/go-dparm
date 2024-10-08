package tcg

import (
	"fmt"
	"time"
	"unsafe"
)

type TcgDeviceType uint32

const (
	UnknownDeviceType TcgDeviceType = 0 + iota
	GenericDevice
	OpalV1Device
	OpalV2Device
	OpalEnterpriseDevice
)

type TcgDevice interface {
	GetSerial() string

	GetDeviceType() TcgDeviceType
	IsAnySSC() bool

	IsLockingSupported() bool
	IsLockingEnabled() bool
	IsLocked() bool
	IsMBREnabled() bool
	IsMBRDone() bool
	IsMediaEncryption() bool

	GetBaseComId() uint16
	GetNumComIds() uint16

	Exec(cmd *TcgCommand, protocol uint8) (*TcgResponse, error)

	GetDefaultPassword() (string, error)

	OpalGetTable(session *TcgSession, table []uint8, startCol, endCol uint16) (*TcgResponse, error)
	RevertTPer(password string, isPsid, isAdminSp bool) error
}

type TcgDeviceImpl struct {
	dev TcgDevice
	dh  *TcgDriveHandle
}

func NewTcgDevice(dch DriveCommandHandler) (TcgDevice, error) {
	tcgDriveHandle := NewTcgDriveHandle(dch)

	base := TcgDeviceImpl{
		dh: tcgDriveHandle,
	}
	base.dev = &base

	switch {
	case tcgDriveHandle.TcgOpalSscV100:
		device := &TcgDeviceOpal1{
			base,
		}
		device.dev = device

		return device, nil
	case tcgDriveHandle.TcgOpalSscV200:
		device := &TcgDeviceOpal2{
			base,
		}
		device.dev = device

		return device, nil
	case tcgDriveHandle.TcgEnterprise:
		device := &TcgDeviceEnterprise{
			base,
		}
		device.dev = device

		return device, nil
	}

	return &base, nil
}

func (p *TcgDeviceImpl) GetSerial() string {
	return p.dh.serial
}

func (p *TcgDeviceImpl) GetDeviceType() TcgDeviceType {
	return GenericDevice
}

func (p *TcgDeviceImpl) IsAnySSC() bool {
	return false
}

func (p *TcgDeviceImpl) IsLockingSupported() bool {
	return p.dh.TcgLocking
}

func (p *TcgDeviceImpl) IsLockingEnabled() bool {
	tcgDh := p.dh
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04>>1)&0x01 != 0
}

func (p *TcgDeviceImpl) IsLocked() bool {
	tcgDh := p.dh
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04>>2)&0x01 != 0
}

func (p *TcgDeviceImpl) IsMBREnabled() bool {
	tcgDh := p.dh
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04>>4)&0x01 != 0
}

func (p *TcgDeviceImpl) IsMBRDone() bool {
	tcgDh := p.dh
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04>>5)&0x01 != 0
}

func (p *TcgDeviceImpl) IsMediaEncryption() bool {
	tcgDh := p.dh
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04>>3)&0x01 != 0
}

func (p *TcgDeviceImpl) Exec(cmd *TcgCommand, protocol uint8) (*TcgResponse, error) {
	timeout := 10000 * time.Millisecond
	beginAt := time.Now()

	if !p.dev.IsAnySSC() {
		return nil, fmt.Errorf("not supported")
	}

	if err := p.dh.SecurityCommand(true, false, protocol, p.dev.GetBaseComId(), unsafe.Slice((*byte)(unsafe.Pointer(cmd.GetCmdPtr())), cmd.GetCmdSize()), 5); err != nil {
		return nil, err
	}

	resp := NewTcgResponse()
	for first := true; first || (resp.header.Cp.Outstanding != 0 && resp.header.Cp.MinTransfer == 0); {
		spentTime := time.Since(beginAt)
		time.Sleep(25 * time.Millisecond)
		resp.Reset()
		first = false

		if err := p.dh.SecurityCommand(false, false, protocol, p.dev.GetBaseComId(), unsafe.Slice((*byte)(unsafe.Pointer(resp.GetRespBuf())), resp.GetRespBufSize()), 5); err != nil {
			return resp, err
		}

		if spentTime > timeout {
			return resp, fmt.Errorf("timeout")
		}
	}

	if err := resp.Commit(); err != nil {
		return resp, err
	}

	return resp, nil
}

func (p *TcgDeviceImpl) GetBaseComId() uint16 {
	return 0
}

func (p *TcgDeviceImpl) GetNumComIds() uint16 {
	return 0
}

func (p *TcgDeviceImpl) GetDefaultPassword() (string, error) {
	return "", fmt.Errorf("not supported")
}

func (p *TcgDeviceImpl) OpalGetTable(session *TcgSession, table []uint8, startCol uint16, endCol uint16) (*TcgResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (p *TcgDeviceImpl) RevertTPer(password string, isPsid bool, isAdminSp bool) error {
	return fmt.Errorf("not supported")
}
