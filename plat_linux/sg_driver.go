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

// Should not be exported in release, exported for testing!
type SgDriverHandle struct {
	common.AtaDriverHandle
	D *SgDriver
	Fd int
	Identity [512]byte
}

func tfToLba(tf *ata.Tf) uint64 {
	var lba24, lbah uint32
	var lba64 uint64

	lba24 = uint32((tf.Lob.Lbah << 16) | (tf.Lob.Lbam << 8) | (tf.Lob.Lbal))
	if tf.IsLba48 != 0 {
		lbah = uint32((tf.Hob.Lbah << 16) | (tf.Hob.Lbam << 8) | (tf.Hob.Lbal))
	} else {
		lbah = uint32(tf.Dev & 0x0f)
	}
	lba64 = (uint64(lbah) << 24) | uint64(lba24)
	return lba64
}

func NewSgDriver() *SgDriver {
	return &SgDriver{}
}

func (d *SgDriver) OpenByPath(path string) (common.DriverHandle, error) {
	fd, err := OpenDevice(path)
	if err != nil {
		return nil, err
	}

	driverHandle, err := d.openImpl(fd)
	if err != nil {
		_ = unix.Close(fd)
	}
	return driverHandle, err
}

func (d *SgDriver) openImpl(fd int) (common.DriverHandle, error) {
	tf := &ata.Tf {
		Command: ATA_IDENTIFY_DEVICE,
	}
	tf.Lob.Nsect = 1 

	dataBuffer := internal.NewAlignedBuffer(512, 512)

	if err := d.doTaskFilecmd(fd, false, false, tf, dataBuffer.GetBuffer(), 3); err != nil {
		return nil, err
	}

	driverHandle := &SgDriverHandle {
		D: d,
		Fd: fd,
	}
	dataBuffer.ResetRead()
	dataBuffer.Read(driverHandle.Identity[:])
	internal.AtaSwapWordEndian(driverHandle.Identity[:])

	return driverHandle, nil
}

func (d *SgDriver) doTaskFilecmd(fd int, rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
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

			// test need modification
			/*if dma {
				cdb.B01 = uint8(internal.Ternary(data != nil, (SG_ATA_PROTO_DMA << 1), (SG_ATA_PROTO_NON_DATA << 1)))
			} else {
				cdb.B01 = uint8(internal.Ternary(data != nil, internal.Ternary(rw, (SG_ATA_PROTO_PIO_OUT << 1), (SG_ATA_PROTO_PIO_IN << 1)), (SG_ATA_PROTO_NON_DATA << 1)))
			}

			if data != nil {
				cdb.B02 |= (SG_ATA_PROTO_DMA << 1) | (SG_ATA_PROTO_NON_DATA << 1)
				cdb.B02 |= uint8(internal.Ternary(rw, (SG_CDB2_TDIR_TO_DEV << 3), (SG_CDB2_TDIR_FROM_DEV << 3)))
			} else {
				cdb.B02 = (SG_CDB2_CHECK_COND << 5)
			}*/

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

		sgParams.DxferLen = uint32(len(data))

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
			copy(buffer[n:], data)

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
				if !rw {
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
	_ = unix.Close(s.Fd)
}

func (s *SgDriverHandle) doTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	return s.D.doTaskFilecmd(s.Fd, rw, dma, tf, data, timeoutSecs)
}

func (s *SgDriverHandle) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return scsiSecurityCommand(s.Fd, rw, dma, protocol, comId, buffer, timeoutSecs)
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
		sgParams.Timeout = uint32(timeoutSecs)
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