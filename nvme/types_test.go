package nvme

import (
	"github.com/jc-lab/go-dparm/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IdentifyPowerState_Size(t *testing.T) {
	assert.Equal(t, 32, test.SizeOf(t, &IdentifyPowerState{}))
}

func Test_IdentifyController_Size(t *testing.T) {
	assert.Equal(t, 4096, test.SizeOf(t, &IdentifyController{}))
}

func Test_SmartLogPage_Size(t *testing.T) {
	assert.Equal(t, 512, test.SizeOf(t, &SmartLogPage{}))
}
