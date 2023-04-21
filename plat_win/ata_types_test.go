package plat_win

import (
	"github.com/jc-lab/go-dparm/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_ATA_PASS_THROUGH_DIRECT_Size(t *testing.T) {
	if internal.IS_64_BIT {
		assert.Equal(t, 48, int(unsafe.Sizeof(ATA_PASS_THROUGH_DIRECT{})))
	} else {
		assert.Equal(t, 40, int(unsafe.Sizeof(ATA_PASS_THROUGH_DIRECT{})))
	}
}
