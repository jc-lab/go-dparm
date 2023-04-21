package scsi

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_SENSE_DATA_Size(t *testing.T) {
	assert.Equal(t, 18, int(unsafe.Sizeof(SENSE_DATA{})))
}
