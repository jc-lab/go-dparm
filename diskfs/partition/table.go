package partition

import (
	"github.com/jc-lab/go-dparm/diskfs/partition/part"
	"github.com/jc-lab/go-dparm/diskfs/util"
)

// Table reference to a partitioning table on disk
type Table interface {
	Type() string
	Write(util.File, int64) error
	GetPartitions() []part.Partition
}
