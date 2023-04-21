package plat_win

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func Test_STORAGE_PROTOCOL_SPECIFIC_DATA_Size(t *testing.T) {
	assert.Equal(t, 40, int(unsafe.Sizeof(STORAGE_PROTOCOL_SPECIFIC_DATA{})))
}

func Test_StorageQueryWithBuffer_Size(t *testing.T) {
	assert.Equal(t, 4144, int(unsafe.Sizeof(StorageQueryWithBuffer{})))
}
