package plat_win

import (
	"github.com/jc-lab/go-dparm/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ATA_PASSTHROUGH12_Size(t *testing.T) {
	assert.Equal(t, 12, test.SizeOf(t, &ATA_PASSTHROUGH12{}))
}

func Test_ATA_PASSTHROUGH16_Size(t *testing.T) {
	assert.Equal(t, 16, test.SizeOf(t, &ATA_PASSTHROUGH16{}))
}
