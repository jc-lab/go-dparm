package plat_linux

import (
	"unsafe"
)

// NVMe ioctl values
var (
	NVME_IOCTL_ID           = IO('N', 0x40)
	NVME_IOCTL_ADMIN_CMD    = IOWR('N', 0x41, unsafe.Sizeof(NvmeAdminCmd{}))
	NVME_IOCTL_SUBMIT_IO    = IOW('N', 0x42, unsafe.Sizeof(UserIo{}))
	NVME_IOCTL_IO_CMD       = IOWR('N', 0x43, unsafe.Sizeof(PassthruCmd{}))
	NVME_IOCTL_RESET        = IO('N', 0x44)
	NVME_IOCTL_SUBSYS_RESET = IO('N', 0x45)
	NVME_RESCAN             = IO('N', 0x46)
	NVME_IOCTL_ADMIN64_CMD  = IOWR('N', 0x47, unsafe.Sizeof(PassthruCmd64{}))
	NVME_IOCTL_IO64_CMD     = IOWR('N', 0x48, unsafe.Sizeof(PassthruCmd64{}))
)
