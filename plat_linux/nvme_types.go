package plat_linux

type UserIo struct {
	Opcode   uint8
	Flags    uint8
	Control  uint16
	Nblocks  uint16
	Rsvd     uint16
	Metadata uint64
	Addr     uint64
	Slba     uint64
	Dsmgmt   uint32
	Reftag   uint32
	Apptag   uint16
	Appmask  uint16
}

type PassthruCmd struct {
	Opcode      uint8
	Flags       uint8
	Rsvd1       uint16
	Nsid        uint32
	Cdw2        uint32
	Cdw3        uint32
	Metadata    uint64
	Addr        uint64
	MetadataLen uint32
	DataLen     uint32
	Cdw10       uint32
	Cdw11       uint32
	Cdw12       uint32
	Cdw13       uint32
	Cdw14       uint32
	Cdw15       uint32
	TimeoutMs   uint32
	Result      uint32
}

type NvmeAdminCmd PassthruCmd

type PassthruCmd64 struct {
	Opcode      uint8
	Flags       uint8
	Rsvd1       uint16
	Nsid        uint32
	Cdw2        uint32
	Cdw3        uint32
	Metadata    uint64
	Addr        uint64
	MetadataLen uint32
	DataLen     uint32
	Cdw10       uint32
	Cdw11       uint32
	Cdw12       uint32
	Cdw13       uint32
	Cdw14       uint32
	Cdw15       uint32
	TimeoutMs   uint32
	Rsvd2       uint32
	Result      uint64
}

type NvmeAdminCmdWithBuffer struct {
	Cmd    NvmeAdminCmd
	Buffer [4096]byte
}
