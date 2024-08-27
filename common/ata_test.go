package common

import (
	_ "embed"
	"errors"
	"github.com/jc-lab/go-dparm/ata"
	"github.com/stretchr/testify/assert"

	"log"
	"testing"
)

//go:embed sample/intel-01.bin
var intel01AtaIdentify []byte

func testAtaHandle(ataIdentityRaw []byte) DriveHandle {
	ataHandle := &DriveHandleImpl{
		Dh: &PseudoAtaDriver{
			AtaIdentityRaw: ataIdentityRaw,
		},
	}
	if err := ataHandle.Init(); err != nil {
		log.Panicln(err)
	}
	return ataHandle
}

func TestAtaIdentify(t *testing.T) {
	handle := testAtaHandle(intel01AtaIdentify)
	driveInfo := handle.GetDriveInfo()

	assert.Equal(t, "INTEL SSDSC2BA400G3C", driveInfo.Model)
	assert.Equal(t, "BTTV5151046J400HGN", driveInfo.Serial)
	assert.Equal(t, "5DV1FJ03", driveInfo.FirmwareRevision)
	assert.Equal(t, uint64(781422768), driveInfo.AtaIdentity.Max48bitLba)
	assert.True(t, driveInfo.AtaIdentity.CommandSetSupport.GetSmartCommands())

	assert.True(t, driveInfo.AtaIdentity.CommandSetSupport.GetSmartCommands())

	assert.True(t, driveInfo.AtaIdentity.Word59.IsSanitizeFeatureSetSupported())
	assert.True(t, driveInfo.AtaIdentity.Word59.IsBlockEraseExtSupported())
	assert.True(t, driveInfo.AtaIdentity.Word59.IsCryptoScrambleExtSupported())
	assert.False(t, driveInfo.AtaIdentity.Word59.IsSanitizeAntifreezeLockExtSupported())
	assert.False(t, driveInfo.AtaIdentity.Word59.IsOverwriteExtSupported())
}

//

type PseudoAtaDriver struct {
	AtaIdentityRaw []byte
}

func (p *PseudoAtaDriver) GetDriverName() string {
	return "pseudo-ata"
}

func (p *PseudoAtaDriver) GetDrivingType() DrivingType {
	return DrivingAtapi
}

func (p *PseudoAtaDriver) ReopenWritable() error {
	return nil
}

func (p *PseudoAtaDriver) Close() {}

func (p *PseudoAtaDriver) SecurityCommand(rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error {
	return errors.New("not supported")
}

func (p *PseudoAtaDriver) GetIdentity() []byte {
	return p.AtaIdentityRaw
}

func (p *PseudoAtaDriver) DoTaskFileCmd(rw bool, dma bool, tf *ata.Tf, data []byte, timeoutSecs int) error {
	return errors.New("not supported")
}
