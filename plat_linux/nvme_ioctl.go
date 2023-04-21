package plat_linux

import (
	"encoding/binary"

	"github.com/jc-lab/go-dparm/nvme"
)

// NVMe ioctl values
var (
	NVME_IOCTL_ID           = IO('N', 0x40)
	NVME_IOCTL_ADMIN_CMD    = IOWR('N', 0x41, uintptr(binary.Size(nvme.NvmeAdminCmd{})))
	NVME_IOCTL_SUBMIT_IO    = IOW('N', 0x42, uintptr(binary.Size(nvme.UserIo{})))
	NVME_IOCTL_IO_CMD       = IOWR('N', 0x43, uintptr(binary.Size(nvme.PassthruCmd{})))
	NVME_IOCTL_RESET        = IO('N', 0x44)
	NVME_IOCTL_SUBSYS_RESET = IO('N', 0x45)
	NVME_RESCAN             = IO('N', 0x46)
	NVME_IOCTL_ADMIN64_CMD  = IOWR('N', 0x47, uintptr(binary.Size(nvme.PassthruCmd64{})))
	NVME_IOCTL_IO64_CMD     = IOWR('N', 0x48, uintptr(binary.Size(nvme.PassthruCmd64{})))
)
