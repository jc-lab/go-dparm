package plat_win

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
