package plat_win

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_DRIVE_LAYOUT_INFORMATION_EX_HEADER_Size(t *testing.T) {
	assert.Equal(t, 8, int(unsafe.Sizeof(DRIVE_LAYOUT_INFORMATION_EX_HEADER{})))
}
