package common

import (
	"errors"
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/nvme"
	"github.com/jc-lab/go-dparm/tcg"
	"github.com/lunixbochs/struc"
	"strings"
)

const trimSet = " \t\r\n\x00"

type DriveHandleImpl struct {
	DriveHandle
	Dh   DriverHandle
	Info DriveInfo
}

var (
	ErrNotSupportThisDriver = errors.New("not supported in this driver")
)

func (p *DriveHandleImpl) Init() error {
	ataDrive, ok := p.Dh.(AtaDriverHandle)
	if ok {
		identityRaw := ataDrive.GetIdentity()
		p.Info.AtaIdentityRaw = identityRaw

		identity := &ata.IdentityDeviceData{}
		strucOpts := internal.GetStrucOptions()
		if err := struc.UnpackWithOptions(internal.NewWrappedBuffer(identityRaw), identity, strucOpts); err != nil {
			return err
		}
		p.Info.AtaIdentity = identity

		internal.AtaSwapWordEndian(identity.ModelNumber[:])
		internal.AtaSwapWordEndian(identity.SerialNumber[:])
		internal.AtaSwapWordEndian(identity.FirmwareRevision[:])

		p.Info.Model = strings.Trim(string(identity.ModelNumber[:]), trimSet)
		p.Info.FirmwareRevision = strings.Trim(string(identity.FirmwareRevision[:]), trimSet)
		rawSerial := identity.SerialNumber[:]
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

	nvmeDrive, ok := p.Dh.(NvmeDriverHandle)
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

func (p *DriveHandleImpl) GetDriverHandle() DriverHandle {
	return p.Dh
}

func (p *DriveHandleImpl) Close() {
	if p.Dh != nil {
		p.Dh.Close()
		p.Dh = nil
	}
}

func (p *DriveHandleImpl) GetDevicePath() string {
	return p.Info.DevicePath
}

func (p *DriveHandleImpl) GetDrivingType() DrivingType {
	return p.Dh.GetDrivingType()
}

func (p *DriveHandleImpl) GetDriverName() string {
	return p.Dh.GetDriverName()
}

func (p *DriveHandleImpl) GetDriveInfo() *DriveInfo {
	return &p.Info
}

func (p *DriveHandleImpl) AtaDoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	impl, ok := p.Dh.(AtaDriverHandle)
	if ok {
		return impl.DoTaskFileCmd(rw, dma, tf, data, timeoutSecs)
	}
	return ErrNotSupportThisDriver
}

func (p *DriveHandleImpl) NvmeGetLogPage(nsid uint32, logId uint32, rae bool, size int) ([]byte, error) {
	impl, ok := p.Dh.(NvmeDriverHandle)
	if ok {
		return impl.NvmeGetLogPage(nsid, logId, rae, size)
	}
	return nil, ErrNotSupportThisDriver
}

func (p *DriveHandleImpl) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	if p.Dh == nil {
		return errors.New("Not supported")
	}

	err := p.Dh.SecurityCommand(rw, dma, protocol, comId, buffer, timeoutSecs)
	if err == nil {
		return nil
	}

	ataDriver, ok := p.Dh.(AtaDriverHandle)
	if ok {
		tf := &ata.Tf{}
		tf.Lob.Feat = protocol
		tf.Lob.Nsect = uint8(len(buffer) / 512)
		if rw {
			tf.Command = internal.Ternary(dma, ata.ATA_OP_TRUSTED_SEND_DMA, ata.ATA_OP_TRUSTED_SEND)
		} else {
			tf.Command = internal.Ternary(dma, ata.ATA_OP_TRUSTED_RECV_DMA, ata.ATA_OP_TRUSTED_RECV)
		}
		tf.Lob.Lbam = uint8(comId)
		tf.Lob.Lbah = uint8(comId >> uint8(8))

		return ataDriver.DoTaskFileCmd(rw, dma, tf, buffer, timeoutSecs)
	}

	nvmeDriver, ok := p.Dh.(NvmeDriverHandle)
	if ok {
		cmd := &nvme.NvmeAdminCmd{}
		cmd.Opcode = uint8(internal.Ternary(rw, nvme.NVME_ADMIN_OP_SECURITY_SEND, nvme.NVME_ADMIN_OP_SECURITY_RECV))
		cmd.DataBuffer = buffer
		cmd.DataLen = uint32(len(buffer))
		cmd.Cdw10 = ((uint32(protocol) & 0xff) << 24) | ((uint32(comId) & 0xffff) << 8)
		cmd.Cdw11 = uint32(len(buffer))

		return nvmeDriver.DoNvmeAdminPassthru(cmd)
	}

	return err
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
