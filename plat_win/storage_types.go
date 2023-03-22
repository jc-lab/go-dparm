package plat_win

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
