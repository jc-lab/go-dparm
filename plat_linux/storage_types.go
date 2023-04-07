package plat_linux

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
	GDIO_DRIVE_CMD = 0x031f
	GDIO_DRIVE_RESET = 0x031c
	GDIO_DRIVE_TASK = 0x031e
	GDIO_DRIVE_TASKFILE = 0x031d
	GDIO_GETGEO = 0x0301
	GDIO_GETGEO_BIG = 0x0330
	GDIO_GET_32BIT = 0x0309
	GDIO_GET_ACOUSTIC = 0x030f
	GDIO_GET_BUSSTATE = 0x031a
	GDIO_GET_DMA = 0x030b
	GDIO_GET_IDENTITY = 0x030d
	GDIO_GET_KEEPSETTINGS = 0x0308
	GDIO_GET_MULTCOUNT = 0x0304
	GDIO_GET_NOWERR = 0x030a
	GDIO_GET_QDMA = 0x0305
	GDIO_GET_UNMASKINTR = 0x0302
	GDIO_OBSOLETE_IDENTITY = 0x0307
	GDIO_SCAN_HWIF = 0x0328
	GDIO_SET_32BIT = 0x0324
	GDIO_SET_ACOUSTIC = 0x032c
	GDIO_SET_BUSSTATE = 0x032d
	GDIO_SET_DMA = 0x0326
	GDIO_SET_KEEPSETTINGS = 0x0323
	GDIO_SET_MULTICOUNT = 0x0321
	GDIO_SET_NOWERR = 0x0325
	GDIO_SET_PIO_MODE = 0x0327
	GDIO_SET_QDMA = 0x032e
	GDIO_SET_UNMASKINTR = 0x0322
	GDIO_SET_WCACHE = 0x032b
	GDIO_TRISTATE_HWIF = 0x031b
	GDIO_UNRESISTER_HWIF = 0x032a
	CDROM_SPEED = 0x5322
)

var (
	BLKGETSIZE64 = IOR(uintptr(0x12), uintptr(114), 64)
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
