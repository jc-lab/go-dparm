//go:build windows
// +build windows

package plat_win

import (
	"errors"
	"github.com/jc-lab/go-dparm/nvme"
	"golang.org/x/sys/windows"
	"unsafe"
)

func scsiMiniportDoNvmeAdminPassthru(handle windows.Handle, cmd *nvme.NvmeAdminCmd) error {
	nptwb := NVME_PASS_THROUGH_IOCTL{}

	nptwb.SrbIoCtrl.ControlCode = NVME_PASS_THROUGH_SRB_IO_CODE
	nptwb.SrbIoCtrl.HeaderLength = uint32(unsafe.Sizeof(nptwb.SrbIoCtrl))
	copyFromAsciiToBuffer(nptwb.SrbIoCtrl.Signature[:], NVME_SIG_STR)
	nptwb.SrbIoCtrl.Timeout = NVME_PT_TIMEOUT
	nptwb.SrbIoCtrl.Length = uint32(unsafe.Sizeof(nptwb) - unsafe.Sizeof(nptwb.SrbIoCtrl))
	nptwb.DataBufferLen = uint32(unsafe.Sizeof(nptwb.DataBuffer))
	nptwb.ReturnBufferLen = uint32(unsafe.Sizeof(nptwb))
	nptwb.Direction = NVME_FROM_DEV_TO_HOST

	if cmd.DataLen > nptwb.DataBufferLen {
		return errors.New("too long data")
	}
	dataRef := cmd.DataBuffer
	if cmd.DataAddr != 0 {
		dataRef = unsafe.Slice((*byte)(unsafe.Pointer(cmd.DataAddr)), cmd.DataLen)
	}
	copy(nptwb.DataBuffer[:], dataRef)

	pcommand := (*NVMe_COMMAND)(unsafe.Pointer(&nptwb.NVMeCmd))
	pcommand.CDW0.OPC = cmd.Opcode
	pcommand.NSID = cmd.Nsid
	pcommand.CDW10_OR_NDP = cmd.Cdw10
	pcommand.CDW11_OR_NDM = cmd.Cdw11
	pcommand.CDW12 = cmd.Cdw12
	pcommand.CDW13 = cmd.Cdw13
	pcommand.CDW14 = cmd.Cdw14
	pcommand.CDW15 = cmd.Cdw15

	var bytesReturned uint32
	if err := windows.DeviceIoControl(
		handle,
		IOCTL_SCSI_MINIPORT,
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		&bytesReturned,
		nil,
	); err != nil {
		return err
	}

	copy(dataRef, nptwb.DataBuffer[:cmd.DataLen])

	return nil
}
