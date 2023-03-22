package go_dparm

import (
	"errors"
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/nvme"
	"github.com/jc-lab/go-dparm/tcg"
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

	// ATA
	AtaDoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error

	// NVME
	NvmeGetLogPage(nsid uint32, logId uint32, rae bool, size int) ([]byte, error)

	// COMMON
	SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error

	TcgDiscovery0() error
}

type DriveHandleImpl struct {
	DriveHandle
	dh   common.DriverHandle
	Info common.DriveInfo
}

var (
	ErrNotSupportThisDriver = errors.New("not supported in this driver")
)

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

	p.Info.TcgRawFeatures = make(map[uint16][]byte)
	_ = p.TcgDiscovery0()

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

func (p *DriveHandleImpl) AtaDoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	impl, ok := p.dh.(common.AtaDriverHandle)
	if ok {
		return impl.DoTaskFileCmd(rw, dma, tf, data, timeoutSecs)
	}
	return ErrNotSupportThisDriver
}

func (p *DriveHandleImpl) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, size int) ([]byte, error) {
	impl, ok := p.dh.(common.NvmeDriverHandle)
	if ok {
		return impl.NvmeGetLogPage(nsid, logId, rae, size)
	}
	return nil, ErrNotSupportThisDriver
}

func (p *DriveHandleImpl) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return p.dh.SecurityCommand(rw, dma, protocol, comId, buffer, timeoutSecs)
}

func (p *DriveHandleImpl) TcgDiscovery0() error {
	alignedBuffer := internal.NewAlignedBuffer(tcg.IO_BUFFER_ALIGNMENT, tcg.MIN_BUFFER_LENGTH)

	if err := p.SecurityCommand(false, false, 0x01, 0x0001, alignedBuffer.GetBuffer(), 3); err != nil {
		if err.Error() == "Not supported" {
			p.Info.TcgSupport = -1
		} else {
			p.Info.TcgSupport = 0
		}
		return err
	}

	p.Info.TcgSupport = 1

	alignedBuffer.ResetRead()
	header := &tcg.Discovery0Header{}
	headerSize, err := struc.Sizeof(header)
	if err != nil {
		return err
	}
	if err := struc.Unpack(alignedBuffer, header); err != nil {
		return err
	}
	bufferRef := alignedBuffer.GetBuffer()

	if len(bufferRef) < int(header.Length) {
		return errors.New("invalid data: length overflow")
	}

	offset := headerSize
	currentUnion := tcg.Discovery0FeatureUnion{}
	for offset < int(header.Length) {
		copy(currentUnion.Buffer[:], bufferRef[offset:])
		basic, err := currentUnion.ToBasic()
		if err != nil {
			return err
		}

		itemBuffer := bufferRef[offset : offset+int(basic.Length)+4]
		fc := basic.FeatureCode

		switch fc {
		case tcg.FcTPer:
			p.Info.TcgTper = true
		case tcg.FcLocking:
			p.Info.TcgLocking = true
		case tcg.FcGeometryReporting:
			p.Info.TcgGeometryReporting = true
		case tcg.FcOpalSscV100:
			p.Info.TcgOpalSscV100 = true
		case tcg.FcOpalSscV200:
			p.Info.TcgOpalSscV200 = true
		case tcg.FcEnterprise:
			p.Info.TcgEnterprise = true
		case tcg.FcSingleUser:
			p.Info.TcgSingleUser = true
		case tcg.FcDataStore:
			p.Info.TcgDataStore = true
		}

		p.Info.TcgRawFeatures[uint16(fc)] = itemBuffer
		offset += int(basic.Length) + 4
	}

	return nil
}
