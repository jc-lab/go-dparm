package nvme

import (
	"github.com/jc-lab/go-dparm/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserIoSize(t *testing.T) {
	assert.Equal(t, 44, test.SizeOf(t, &UserIo{}))
}
