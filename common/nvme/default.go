package common_nvme

import (
	"github.com/jc-lab/go-dparm/common"
	"github.com/jc-lab/go-dparm/internal"
	"github.com/jc-lab/go-dparm/nvme"
	"unsafe"
)

func NvmeGetLogPageByAdminPassthru(handle common.NvmeDriverHandle, nsid uint32, logId uint32, rae bool, dataSize int) ([]byte, error) {
	var rootError error

	offset, xferLen := 0, dataSize
	lsp, lpo, lsi := nvme.NVME_NO_LOG_LSP, offset, 0

	dataBuffer := make([]byte, dataSize)

	for {
		if offset >= dataSize {
			return dataBuffer, rootError
		}

		xferLen = dataSize - offset
		if xferLen > 4096 {
			xferLen = 4096
		}

		numd := uint32((dataSize >> 2) - 1)
		numdh := uint32((numd >> 16) & 0xffff)
		numdl := uint32(numd & 0xffff)
		cdw10 := logId | (numdl << 16) | uint32(internal.Ternary(rae, (1<<15), 0)) | (uint32(lsp) << 8)

		cmd := &nvme.NvmeAdminCmd{}
		cmd.Opcode = uint8(nvme.NVME_ADMIN_OP_GET_LOG_PAGE)
		cmd.Nsid = nsid
		cmd.Addr = uintptr(unsafe.Pointer(&dataBuffer[0]))
		cmd.DataLen = uint32(dataSize)
		cmd.Cdw10 = cdw10
		cmd.Cdw11 = numdh | uint32(lsi<<16)
		cmd.Cdw12 = uint32(lpo)
		cmd.Cdw13 = uint32(lpo >> 32)
		cmd.Cdw14 = 0

		rootError = handle.DoNvmeAdminPassthru(cmd)
		if rootError == nil {
			return dataBuffer, nil
		}

		offset += xferLen
	}
}
