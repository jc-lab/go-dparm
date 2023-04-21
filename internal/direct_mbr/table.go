package direct_mbr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/diskfs/go-diskfs/partition/mbr"

	"github.com/diskfs/go-diskfs/util"
)

// Table represents an MBR partition table to be applied to a disk or read from a disk
type TableEx struct {
	mbr.Table
	MbrIdentifier uint32 // ideitifier of a mbr device #offset(0x01b8~0x01bb)
}

const (
	mbrSize               = 512
	logicalSectorSize     = 512
	physicalSectorSize    = 512
	mbrIdentifierStart    = 440
	mbrIdentifierEnd      = 444
	partitionEntriesStart = 446
	partitionEntriesCount = 4
	signatureStart        = 510
)

// partitionEntrySize standard size of an MBR partition
const partitionEntrySize = 16

func getMbrSignature() []byte {
	return []byte{0x55, 0xaa}
}

// tableFromBytes read a partition table from a byte slice
func tableFromBytes(b []byte) (*TableEx, error) {
	// check length
	if len(b) != mbrSize {
		return nil, fmt.Errorf("data for partition was %d bytes instead of expected %d", len(b), mbrSize)
	}
	mbrSignature := b[signatureStart:]

	// validate signature
	if !bytes.Equal(mbrSignature, getMbrSignature()) {
		return nil, fmt.Errorf("invalid MBR Signature %v", mbrSignature)
	}

	// get mbr identifier
	identity := binary.LittleEndian.Uint32(b[mbrIdentifierStart:mbrIdentifierEnd])

	//parts := make([]*mbr.Partition, 0, partitionEntriesCount)
	//count := int(partitionEntriesCount)
	//for i := 0; i < count; i++ {
	//	// write the primary partition entry
	//	start := partitionEntriesStart + i*partitionEntrySize
	//	end := start + partitionEntrySize
	//	p, err := partitionFromBytes(b[start:end], logicalSectorSize, physicalSectorSize)
	//	if err != nil {
	//		return nil, fmt.Errorf("error reading partition entry %d: %v", i, err)
	//	}
	//	parts = append(parts, p)
	//}

	table := &TableEx{
		Table: mbr.Table{
			LogicalSectorSize:  logicalSectorSize,
			PhysicalSectorSize: 512,
		},
		MbrIdentifier: identity,
	}

	return table, nil
}

// Read read a partition table from a disk, given the logical block size and physical block size
func Read(f util.File, logicalBlockSize, physicalBlockSize int) (*TableEx, error) {
	// read the data off of the disk
	b := make([]byte, mbrSize)
	read, err := f.ReadAt(b, 0)
	if err != nil {
		return nil, fmt.Errorf("error reading MBR from file: %v", err)
	}
	if read != len(b) {
		return nil, fmt.Errorf("read only %d bytes of MBR from file instead of expected %d", read, len(b))
	}
	return tableFromBytes(b)
}
