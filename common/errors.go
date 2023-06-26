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

func (e *DparmError) Error() string {
	return e.Message
}

type NestedError struct {
	error
	message string
	cause   error
}

func NewNestedError(msg string, cause error) *NestedError {
	return &NestedError{
		message: msg + ": " + cause.Error(),
		cause:   cause,
	}
}

func (e *NestedError) Cause() error {
	return e.cause
}

func (e *NestedError) Error() string {
	return e.message
}
