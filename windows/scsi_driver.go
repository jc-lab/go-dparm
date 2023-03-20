//go:build windows
// +build windows

package windows

import (
	"bytes"
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/lunixbochs/struc"
	"log"
	"syscall"
	"unsafe"
)

type ScsiDriver struct {
	common.Driver
}

type ScsiDriverHandle struct {
	handle syscall.Handle
}

func NewScsiDriver() *ScsiDriver {
	return &ScsiDriver{}
}

func (d *ScsiDriver) OpenByPath(path string) (common.DriveHandle, error) {
	handle, err := openDevice(path)
	if err != nil {
		return nil, err
	}

	driverHandle, err := d.openImpl(handle)
	if err != nil {
		_ = syscall.CloseHandle(handle)
	}
	return driverHandle, err
}

func (d *ScsiDriver) OpenByWindowsPhysicalDrive(path *common.WindowsPhysicalDrive) (common.DriveHandle, error) {
	handle, err := openDevice(path.PhysicalDiskPath)
	if err != nil {
		return nil, err
	}

	driverHandle, err := d.openImpl(handle)
	if err != nil {
		_ = syscall.CloseHandle(handle)
	}
	return driverHandle, err
}

func (d *ScsiDriver) openImpl(handle syscall.Handle) (common.DriveHandle, error) {
	driverHandle := &ScsiDriverHandle{
		handle: handle,
	}

	tf := &ata.Tf{
		Command: 0xec,
	}
	identity := &ata.IdentityDeviceData{}
	dataSize, err := struc.Sizeof(identity)
	if err != nil {
		return nil, err
	}

	dataBuffer := make([]byte, dataSize)
	if err := d.doTaskFileCmd(driverHandle.handle, false, false, tf, dataBuffer, 10); err != nil {
		return nil, err
	}

	if err := struc.Unpack(bytes.NewReader(dataBuffer), &identity); err != nil {
		return nil, err
	}

	return driverHandle, nil
}

func (d *ScsiDriver) doTaskFileCmd(handle syscall.Handle, rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	strucOpts := internal.GetStrucOptions()
	var rootError error = nil

	if rw && data != nil {
		for i := range data {
			data[i] = 0
		}
	}

	scsiParams := &SCSI_PASS_THROUGH_DIRECT_WITH_SENSE_BUF{}
	var payloadBuffer *internal.AlignedBuffer

	for retry := 0; retry < 2; retry++ {
		scsiParams.TimeOutValue = uint32(timeoutSecs)
		scsiParams.DataIn = internal.Ternary(rw, SCSI_IOCTL_DATA_OUT, SCSI_IOCTL_DATA_IN)

		n, err := struc.Sizeof(&scsiParams.SenseInfo)
		if err != nil {
			return err
		}
		scsiParams.SenseInfoLength = byte(n)

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
			n, err = struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			scsiParams.CdbLength = byte(n)

			var buf bytes.Buffer
			err = struc.PackWithOptions(&buf, cdb, strucOpts)
			if err != nil {
				return err
			}
			copy(scsiParams.Cdb[:], buf.Bytes())
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
			n, err = struc.Sizeof(cdb)
			if err != nil {
				return err
			}
			scsiParams.CdbLength = byte(n)

			var buf bytes.Buffer
			err = struc.PackWithOptions(&buf, cdb, strucOpts)
			if err != nil {
				return err
			}
			copy(scsiParams.Cdb[:], buf.Bytes())
		}

		n, err = struc.SizeofWithOptions(&SCSI_PASS_THROUGH_DIRECT{}, strucOpts)
		if err != nil {
			return err
		}
		scsiParams.Length = uint16(n)
		scsiParams.SenseInfoOffset = uint32(n)

		totalSize, err := struc.SizeofWithOptions(scsiParams, strucOpts)
		if err != nil {
			return err
		}

		scsiParams.DataTransferLength = uint32(len(data))
		if retry == 0 {
			scsiParams.DataBuffer = uintptr(unsafe.Pointer(&data[0]))
		} else {
			scsiParams.DataBuffer = uintptr(totalSize) // DataBufferOffset
			totalSize += len(data)
		}

		payloadBuffer = internal.NewAlignedBuffer(4096, totalSize)
		if err = struc.PackWithOptions(payloadBuffer, scsiParams, strucOpts); err != nil {
			return err
		}

		var bytesReturned uint32
		if retry == 0 {
			if err = syscall.DeviceIoControl(
				handle,
				IOCTL_SCSI_PASS_THROUGH_DIRECT,
				(*byte)(unsafe.Pointer(&scsiParams)),
				uint32(payloadBuffer.GetCapacity()),
				(*byte)(unsafe.Pointer(&scsiParams)),
				uint32(unsafe.Sizeof(scsiParams)),
				&bytesReturned,
				nil,
			); err != nil {
				rootError = err
			} else {
				break
			}
		} else {
			payloadBuffer.Write(data)

			if err = syscall.DeviceIoControl(
				handle,
				IOCTL_SCSI_PASS_THROUGH,
				payloadBuffer.GetPointer(),
				uint32(payloadBuffer.GetCapacity()),
				payloadBuffer.GetPointer(),
				uint32(payloadBuffer.GetCapacity()),
				&bytesReturned,
				nil,
			); err != nil {
				rootError = err
			} else {
				rootError = nil
			}
		}
	}

	if rootError == nil {
		// TODO: Fix it, maybe `if (sense_data.SenseKey)` is correct
		payloadBuffer.ResetRead()
		err := struc.UnpackWithOptions(payloadBuffer, scsiParams, strucOpts)
		if err != nil {
			// TODO: Handle it
			log.Println(err)
		} else {
			if scsiParams.SenseInfo.IsValid() && scsiParams.SenseInfo.GetSenseKey() != 0 {
				return &common.DparmError{
					SenseData: &scsiParams.SenseInfo,
				}
			}
		}
	}

	return rootError
}

func (s ScsiDriverHandle) GetDriverName() string {
	//TODO implement me
	panic("implement me")
}

func (s ScsiDriverHandle) MergeDriveInfo(data common.DriveInfo) {
	//TODO implement me
	panic("implement me")
}

func (s ScsiDriverHandle) GetDrivingType() common.DrivingType {
	//TODO implement me
	panic("implement me")
}

func (s ScsiDriverHandle) ReopenWritable() error {
	//TODO implement me
	panic("implement me")
}

func (s ScsiDriverHandle) Close() {
	//TODO implement me
	panic("implement me")
}
