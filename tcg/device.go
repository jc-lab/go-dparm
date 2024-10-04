package tcg

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/jc-lab/go-dparm/common"
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
	GetDriveHandle() common.DriveHandle

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
	TcgDevice
	dh common.DriveHandle
}

func NewTcgDevice(driveHandle common.DriveHandle) (TcgDevice, error) {
	tcgDriveHandle := NewTcgDriveHandle(driveHandle)

	if err := tcgDriveHandle.TcgDiscovery0(); err != nil {
		return nil, err
	}

	return &TcgDeviceImpl{
		dh: tcgDriveHandle,
	}, nil
}

func (p *TcgDeviceImpl) GetDriveHandle() common.DriveHandle {
	return p.dh
}

func (p *TcgDeviceImpl) GetDeviceType() TcgDeviceType {
	return GenericDevice
}

func (p *TcgDeviceImpl) IsAnySSC() bool {
	return false
}

func (p *TcgDeviceImpl) IsLockingSupported() bool {
	return p.dh.(*TcgDriveHandle).TcgLocking
}

func (p *TcgDeviceImpl) IsLockingEnabled() bool {
	tcgDh := p.dh.(*TcgDriveHandle)
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04 >> 1) & 0x01 != 0
}

func (p *TcgDeviceImpl) IsLocked() bool {
	tcgDh := p.dh.(*TcgDriveHandle)
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04 >> 2) & 0x01 != 0
}

func (p *TcgDeviceImpl) IsMBREnabled() bool {
	tcgDh := p.dh.(*TcgDriveHandle)
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04 >> 4) & 0x01 != 0
}

func (p *TcgDeviceImpl) IsMBRDone() bool {
	tcgDh := p.dh.(*TcgDriveHandle)
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04 >> 5) & 0x01 != 0
}

func (p *TcgDeviceImpl) IsMediaEncryption() bool {
	tcgDh := p.dh.(*TcgDriveHandle)
	if !tcgDh.TcgLocking {
		return false
	}
	data, ok := tcgDh.TcgRawFeatures[uint16(FcLocking)]
	if !ok {
		return false
	}

	feature := (*Discovery0LockingFeature)(unsafe.Pointer(&data[0]))

	return (feature.B04 >> 3) & 0x01 != 0
}

func (p *TcgDeviceImpl) Exec(cmd *TcgCommand, protocol uint8) (*TcgResponse, error) {
	timeout := 10000 * time.Millisecond
	beginAt := time.Now()

	if !p.IsAnySSC() {
		return nil, fmt.Errorf("not supported")
	}

	if err := p.dh.SecurityCommand(true, false, protocol, p.TcgDevice.GetBaseComId(), unsafe.Slice((*byte)(unsafe.Pointer(cmd.GetCmdPtr())), cmd.GetCmdSize()), 5); err != nil {
		return nil, err
	}

	resp := NewTcgResponse()
	for resp.header.Cp.Outstanding == 0 || resp.header.Cp.MinTransfer != 0 {
		spentTime := time.Since(beginAt)
		time.Sleep(25 * time.Millisecond)
		resp.Reset()

		if err := p.dh.SecurityCommand(false, false, protocol, p.TcgDevice.GetBaseComId(), unsafe.Slice((*byte)(unsafe.Pointer(resp.GetRespBuf())), resp.GetRespBufSize()), 5); err != nil {
			return resp, err
		}

		if spentTime > timeout {
			// TIMEOUT!
			break;
		}
	}

	if resp.header.Cp.Outstanding != 0 && resp.header.Cp.MinTransfer == 0 {
		return resp, fmt.Errorf("timeout")
	}

	if err := resp.Commit(); err != nil {
		return resp, err
	}

	return resp, nil
}
