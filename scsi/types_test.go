package scsi

import (
	"github.com/jc-lab/go-dparm/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SENSE_DATA_Size(t *testing.T) {
	assert.Equal(t, 18, test.SizeOf(t, &SENSE_DATA{}))
}
