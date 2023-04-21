package plat_linux

import (
	"github.com/jc-lab/go-dparm/nvme"
)

type NvmeAdminCmdWithBuffer struct {
	Cmd    nvme.NvmeAdminCmd
	Buffer [4096]byte
}
