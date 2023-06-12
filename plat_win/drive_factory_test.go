//go:build windows
// +build windows

package plat_win

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGUID_DEVINTERFACE_DISK(t *testing.T) {
	assert.Equal(t, "{53F56307-B6BF-11D0-94F2-00A0C91EFB8B}", GUID_DEVINTERFACE_DISK.String())
}
