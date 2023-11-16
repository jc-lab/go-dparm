//go:build windows
// +build windows

package plat_win

import (
	"github.com/jc-lab/go-dparm/common"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

const (
	IOCTL_STORAGE_QUERY_PROPERTY     = 0x2d1400
	IOCTL_STORAGE_GET_DEVICE_NUMBER  = 0x2d1080
	IOCTL_DISK_GET_DRIVE_GEOMETRY_EX = 0x700a0
	IOCTL_DISK_GET_DRIVE_LAYOUT_EX   = 0x00070050
)

type WinBasicInfo struct {
	StorageDeviceNumber *STORAGE_DEVICE_NUMBER
	DiskGeometryEx      *DISK_GEOMETRY_EX
	PartitionStyle      common.PartitionStyle
	MbrSignature        uint32
	GptDiskId           string
	Partitions          []common.Partition
}

func OpenDevice(path string) (windows.Handle, error) {
	sPath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	return windows.CreateFile(
		sPath,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
}

func ReadBasicInfo(handle windows.Handle) *WinBasicInfo {
	result := &WinBasicInfo{}
	storageDeviceNumber := STORAGE_DEVICE_NUMBER{}
	diskGeometryEx := DISK_GEOMETRY_EX{}

	var bytesReturned uint32
	var diskNumber int = -1

	err := windows.DeviceIoControl(
		handle,
		IOCTL_STORAGE_GET_DEVICE_NUMBER,
		nil,
		0,
		(*byte)(unsafe.Pointer(&storageDeviceNumber)),
		uint32(unsafe.Sizeof(storageDeviceNumber)),
		&bytesReturned,
		nil,
	)
	if err == nil {
		diskNumber = int(storageDeviceNumber.DeviceNumber)
		result.StorageDeviceNumber = &storageDeviceNumber
	}

	err = windows.DeviceIoControl(
		handle,
		IOCTL_DISK_GET_DRIVE_GEOMETRY_EX,
		nil,
		0,
		(*byte)(unsafe.Pointer(&diskGeometryEx)),
		uint32(unsafe.Sizeof(diskGeometryEx)),
		&bytesReturned,
		nil,
	)
	if err == nil {
		result.DiskGeometryEx = &diskGeometryEx
	}

	data, err := readDeviceIoControl(handle, IOCTL_DISK_GET_DRIVE_LAYOUT_EX, nil, 0)
	if err == nil {
		header := (*DRIVE_LAYOUT_INFORMATION_EX_HEADER)(unsafe.Pointer(&data[0]))
		next := data[int(unsafe.Sizeof(*header)):]
		entryOffset := GetSizeOf_DRIVE_LAYOUT_INFORMATION()
		entryData := next[entryOffset:]
		entrySize := unsafe.Sizeof(PARTITION_INFORMATION_EX{})

		if header.PartitionStyle == PartitionStyleGpt {
			info := (*DRIVE_LAYOUT_INFORMATION_GPT)(unsafe.Pointer(&next[0]))
			result.PartitionStyle = common.PartitionStyleGpt
			result.GptDiskId = info.DiskId.String()
			result.GptDiskId = strings.Trim(result.GptDiskId, "{}")
			result.GptDiskId = strings.ToLower(result.GptDiskId)

			for i := 0; i < int(header.PartitionCount); i++ {
				if len(entryData) < int(entrySize) {
					break
				}
				partitionEntry := (*PARTITION_INFORMATION_EX)(unsafe.Pointer(&entryData[0]))
				entryData = entryData[entrySize:]

				winGptInfo := partitionEntry.GetGpt()
				partitionImpl := &PartitionImpl{
					DiskExtent: DISK_EXTENT{
						DiskNumber:     uint32(diskNumber),
						StartingOffset: uint64(partitionEntry.StartingOffset),
						ExtentLength:   uint64(partitionEntry.PartitionLength),
					},
					PartitionInfo: *partitionEntry,
					gptInfo: &common.GptPartitionInfo{
						PartitionType: strings.ToUpper(winGptInfo.PartitionType.String()),
						PartitionId:   strings.ToUpper(winGptInfo.PartitionId.String()),
					},
				}
				result.Partitions = append(result.Partitions, partitionImpl)
			}

		} else if header.PartitionStyle == PartitionStyleMbr {
			info := (*DRIVE_LAYOUT_INFORMATION_MBR)(unsafe.Pointer(&next[0]))
			result.PartitionStyle = common.PartitionStyleMbr
			result.MbrSignature = info.Signature

			for i := 0; i < int(header.PartitionCount); i++ {
				if len(entryData) < int(entrySize) {
					break
				}
				partitionEntry := (*PARTITION_INFORMATION_EX)(unsafe.Pointer(&entryData[0]))
				entryData = entryData[entrySize:]

				winMbrInfo := partitionEntry.GetMbr()
				partitionImpl := &PartitionImpl{
					DiskExtent: DISK_EXTENT{
						DiskNumber:     uint32(diskNumber),
						StartingOffset: uint64(partitionEntry.StartingOffset),
						ExtentLength:   uint64(partitionEntry.PartitionLength),
					},
					PartitionInfo: *partitionEntry,
					mbrInfo: &common.MbrPartitionInfo{
						PartitionType: winMbrInfo.PartitionType,
						BootIndicator: winMbrInfo.BootIndicator,
					},
				}
				result.Partitions = append(result.Partitions, partitionImpl)
			}
		}
	}

	return result
}

type StorageDeviceDescription struct {
	DeviceType         byte
	DeviceTypeModifier byte
	RemovableMedia     bool
	VendorId           string
	ProductId          string
	ProductRevision    string
	SerialNumber       string
	BusType            STORAGE_BUS_TYPE
}

func ReadStorageQuery(handle windows.Handle) (*StorageDeviceDescription, error) {
	query := STORAGE_PROPERTY_QUERY_WITH_DUMMY{}
	query.QueryType = PropertyStandardQuery
	query.PropertyId = StorageDeviceProperty

	buffer, err := readDeviceIoControl(
		handle,
		IOCTL_STORAGE_QUERY_PROPERTY,
		(*byte)(unsafe.Pointer(&query)),
		uint32(unsafe.Sizeof(query)),
	)
	if err != nil {
		return nil, err
	}

	resp := (*STORAGE_DEVICE_DESCRIPTOR)(unsafe.Pointer(&buffer[0]))

	// STORAGE_DEVICE_DESCRIPTOR

	result := &StorageDeviceDescription{
		DeviceType:         resp.DeviceType,
		DeviceTypeModifier: resp.DeviceTypeModifier,
		RemovableMedia:     resp.RemovableMedia,
		BusType:            resp.BusType,
	}

	result.VendorId = strings.Trim(readNullTerminatedAscii(buffer, int(resp.VendorIdOffset)), " ")
	result.ProductId = strings.Trim(readNullTerminatedAscii(buffer, int(resp.ProductIdOffset)), " ")
	result.SerialNumber = strings.Trim(readNullTerminatedAscii(buffer, int(resp.SerialNumberOffset)), " ")
	result.ProductRevision = strings.Trim(readNullTerminatedAscii(buffer, int(resp.ProductRevisionOffset)), " ")

	return result, nil
}

func readNullTerminatedAscii(buf []byte, offset int) string {
	if offset <= 0 {
		return ""
	}
	buf = buf[offset:]
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return string(buf[:i])
		}
	}
	return ""
}

func readDeviceIoControl(handle windows.Handle, ioctl uint32, inBuffer *byte, inSize uint32) ([]byte, error) {
	var bytesReturned uint32

	buffer := make([]byte, 4096)
	err := windows.DeviceIoControl(handle, ioctl, inBuffer, inSize, &buffer[0], uint32(len(buffer)), &bytesReturned, nil)
	errno, ok := err.(syscall.Errno)
	if ok && errno == windows.ERROR_INSUFFICIENT_BUFFER {
		buffer = make([]byte, bytesReturned)
		err = windows.DeviceIoControl(handle, ioctl, inBuffer, inSize, &buffer[0], uint32(len(buffer)), &bytesReturned, nil)
	}
	if err == nil {
		return buffer[:bytesReturned], nil
	}

	return nil, errno
}


func copyToPointer(dest unsafe.Pointer, src []byte, len int) {
	destRef := unsafe.Slice((*byte)(dest), len)
	copy(destRef, src[:len])
}

func copyFromPointer(dest []byte, src unsafe.Pointer, len int) {
	srcRef := unsafe.Slice((*byte)(src), len)
	copy(dest, srcRef[:len])
}

func copyFromAsciiToBuffer(dest []byte, text string) {
	c := len(text)
	for i := 0; i < c; i++ {
		dest[i] = text[i]
	}
}

func zerofill(buf []uint16) {
	for i := range buf {
		buf[i] = 0
	}
}

func wcslen(buf []uint16) int {
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return i
		}
	}
	return len(buf)
}
