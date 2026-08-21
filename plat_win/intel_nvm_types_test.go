package plat_win

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_INTEL_NVME_PASS_THROUGH_Layout(t *testing.T) {
	v := INTEL_NVME_PASS_THROUGH{}

	assert.Equal(t, 0x1C, int(unsafe.Sizeof(v.SrbIoCtrl)))
	assert.Equal(t, 0x1C, int(unsafe.Offsetof(v.Version)))
	assert.Equal(t, 0x1D, int(unsafe.Offsetof(v.PathId)))
	assert.Equal(t, 0x1E, int(unsafe.Offsetof(v.TargetId)))
	assert.Equal(t, 0x1F, int(unsafe.Offsetof(v.Lun)))
	assert.Equal(t, 0x20, int(unsafe.Offsetof(v.Cmd)))
	assert.Equal(t, 0x40, int(unsafe.Sizeof(v.Cmd)))
	assert.Equal(t, 0x60, int(unsafe.Offsetof(v.IsIOCommandSet)))
	assert.Equal(t, 0x64, int(unsafe.Offsetof(v.CplEntry)))
	assert.Equal(t, 0x74, int(unsafe.Offsetof(v.DataBufferOffset)))
	assert.Equal(t, 0x78, int(unsafe.Offsetof(v.DataBufferLength)))
	assert.Equal(t, 0x7C, int(unsafe.Offsetof(v.Rsvd2)))
	assert.Equal(t, 0xA4, int(unsafe.Offsetof(v.DataBuffer)))
	assert.Equal(t, 0xA4+0x1000, int(unsafe.Sizeof(v)))
}

func Test_INTEL_NVME_PASS_THROUGH_Setters(t *testing.T) {
	v := INTEL_NVME_PASS_THROUGH{}

	v.SetOpcode(NVME_ADMIN_OP_IDENTIFY)
	v.SetFuse(0x02)
	v.SetNsid(0xFFFFFFFF)
	v.SetPrp1(0x1122334455667788)
	v.SetCdw(10, NVME_IDENTIFY_CNS_CONTROLLER)
	v.SetCdw(15, 0xDEADBEEF)

	assert.Equal(t, uint32(0x00000206), v.Cmd[INTEL_NVME_CMD_DW_CDW0])
	assert.Equal(t, uint32(0xFFFFFFFF), v.Cmd[INTEL_NVME_CMD_DW_NSID])
	assert.Equal(t, uint32(0x55667788), v.Cmd[INTEL_NVME_CMD_DW_PRP1_LO])
	assert.Equal(t, uint32(0x11223344), v.Cmd[INTEL_NVME_CMD_DW_PRP1_HI])
	assert.Equal(t, uint32(1), v.GetCdw(10))
	assert.Equal(t, uint32(0xDEADBEEF), v.GetCdw(15))
}

func Test_INTEL_NVME_PASS_THROUGH_CompletionStatus(t *testing.T) {
	v := INTEL_NVME_PASS_THROUGH{}

	assert.Equal(t, uint32(0), v.GetCompletionStatus())

	// DW3: SC=0x02 (Invalid Field) at bit 17~24
	v.CplEntry[3] = 0x02 << 17
	assert.Equal(t, uint32(0x02), v.GetCompletionStatus())
}
