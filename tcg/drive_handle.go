package tcg

import (
	"fmt"

	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/lunixbochs/struc"
)

type TcgDriveHandle struct {
	common.DriveHandle
	TcgSupport           int
	TcgTper              bool
	TcgLocking           bool
	TcgGeometryReporting bool
	TcgOpalSscV100       bool
	TcgOpalSscV200       bool
	TcgEnterprise        bool
	TcgSingleUser        bool
	TcgDataStore         bool

	TcgRawFeatures map[uint16][]byte
}

func NewTcgDriveHandle(dh common.DriveHandle) *TcgDriveHandle {
	return &TcgDriveHandle{
		DriveHandle:    dh,
		TcgRawFeatures: make(map[uint16][]byte),
	}
}

func (p *TcgDriveHandle) TcgDiscovery0() error {
	alignedBuffer := internal.NewAlignedBuffer(IO_BUFFER_ALIGNMENT, MIN_BUFFER_LENGTH)

	if err := p.SecurityCommand(false, false, 0x01, 0x0001, alignedBuffer.GetBuffer(), 3); err != nil {
		if err.Error() == "not supported" {
			p.TcgSupport = -1
		} else {
			p.TcgSupport = 0
		}
		return err
	}

	p.TcgSupport = 1

	alignedBuffer.ResetRead()
	header := &Discovery0Header{}
	headerSize, err := struc.Sizeof(header)
	if err != nil {
		return err
	}
	if err := struc.Unpack(alignedBuffer, header); err != nil {
		return err
	}
	bufferRef := alignedBuffer.GetBuffer()

	if len(bufferRef) < int(header.Length) {
		return fmt.Errorf("invalid data: length overflow")
	}

	offset := headerSize
	currentUnion := Discovery0FeatureUnion{}
	for offset < int(header.Length) {
		copy(currentUnion.Buffer[:], bufferRef[offset:])
		basic, err := currentUnion.ToBasic()
		if err != nil {
			return err
		}

		itemBuffer := bufferRef[offset : offset+int(basic.Length)+4]
		fc := basic.FeatureCode

		switch fc {
		case FcTPer:
			p.TcgTper = true
		case FcLocking:
			p.TcgLocking = true
		case FcGeometryReporting:
			p.TcgGeometryReporting = true
		case FcOpalSscV100:
			p.TcgOpalSscV100 = true
		case FcOpalSscV200:
			p.TcgOpalSscV200 = true
		case FcEnterprise:
			p.TcgEnterprise = true
		case FcSingleUser:
			p.TcgSingleUser = true
		case FcDataStore:
			p.TcgDataStore = true
		}

		p.TcgRawFeatures[uint16(fc)] = itemBuffer
		offset += int(basic.Length) + 4
	}

	return nil
}
