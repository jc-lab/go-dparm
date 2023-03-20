package windows

import (
	"github.com/jc-lab/go-dparm/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_SCSI_PASS_THROUGH_DIRECT_Size(t *testing.T) {
	if internal.IS_64_BIT {
		assert.Equal(t, 56, int(unsafe.Sizeof(SCSI_PASS_THROUGH_DIRECT{})))
	} else {
		assert.Equal(t, 44, int(unsafe.Sizeof(SCSI_PASS_THROUGH_DIRECT{})))
	}
}
