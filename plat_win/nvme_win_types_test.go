package plat_win

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_NVMe_COMMAND_DWORD_0_Size(t *testing.T) {
	assert.Equal(t, 4, int(unsafe.Sizeof(NVMe_COMMAND_DWORD_0{})))
}

func Test_NVMe_COMMAND_Size(t *testing.T) {
	assert.Equal(t, 64, int(unsafe.Sizeof(NVMe_COMMAND{})))
}
