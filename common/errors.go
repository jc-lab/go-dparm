package common

import (
	"github.com/jc-lab/go-dparm/scsi"
)

type DparmError struct {
	error
	SenseData *scsi.SENSE_DATA
}
