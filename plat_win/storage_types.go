package plat_win

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type PartitionStyle uint32

const (
	PartitionStyleMbr PartitionStyle = 0
	PartitionStyleGpt PartitionStyle = 1
	PartitionStyleRaw PartitionStyle = 2
)

const (
	MAX_PARTITIONS = 64
)

const (
	IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS = 0x560000
	IOCTL_DISK_GET_PARTITION_INFO_EX     = 0x00070048
)

type DEVICE_TYPE = uint32
type MEDIA_TYPE = uint32

type STORAGE_DEVICE_NUMBER struct {
	DeviceType      DEVICE_TYPE
	DeviceNumber    uint32
	PartitionNumber uint32
}

type DISK_GEOMETRY struct {
	Cylinders         uint64
	MediaType         MEDIA_TYPE
	TracksPerCylinder uint32
	SectorsPerTrack   uint32
	BytesPerSector    uint32
}

type DISK_GEOMETRY_EX struct {
	Geometry DISK_GEOMETRY
	DiskSize uint64
}

type DRIVE_LAYOUT_INFORMATION_EX_HEADER struct {
	PartitionStyle PartitionStyle
	PartitionCount uint32
}

type DRIVE_LAYOUT_INFORMATION_MBR struct {
	Signature uint32
	CheckSum  uint32
}

type DISK_EXTENT struct {
	DiskNumber     uint32
	StartingOffset uint64
	ExtentLength   uint64
}

type VOLUME_DISK_EXTENTS struct {
	NumberOfDiskExtents uint32
	Extents             [1]DISK_EXTENT
}

type PARTITION_INFORMATION_MBR struct {
	PartitionType       byte
	BootIndicator       bool
	RecognizedPartition bool
	HiddenSectors       uint32
	PartitionId         windows.GUID
}

type PARTITION_INFORMATION_GPT struct {
	PartitionType windows.GUID
	PartitionId   windows.GUID
	Attributes    uint64
	Name          [36]uint16
}

type PARTITION_INFORMATION_EX struct {
	PartitionStyle   PartitionStyle
	StartingOffset   int64
	PartitionLength  int64
	PartitionNumber  int32
	RewritePartition bool
	Rev01            bool
	Rev02            bool
	Rev03            bool
	PartitionInfo    [112]byte
}

func (p *PARTITION_INFORMATION_EX) GetMbr() *PARTITION_INFORMATION_MBR {
	if p.PartitionStyle == PartitionStyleGpt {
		return (*PARTITION_INFORMATION_MBR)(unsafe.Pointer(&p.PartitionInfo[0]))
	}
	return nil
}

func (p *PARTITION_INFORMATION_EX) GetGpt() *PARTITION_INFORMATION_GPT {
	if p.PartitionStyle == PartitionStyleGpt {
		return (*PARTITION_INFORMATION_GPT)(unsafe.Pointer(&p.PartitionInfo[0]))
	}
	return nil
}

type STORAGE_BUS_TYPE byte

const (
	BusTypeUnknown STORAGE_BUS_TYPE = iota + 0
	BusTypeScsi
	BusTypeAtapi
	BusTypeAta
	BusType1394
	BusTypeSsa
	BusTypeFibre
	BusTypeUsb
	BusTypeRAID
	BusTypeiScsi
	BusTypeSas
	BusTypeSata
	BusTypeSd
	BusTypeMmc
	BusTypeVirtual
	BusTypeFileBackedVirtual
	BusTypeSpaces
	BusTypeNvme
	BusTypeSCM
	BusTypeUfs
	BusTypeMax
	BusTypeMaxReserved STORAGE_BUS_TYPE = 0x7F
)

type STORAGE_DEVICE_DESCRIPTOR struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        bool
	CommandQueueing       bool
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               STORAGE_BUS_TYPE
	RawPropertiesLength   uint32
}

func StorageBusTypeToString(busType STORAGE_BUS_TYPE) string {
	switch busType {
	case BusTypeUnknown:
		return "Unknown"
	case BusTypeScsi:
		return "Scsi"
	case BusTypeAtapi:
		return "Atapi"
	case BusTypeAta:
		return "Ata"
	case BusType1394:
		return "1394"
	case BusTypeSsa:
		return "Ssa"
	case BusTypeFibre:
		return "Fibre"
	case BusTypeUsb:
		return "Usb"
	case BusTypeRAID:
		return "RAID"
	case BusTypeiScsi:
		return "iScsi"
	case BusTypeSas:
		return "Sas"
	case BusTypeSata:
		return "Sata"
	case BusTypeSd:
		return "Sd"
	case BusTypeMmc:
		return "Mmc"
	case BusTypeVirtual:
		return "Virtual"
	case BusTypeFileBackedVirtual:
		return "FileBackedVirtual"
	case BusTypeSpaces:
		return "Spaces"
	case BusTypeNvme:
		return "Nvme"
	case BusTypeSCM:
		return "SCM"
	case BusTypeUfs:
		return "Ufs"
	}
	return ""
}
