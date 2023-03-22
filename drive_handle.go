package go_dparm

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/nvme"
	"github.com/lunixbochs/struc"
	"strings"
)

const trimSet = " \t\r\n\x00"

type DriveHandle interface {
	GetDriverHandle() common.DriverHandle
	Close()
	GetDevicePath() string
	GetDrivingType() common.DrivingType
	GetDriverName() string

	GetDriveInfo() *common.DriveInfo
}

type DriveHandleImpl struct {
	DriveHandle
	dh   common.DriverHandle
	Info common.DriveInfo
}

func (p *DriveHandleImpl) init() error {
	ataDrive, ok := p.dh.(common.AtaDriverHandle)
	if ok {
		identityRaw := ataDrive.GetIdentity()
		p.Info.AtaIdentityRaw = identityRaw

		identity := &ata.IdentityDeviceData{}
		if err := struc.Unpack(internal.NewWrappedBuffer(identityRaw), identity); err != nil {
			return err
		}
		p.Info.AtaIdentity = identity

		p.Info.Model = strings.Trim(string(ata.FixAtaStringOrder(identity.ModelNumber[:], true)), trimSet)
		p.Info.FirmwareRevision = strings.Trim(string(ata.FixAtaStringOrder(identity.FirmwareRevision[:], true)), trimSet)
		rawSerial := ata.FixAtaStringOrder(identity.SerialNumber[:], false)
		copy(p.Info.RawSerial[:], rawSerial)
		p.Info.Serial = strings.Trim(string(rawSerial), trimSet)
		p.Info.SmartEnabled = identity.CommandSetSupport.GetSmartCommands() && identity.CommandSetActive.GetSmartCommands()
		p.Info.SsdCheckWeight = 0
		if identity.NominalMediaRotationRate == 0 || identity.NominalMediaRotationRate == 1 {
			p.Info.SsdCheckWeight++
		}
		if identity.DataSetManagementFeature.GetTrim() {
			p.Info.SsdCheckWeight++
		}
		p.Info.IsSsd = p.Info.SsdCheckWeight > 0
	}

	nvmeDrive, ok := p.dh.(common.NvmeDriverHandle)
	if ok {
		identityRaw := nvmeDrive.GetIdentity()
		p.Info.NvmeIdentityRaw = identityRaw

		identity := &nvme.IdentifyController{}
		if err := struc.Unpack(internal.NewWrappedBuffer(identityRaw), identity); err != nil {
			return err
		}
		p.Info.NvmeIdentity = identity

		p.Info.Model = strings.Trim(string(identity.Mn[:]), trimSet)
		p.Info.FirmwareRevision = strings.Trim(string(identity.Fr[:]), trimSet)
		copy(p.Info.RawSerial[:], identity.Sn[:])
		p.Info.Serial = strings.Trim(string(identity.Sn[:]), trimSet)
		p.Info.SmartEnabled = true
		p.Info.IsSsd = true
		p.Info.SsdCheckWeight = 0
	}

	return nil
}

func (p *DriveHandleImpl) GetDriverHandle() common.DriverHandle {
	return p.dh
}

func (p *DriveHandleImpl) Close() {
	p.dh.Close()
}

func (p *DriveHandleImpl) GetDevicePath() string {
	return p.Info.DevicePath
}

func (p *DriveHandleImpl) GetDrivingType() common.DrivingType {
	return p.dh.GetDrivingType()
}

func (p *DriveHandleImpl) GetDriverName() string {
	return p.dh.GetDriverName()
}

func (p *DriveHandleImpl) GetDriveInfo() *common.DriveInfo {
	return &p.Info
}
