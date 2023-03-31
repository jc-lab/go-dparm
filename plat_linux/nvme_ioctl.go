package plat_linux

import (
	"unsafe"
)

type NvmeIoctlUserIo struct {
	Opcode uint8
	Flags uint8
	Control uint16
	Nblocks uint16
	Rsvd uint16
	Metadata uint64
	Addr uint64
	Slba uint64
	Dsmgmt uint32
	Reftag uint32
	Apptag uint16
	Appmask uint16
}

type NvmeIoctlPassthruCmd struct {
	Opcode uint8
	Flags uint8
	Rsvd1 uint16
	Nsid uint32
	Cdw2 uint32
	Cdw3 uint32
	Metadata uint64
	Addr uint64
	MetadataLen uint32
	DataLen uint32
	Cdw10 uint32
	Cdw11 uint32
	Cdw12 uint32
	Cdw13 uint32
	Cdw14 uint32
	Cdw15 uint32
	TimeoutMs uint32
	Result uint32
}

type NvmeIoctlAdminCmd NvmeIoctlPassthruCmd

type NvmeIoctlPassthruCmd64 struct {
	Opcode uint8
	Flags uint8
	Rsvd1 uint16
	Nsid uint32
	Cdw2 uint32
	Cdw3 uint32
	Metadata uint64
	Addr uint64
	MetadataLen uint32
	DataLen uint32
	Cdw10 uint32
	Cdw11 uint32
	Cdw12 uint32
	Cdw14 uint32
	Cdw15 uint32
	TimeoutMs uint32
	Rsvd2 uint32
	Result uint64
}

// NVMe ioctl values
var (
	NVME_IOCTL_ID = IO('N', 0x40)
	NVME_IOCTL_ADMIN_CMD = IOWR('N', 0x41, unsafe.Sizeof(NvmeIoctlAdminCmd{}))
	NVME_IOCTL_SUBMIT_IO = IOW('N', 0x42, unsafe.Sizeof(NvmeIoctlUserIo{}))
	NVME_IOCTL_IO_CMD = IOWR('N', 0x43, unsafe.Sizeof(NvmeIoctlPassthruCmd{}))
	NVME_IOCTL_RESET = IO('N', 0x44)
	NVME_IOCTL_SUBSYS_RESET = IO('N', 0x45)
	NVME_RESCAN = IO('N', 0x46)
	NVME_IOCTL_ADMIN64_CMD = IOWR('N', 0x47, unsafe.Sizeof(NvmeIoctlPassthruCmd64{}))
	NVME_IOCTL_IO64_CMD = IOWR('N', 0x48, unsafe.Sizeof(NvmeIoctlPassthruCmd64{}))
)