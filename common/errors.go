package common

import (
	"github.com/jc-lab/go-dparm/scsi"
)

type DparmError struct {
	error
	Message      string
	DriverStatus byte
	SenseData    *scsi.SENSE_DATA
}

func (e DparmError) Error() string {
	return e.Message
}
