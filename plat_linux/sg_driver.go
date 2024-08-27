//go:build linux
// +build linux

package plat_linux

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/scsi"
	"github.com/lunixbochs/struc"
)

const (
	SG_IO = 0x2285
)

type SgDriver struct {
	LinuxDriver
}

type SgDriverHandle struct {
	common.AtaDriverHandle
	d        *SgDriver
	fd       int
	identity [512]byte
}

func tfToLba(tf *ata.Tf) uint64 {
	var lba24, lbah uint32
	var lba64 uint64

	lba24 = (uint32(tf.Lob.Lbah) << 16) | (uint32(tf.Lob.Lbam) << 8) | (uint32(tf.Lob.Lbal))
	if tf.IsLba48 != 0 {
		lbah = ((uint32(tf.Hob.Lbah) << 16) | (uint32(tf.Hob.Lbam) << 8) | (uint32(tf.Hob.Lbal)))
	} else {
		lbah = uint32(tf.Dev & 0x0f)
	}
	lba64 = (uint64(lbah) << 24) | uint64(lba24)
	return lba64
}

func NewSgDriver() *SgDriver {
	return &SgDriver{}
}

func (d *SgDriver) OpenByFd(fd int) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(fd)
	if err != nil {
		return nil, err
	}
	return driverHandle, err
}

func (d *SgDriver) openImpl(fd int) (*SgDriverHandle, error) {
	tf := &ata.Tf{
		Command: ATA_IDENTIFY_DEVICE,
	}
	tf.Lob.Nsect = 1

	dataBuffer := internal.NewAlignedBuffer(512, 512)

	if err := d.doTaskFileCmd(fd, false, false, tf, dataBuffer.GetBuffer(), 3); err != nil {
		return nil, err
	}

	driverHandle := &SgDriverHandle{
		d:  d,
		fd: fd,
	}
	dataBuffer.ResetRead()
	dataBuffer.Read(driverHandle.identity[:])

	return driverHandle, nil
}

func (d *SgDriver) doTaskFileCmd(fd int, rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	strucOpts := internal.GetStrucOptions()
	var rootError error = nil

	if rw && data != nil {
		for i := range data {
			data[i] = 0
		}
	}

	sgParams := SG_IO_HDR_WITH_SENSE_BUF{}
	dataBuffer := make([]byte, 512)

	for retry := 0; retry < 2; retry++ {
		sgParams.InterfaceID = 'S'
		sgParams.Timeout = uint32(timeoutSecs)
		sgParams.DxferDirection = int32(internal.Ternary(data != nil, internal.Ternary(rw, SG_DXFER_TO_DEV, SG_DXFER_FROM_DEV), SG_DXFER_NONE))
		sgParams.Dxferp = uintptr(unsafe.Pointer(&dataBuffer))
		sgParams.MxSbLen = uint8(unsafe.Sizeof(sgParams.SenseData))
		sgParams.Sbp = &sgParams.SenseData[0]
		sgParams.PackID = int32(tfToLba(tf))

		if tf.IsLba48 != 0 {
			cdb := &ATA_PASSTHROUGH16{
				OperationCode:   SCSIOP_ATA_PASSTHROUGH16,
				Features7_0:     tf.Lob.Feat,
				Features15_8:    tf.Hob.Feat,
				SectorCount7_0:  tf.Lob.Nsect,
				SectorCount15_8: tf.Hob.Nsect,
				LbaLow7_0:       tf.Lob.Lbal,
				LbaLow15_8:      tf.Hob.Lbal,
				LbaMid7_0:       tf.Lob.Lbam,
				LbaMid15_8:      tf.Hob.Lbam,
				LbaHigh7_0:      tf.Lob.Lbah,
				LbaHigh15_8:     tf.Hob.Lbah,
				Device:          tf.Dev,
				Command:         uint8(tf.Command),
				Control:         0, // always zero
			}
			cdb.B01 |= SG_ATA_LBA48
			if dma {
				cdb.SetProtocol(uint8(internal.Ternary(data != nil, SG_ATA_PROTO_DMA, SG_ATA_PROTO_NON_DATA)))
			} else {
				cdb.SetProtocol(uint8(internal.Ternary(data != nil, (internal.Ternary(rw, SG_ATA_PROTO_PIO_OUT, SG_ATA_PROTO_PIO_IN)), SG_ATA_PROTO_NON_DATA)))
			}

			if data != nil {
				cdb.SetTLength(SG_CDB2_TLEN_NSECT)
				cdb.SetByteBlock(true)
				cdb.SetTDir(!rw)
			} else {
				cdb.SetCkCond(true)
			}

			n, err := struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			sgParams.CmdLen = byte(n)

			cmdBuffer := make([]byte, sgParams.CmdLen)
			sgParams.Cmdp = &cmdBuffer[0]

			writer := internal.NewWrappedBuffer(cmdBuffer[:])
			err = struc.PackWithOptions(writer, cdb, strucOpts)
			if err != nil {
				return err
			}
		} else {
			cdb := &ATA_PASSTHROUGH12{
				OperationCode: SCSIOP_ATA_PASSTHROUGH12,
				Features:      tf.Lob.Feat,
				SectorCount:   tf.Lob.Nsect,
				LbaLow:        tf.Lob.Lbal,
				LbaMid:        tf.Lob.Lbam,
				LbaHigh:       tf.Lob.Lbah,
				Device:        tf.Dev,
				Command:       uint8(tf.Command),
				Control:       0, // always zero
			}

			if dma {
				cdb.SetProtocol(uint8(internal.Ternary(data != nil, SG_ATA_PROTO_DMA, SG_ATA_PROTO_NON_DATA)))
			} else {
				cdb.SetProtocol(uint8(internal.Ternary(data != nil, (internal.Ternary(rw, SG_ATA_PROTO_PIO_OUT, SG_ATA_PROTO_PIO_IN)), SG_ATA_PROTO_NON_DATA)))
			}

			if data != nil {
				cdb.SetTLength(SG_CDB2_TLEN_NSECT)
				cdb.SetByteBlock(true)
				cdb.SetTDir(!rw)
			} else {
				cdb.SetCkCond(true)
			}

			n, err := struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			sgParams.CmdLen = byte(n)

			cmdBuffer := make([]byte, sgParams.CmdLen)
			sgParams.Cmdp = &cmdBuffer[0]

			writer := internal.NewWrappedBuffer(cmdBuffer[:])
			err = struc.PackWithOptions(writer, cdb, strucOpts)
			if err != nil {
				return err
			}
		}

		sgParams.DxferLen = 0
		if data != nil {
			sgParams.DxferLen = uint32(len(data))
		}

		if retry == 0 {
			var alignedBuffer *internal.AlignedBuffer = nil
			if data != nil {
				sgParams.Dxferp = uintptr(unsafe.Pointer(&data[0]))
				if !internal.IsAlignedPointer(512, sgParams.Dxferp) {
					alignedBuffer = internal.NewAlignedBuffer(512, len(data))
					if rw {
						alignedBuffer.ResetWrite()
						alignedBuffer.Write(data)
					}
					sgParams.Dxferp = uintptr(unsafe.Pointer(alignedBuffer.GetPointer()))
				}
			}

			_, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(fd),
				uintptr(SG_IO),
				uintptr(unsafe.Pointer(&sgParams)),
			)
			if err != unix.Errno(0) {
				rootError = err
			} else {
				if !rw && alignedBuffer != nil {
					alignedBuffer.ResetRead()
					alignedBuffer.Read(data)
				}
				break
			}
		} else {
			n := int(unsafe.Sizeof(sgParams))
			sgParams.Dxferp = uintptr(n)

			buffer := make([]byte, n+len(data))
			copyFromPointer(buffer, unsafe.Pointer(&sgParams), int(unsafe.Sizeof(sgParams)))
			if data != nil {
				copy(buffer[n:], data)
			}

			_, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(fd),
				uintptr(SG_IO),
				uintptr(unsafe.Pointer(&sgParams)),
			)
			if err != unix.Errno(0) {
				rootError = err
			} else {
				rootError = nil
				copyToPointer(unsafe.Pointer(&sgParams), buffer[:], int(unsafe.Sizeof(sgParams)))
				if !rw && data != nil {
					copy(data, buffer[n:])
				}
			}
		}
	}

	var senseInfo scsi.SENSE_DATA
	copyToPointer(unsafe.Pointer(&senseInfo), sgParams.SenseData[:], int(unsafe.Sizeof(senseInfo)))

	if rootError == nil {
		// ?? if sgParams.SenseInfo.IsValid() && sgParams.SenseInfo.GetSenseKey() != 0
		if senseInfo.GetSenseKey() != 0 {
			return &common.DparmError{
				DriverStatus: sgParams.Status,
				Message: fmt.Sprintf("SCSI status: %02x, Sense key: %#02x, ASC: %#02x, ASCQ: %#02x",
					sgParams.Status,
					senseInfo.GetSenseKey(), senseInfo.AdditionalSenseCode, senseInfo.AdditionalSenseCodeQualifier),
				SenseData: &senseInfo,
			}
		}
	}

	desc := sgParams.SenseData[8:]

	tf.IsLba48 = desc[2] & 1
	tf.Error = desc[3]
	tf.Lob.Nsect = desc[5]
	tf.Lob.Lbal = desc[7]
	tf.Lob.Lbam = desc[9]
	tf.Lob.Lbah = desc[11]
	tf.Dev = desc[12]
	tf.Status = desc[13]
	tf.Hob.Feat = 0
	if tf.IsLba48 != 0 {
		tf.Hob.Nsect = desc[4]
		tf.Hob.Lbal = desc[6]
		tf.Hob.Lbam = desc[8]
		tf.Hob.Lbah = desc[10]
	} else {
		tf.Hob.Nsect = 0
		tf.Hob.Lbal = 0
		tf.Hob.Lbam = 0
		tf.Hob.Lbah = 0
	}

	return rootError
}

func (s *SgDriverHandle) GetDriverName() string {
	return "LinuxScsiDriver"
}

func (s *SgDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingAtapi
}

func (s *SgDriverHandle) ReopenWritable() error {
	return nil
}

func (s *SgDriverHandle) Close() {
	_ = unix.Close(s.fd)
}

func (s *SgDriverHandle) GetIdentity() []byte {
	return s.identity[:]
}

func (s *SgDriverHandle) DoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	return s.d.doTaskFileCmd(s.fd, rw, dma, tf, data, timeoutSecs)
}

func (s *SgDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return scsiSecurityCommand(s.fd, rw, dma, protocol, comId, buffer, timeoutSecs)
}

func scsiSecurityCommand(fd int, rw bool, dma bool, protocol uint8, comId uint16, data []byte, timeoutSecs int) error {
	var rootError error = nil

	if rw && data != nil {
		for i := range data {
			data[i] = 0
		}
	}

	sgParams := SG_IO_HDR_WITH_SENSE_BUF{}

	for retry := 0; retry < 2; retry++ {
		sgParams.InterfaceID = 'S'
		sgParams.Timeout = uint32(timeoutSecs) * 1000
		sgParams.DxferDirection = int32(internal.Ternary(rw, SG_DXFER_FROM_DEV, SG_DXFER_TO_DEV))
		sgParams.SbLenWr = byte(unsafe.Sizeof(sgParams.SenseData))
		sgParams.Sbp = &sgParams.SenseData[0]

		cdb := &SCSI_SECURITY_PROTOCOL{}
		if rw {
			cdb.OperationCode = SCSIOP_SECURITY_PROTOCOL_OUT
		} else {
			cdb.OperationCode = SCSIOP_SECURITY_PROTOCOL_IN
		}
		cdb.Protocol = protocol
		cdb.ProtocolSp = comId
		cdb.Length = uint32(len(data))

		sizeOfCdb, err := struc.Sizeof(cdb)
		if err != nil {
			return err
		}

		sgParams.CmdLen = byte(sizeOfCdb)
		sgParams.DxferLen = uint32(len(data))

		cmdBuffer := make([]byte, sgParams.CmdLen)
		sgParams.Cmdp = &cmdBuffer[0]

		if err := struc.Pack(internal.NewWrappedBuffer([]byte{*sgParams.Cmdp}), cdb); err != nil {
			return err
		}

		if retry == 0 {
			var alignedBuffer *internal.AlignedBuffer = nil
			sgParams.Dxferp = uintptr(unsafe.Pointer(&data[0]))
			if !internal.IsAlignedPointer(512, sgParams.Dxferp) {
				alignedBuffer = internal.NewAlignedBuffer(512, len(data))
				if rw {
					alignedBuffer.ResetWrite()
					alignedBuffer.Write(data)
				}
				sgParams.Dxferp = uintptr(unsafe.Pointer(alignedBuffer.GetPointer()))
			}

			if _, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(fd),
				uintptr(SG_IO),
				uintptr(unsafe.Pointer(&sgParams)),
			); err != 0 {
				rootError = err
			} else {
				if !rw && alignedBuffer != nil {
					alignedBuffer.ResetRead()
					alignedBuffer.Read(data)
				}
				break
			}
		} else {
			n := int(unsafe.Sizeof(sgParams))
			sgParams.Dxferp = uintptr(n)

			buffer := make([]byte, n+len(data))
			copyFromPointer(buffer, unsafe.Pointer(&sgParams), int(unsafe.Sizeof(sgParams)))
			copy(buffer[n:], data)

			if _, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(fd),
				uintptr(SG_IO),
				uintptr(unsafe.Pointer(&sgParams)),
			); err != 0 {
				rootError = err
			} else {
				rootError = nil
				copyToPointer(unsafe.Pointer(&sgParams), buffer[:], int(unsafe.Sizeof(sgParams)))
				if !rw {
					copy(data, buffer[n:])
				}
			}
		}
	}

	var senseInfo scsi.SENSE_DATA
	copyToPointer(unsafe.Pointer(&senseInfo), sgParams.SenseData[:], int(unsafe.Sizeof(senseInfo)))

	if rootError == nil {
		// ?? if scsiParams.SenseInfo.IsValid() && scsiParams.SenseInfo.GetSenseKey() != 0
		if senseInfo.GetSenseKey() != 0 {
			return &common.DparmError{
				DriverStatus: sgParams.Status,
				Message: fmt.Sprintf("SCSI Status: %02x, Sense Key: %#02x, ASC: %#02x, ASCQ: %#02x",
					sgParams.Status,
					senseInfo.GetSenseKey(), senseInfo.AdditionalSenseCode, senseInfo.AdditionalSenseCodeQualifier),
				SenseData: &senseInfo,
			}
		}
	}

	return rootError
}
