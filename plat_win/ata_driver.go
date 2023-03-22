//go:build windows
// +build windows

package plat_win

import (
	"fmt"
	"github.com/jc-lab/go-dparm/ata"
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	IOCTL_ATA_PASS_THROUGH        = 0x4d02c
	IOCTL_ATA_PASS_THROUGH_DIRECT = 0x4d030
)

const (
	ATA_FLAGS_DRDY_REQUIRED uint16 = 1 << 0
	ATA_FLAGS_DATA_IN       uint16 = 1 << 1
	ATA_FLAGS_DATA_OUT      uint16 = 1 << 2
	ATA_FLAGS_48BIT_COMMAND uint16 = 1 << 3
	ATA_FLAGS_USE_DMA       uint16 = 1 << 4
)

type AtaDriver struct {
	WinDriver
}

type AtaDriverHandle struct {
	common.AtaDriverHandle
	d        *AtaDriver
	handle   windows.Handle
	identity [512]byte
}

func NewAtaDriver() *AtaDriver {
	return &AtaDriver{}
}

func (d *AtaDriver) OpenByHandle(handle windows.Handle) (common.DriverHandle, error) {
	driverHandle, err := d.openImpl(handle)
	return driverHandle, err
}

func (d *AtaDriver) openImpl(handle windows.Handle) (*AtaDriverHandle, error) {
	tf := &ata.Tf{
		Command: ATA_IDENTIFY_DEVICE,
	}
	//tf.Lob.Nsect = 1

	dataBuffer := internal.NewAlignedBuffer(512, 512)
	if err := d.doTaskFileCmd(handle, false, false, tf, dataBuffer.GetBuffer(), 10); err != nil {
		return nil, err
	}

	driverHandle := &AtaDriverHandle{
		d:      d,
		handle: handle,
	}

	dataBuffer.ResetRead()
	dataBuffer.Read(driverHandle.identity[:])
	internal.AtaSwapWordEndian(driverHandle.identity[:])

	return driverHandle, nil
}

func (d *AtaDriver) doTaskFileCmd(handle windows.Handle, rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	var rootError error = nil

	if rw && data != nil {
		for i := range data {
			data[i] = 0
		}
	}

	ataParams := ATA_PASS_THROUGH_DIRECT{}

	for retry := 0; retry < 2; retry++ {
		ataParams.Length = uint16(unsafe.Sizeof(ataParams))
		ataParams.TimeOutValue = uint32(timeoutSecs)
		ataParams.AtaFlags = internal.Ternary(rw, ATA_FLAGS_DATA_OUT, ATA_FLAGS_DATA_IN)

		if tf.IsLba48 != 0 {
			ataParams.AtaFlags |= ATA_FLAGS_48BIT_COMMAND
			ataParams.PreviousTaskFile[0] = tf.Hob.Feat
			ataParams.PreviousTaskFile[1] = tf.Hob.Nsect
			ataParams.PreviousTaskFile[2] = tf.Hob.Lbal
			ataParams.PreviousTaskFile[3] = tf.Hob.Lbam
			ataParams.PreviousTaskFile[4] = tf.Hob.Lbah
			ataParams.PreviousTaskFile[5] = 0
			ataParams.PreviousTaskFile[6] = 0
			ataParams.PreviousTaskFile[7] = 0
		}

		ataParams.CurrentTaskFile[0] = tf.Lob.Feat
		ataParams.CurrentTaskFile[1] = tf.Lob.Nsect
		ataParams.CurrentTaskFile[2] = tf.Lob.Lbal
		ataParams.CurrentTaskFile[3] = tf.Lob.Lbam
		ataParams.CurrentTaskFile[4] = tf.Lob.Lbah
		ataParams.CurrentTaskFile[5] = tf.Dev
		ataParams.CurrentTaskFile[6] = uint8(tf.Command)
		ataParams.CurrentTaskFile[7] = 0 // always zero

		ataParams.DataTransferLength = uint32(len(data))

		var bytesReturned uint32
		if retry == 0 {
			var alignedBuffer *internal.AlignedBuffer = nil
			ataParams.DataBuffer = uintptr(unsafe.Pointer(&data[0]))
			if !internal.IsAlignedPointer(512, ataParams.DataBuffer) {
				alignedBuffer = internal.NewAlignedBuffer(512, len(data))
				if rw {
					alignedBuffer.ResetWrite()
					alignedBuffer.Write(data)
				}
				ataParams.DataBuffer = uintptr(unsafe.Pointer(alignedBuffer.GetPointer()))
			}

			if err := windows.DeviceIoControl(
				handle,
				IOCTL_ATA_PASS_THROUGH_DIRECT,
				(*byte)(unsafe.Pointer(&ataParams)),
				uint32(unsafe.Sizeof(ataParams)),
				(*byte)(unsafe.Pointer(&ataParams)),
				uint32(unsafe.Sizeof(ataParams)),
				&bytesReturned,
				nil,
			); err != nil {
				rootError = err
			} else {
				if !rw && alignedBuffer != nil {
					alignedBuffer.ResetRead()
					alignedBuffer.Read(data)
				}
				break
			}
		} else {
			n := int(unsafe.Sizeof(ataParams))
			ataParams.DataBuffer = uintptr(n)

			srcAtaParams := unsafe.Slice((*byte)(unsafe.Pointer(&ataParams)), unsafe.Sizeof(ataParams))
			buffer := make([]byte, n+len(data))
			copy(buffer, srcAtaParams)
			copy(buffer[n:], data)

			if err := windows.DeviceIoControl(
				handle,
				IOCTL_ATA_PASS_THROUGH,
				&buffer[0],
				uint32(len(buffer)),
				&buffer[0],
				uint32(len(buffer)),
				&bytesReturned,
				nil,
			); err != nil {
				rootError = err
			} else {
				rootError = nil
				if !rw {
					copy(data, buffer[n:])
				}
			}
		}
	}

	if rootError == nil {
		status := ataParams.CurrentTaskFile[6]
		if (status & (0x01 /* ERR */ | 0x08 /* DRQ */)) != 0 {
			return &common.DparmError{
				DriverStatus: status,
				Message:      fmt.Sprintf("ATA Status: %02x", status),
			}
		}
	}

	return rootError
}

func (s *AtaDriverHandle) GetDriverName() string {
	return "WindowsAtaDriver"
}

func (s *AtaDriverHandle) GetDrivingType() common.DrivingType {
	return common.DrivingAtapi
}

func (s *AtaDriverHandle) ReopenWritable() error {
	return nil
}

func (s *AtaDriverHandle) Close() {
	_ = windows.CloseHandle(s.handle)
}

func (s *AtaDriverHandle) GetIdentity() []byte {
	return s.identity[:]
}

func (s *AtaDriverHandle) doTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	return s.d.doTaskFileCmd(s.handle, rw, dma, tf, data, timeoutSecs)
}
