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
		if _, err = unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0); err == nil {
			_ = unix.Close(fd)
		}
	}
	return driverHandle, err
}

func (d *SgDriver) openImpl(fd int) (common.DriverHandle, error) {
	tf := &ata.Tf {
		Command: ATA_IDENTIFY_DEVICE,
	}
	//tf.Lob.Nsect = 1 

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
		sgParams.DxferDirection = int32(internal.Ternary(rw, SG_DXFER_FROM_DEV, SG_DXFER_TO_DEV))
		sgParams.Dxferp = uintptr(unsafe.Pointer(&dataBuffer))
		sgParams.MxSbLen = uint8(unsafe.Sizeof(sgParams.SenseData))
		sgParams.Sbp = &sgParams.SenseData[0]

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
			n, err := struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			sgParams.CmdLen = byte(n)

			writer := internal.NewWrappedBuffer([]byte{*sgParams.Cmdp})
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

			n, err := struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			sgParams.CmdLen = byte(n)

			cmdBuffer := make([]byte, sgParams.CmdLen)

			if dma {
				cmdBuffer[1] = SG_ATA_PROTO_DMA // Sense data always included(?)
			} else {
				cmdBuffer[1] = byte(internal.Ternary(rw, SG_ATA_PROTO_PIO_OUT, SG_ATA_PROTO_NON_DATA))
			}

			if sgParams.CmdLen == 16 {
				cmdBuffer[0] = SG_ATA_16
				cmdBuffer[4] = tf.Lob.Feat
				cmdBuffer[6] = tf.Lob.Nsect
				cmdBuffer[8] = tf.Lob.Lbal
				cmdBuffer[10] = tf.Lob.Lbam
				cmdBuffer[12] = tf.Lob.Lbah
				cmdBuffer[13] = tf.Dev
				cmdBuffer[14] = byte(tf.Command)
				if tf.IsLba48 == 1 {
					cmdBuffer[1] |= SG_ATA_LBA48
					cmdBuffer[3] = tf.Hob.Feat
					cmdBuffer[5] = tf.Hob.Nsect
					cmdBuffer[7] = tf.Hob.Lbal
					cmdBuffer[9] = tf.Hob.Lbam
					cmdBuffer[11] = tf.Hob.Lbah
				}
			} else if sgParams.CmdLen == 12 {
				cmdBuffer[0] = SG_ATA_12
				cmdBuffer[3] = tf.Lob.Feat
				cmdBuffer[4] = tf.Lob.Nsect
				cmdBuffer[5] = tf.Lob.Lbal
				cmdBuffer[6] = tf.Lob.Lbam
				cmdBuffer[7] = tf.Lob.Lbah
				cmdBuffer[8] = tf.Dev
				cmdBuffer[9] = byte(tf.Command)
			}

			sgParams.Cmdp = &cmdBuffer[0]

			writer := internal.NewWrappedBuffer([]byte{*sgParams.Cmdp})
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

			if _, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(fd),
				uintptr(SG_IO),
				uintptr(unsafe.Pointer(&sgParams)),
			); err != 0 || sgParams.DriverStatus != 0 {
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
			); err != 0 || sgParams.DriverStatus != 0 {
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
	copyToPointer(unsafe.Pointer(&sgParams), sgParams.SenseData[:], int(unsafe.Sizeof(senseInfo)))

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