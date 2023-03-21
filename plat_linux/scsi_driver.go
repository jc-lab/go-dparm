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
	IOCTL_SCSI_PASS_THROUGH        = 0x4D004
	IOCTL_SCSI_PASS_THROUGH_DIRECT = 0x4D014
)

type ScsiDriver struct {
	LinuxDriver
}

type ScsiDriverHandle struct {
	common.AtaDriverHandle
	d *ScsiDriver
	handle int
}

func NewScsiDriver() *ScsiDriver {
	return &ScsiDriver{}
}

func (d *ScsiDriver) OpenByPath(path string) (common.DriveHandle, error) {
	handle, err := OpenDevice(path)
	if err != nil {
		return nil, err
	}

	driverHandle, err := d.openImpl(handle)
	if err != nil {
		_ = unix.Close(handle)
	}
	return driverHandle, err
}

func (d *ScsiDriver) openImpl(handle int) (common.DriveHandle, error) {
	driverHandle := &ScsiDriverHandle {
		handle: handle,
	}

	tf := &ata.Tf {
		Command: ATA_IDENTIFY_DEVICE,
	}
	//tf.Lob.Nsect = 1
	identity := &ata.IdentityDeviceData{}
	dataSize, err := struc.Sizeof(identity)
	if err != nil {
		return nil, err
	}

	dataBuffer := internal.NewAlignedBuffer(512, dataSize)
	if err := d.doTaskFilecmd(driverHandle.handle, false, false, tf, dataBuffer.GetBuffer(), 10); err != nil {
		return nil, err
	}

	if err := struc.Unpack(dataBuffer, &identity); err != nil {
		return nil, err
	}

	return &ScsiDriverHandle{
		d: d,
		handle: handle,
	}, nil
}

func (d *ScsiDriver) doTaskFilecmd(handle int, rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	strucOpts := internal.GetStrucOptions()
	var rootError error = nil

	if rw && data != nil {
		for i := range data {
			data[i] = 0
		}
	}

	scsiParams := SCSI_PASS_THROUGH_DIRECT_WITH_SENSE_BUF{}

	for retry := 0; retry < 2; retry++ {
		scsiParams.Length = uint16(unsafe.Sizeof(SCSI_PASS_THROUGH_DIRECT{}))
		scsiParams.TimeOutValue = uint32(timeoutSecs)
		scsiParams.DataIn = internal.Ternary(rw, SCSI_IOCTL_DATA_OUT, SCSI_IOCTL_DATA_IN)
		scsiParams.SenseInfoLength = byte(unsafe.Sizeof(scsiParams.SenseData))
		scsiParams.SenseInfoOffset = uint32(unsafe.Offsetof(scsiParams.SenseData))

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
		scsiParams.CdbLength = byte(n)

		writer := internal.NewWrappedBuffer(scsiParams.Cdb[:])
		err = struc.PackWithOptions(writer, cdb, strucOpts)
		if err != nil {
			return err
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
			//cdb.SetProtocol(PIO_DATA_IN)
			////cdb.SetTDir(true)
			//cdb.SetByteBlock(true)
			//cdb.SetTLength(2)
			n, err := struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			scsiParams.CdbLength = byte(n)

			writer := internal.NewWrappedBuffer(scsiParams.Cdb[:])
			err = struc.PackWithOptions(writer, cdb, strucOpts)
			if err != nil {
				return err
			}
		}

		scsiParams.DataTransferLength = uint32(len(data))

		if retry == 0 {
			var alignedBuffer *internal.AlignedBuffer = nil
			scsiParams.DataBuffer = uintptr(unsafe.Pointer(&data[0]))
			if !internal.IsAlignedPointer(512, scsiParams.DataBuffer) {
				alignedBuffer = internal.NewAlignedBuffer(512, len(data))
				if rw {
					alignedBuffer.ResetWrite()
					alignedBuffer.Write(data)
				}
				scsiParams.DataBuffer = uintptr(unsafe.Pointer(alignedBuffer.GetPointer()))
			}
			if _, _, err := unix.Syscall(
				unix.SYS_IOCTL,
				uintptr(handle),
				uintptr(IOCTL_SCSI_PASS_THROUGH),
				uintptr(unsafe.Pointer(&scsiParams)),
				) 
				err != 0 {
					rootError = err
				} else {
					if !rw && (alignedBuffer != nil) {
						alignedBuffer.ResetRead()
						alignedBuffer.Read(data)
					}
					break
				}
			} else {
				n := int(unsafe.Sizeof(scsiParams))
				scsiParams.DataBuffer = uintptr(n)

				buffer := make([]byte, n + len(data))
				copyFromPointer(buffer, unsafe.Pointer(&scsiParams), int(unsafe.Sizeof(scsiParams)))
				copy(buffer[n:], data)

				if _, _, err := unix.Syscall(
					unix.SYS_IOCTL,
					uintptr(handle),
					uintptr(IOCTL_SCSI_PASS_THROUGH),
					uintptr(unsafe.Pointer(&buffer)),
				); err != 0 {
					rootError = err
				} else {
					rootError = nil
					copyToPointer(unsafe.Pointer(&scsiParams), buffer[:], int(unsafe.Sizeof(scsiParams)))
					if !rw {
						copy(data, buffer[n:])
					}
				}	
			}
		}
	}

	var senseInfo scsi.SENSE_DATA
	copyToPointer(unsafe.Pointer(&senseInfo), scsiParams.SenseData[:], int(unsafe.Sizeof(senseInfo)))

	if rootError == nil {
		// ?? if scsiParams.SenseInfo.IsValid() && scsiParams.SenseInfo.GetSenseKey() != 0
		if senseInfo.GetSenseKey() != 0 {
			return &common.DparmError{
				DriverStatus: scsiParams.ScsiStatus,
				Message: fmt.Sprintf("SCSI Status: %02x, Sense Key: %#02x, ASC: %#02x, ASCQ: %#02x",
					scsiParams.ScsiStatus,
					senseInfo.GetSenseKey(), senseInfo.AdditionalSenseCode, senseInfo.AdditionalSenseCodeQualifier),
				SenseData: &senseInfo,
			}
		}
	}

	return rootError
}


func (s *ScsiDriverHandle) GetDriverName() string {
	return "LinuxScsiDriver"
}

func (s *ScsiDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingAtapi
}

func (s *ScsiDriverHandle) ReopenWritable() error {
	return nil
}

func (s *ScsiDriverHandle) Close() {
	_ = unix.Close(s.handle)
}

func (s *ScsiDriverHandle) doTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	return s.d.doTaskFilecmd(s.handle, rw, dma, tf, data, timeoutSecs)
}
