package tcg

import (
	"encoding/binary"

	"unsafe"
)

type TcgDeviceOpal1 struct {
	TcgDeviceImpl
}

func (p *TcgDeviceOpal1) IsAnySSC() bool {
	return true
}

func (p *TcgDeviceOpal1) GetDeviceType() TcgDeviceType {
	return OpalV1Device
}

func (p *TcgDeviceOpal1) GetBaseComId() uint16 {
	tcgDh := p.dh.(*TcgDriveHandle)
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcOpalSscV100)]
	if !ok {
		return 0
	}

	feature := (*Discovery0OpalSSCFeatureV100)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.BaseComID)), 2))
}

func (p *TcgDeviceOpal1) GetNumComIds() uint16 {
	tcgDh := p.dh.(*TcgDriveHandle)
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcOpalSscV100)]
	if !ok {
		return 0
	}

	feature := (*Discovery0OpalSSCFeatureV100)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.NumComIDs)), 2))
}

func (p *TcgDeviceOpal1) RevertTPer(password string, isPsid, isAdminSp bool) error {
	return revertTPer(p, password, isPsid)
}

func (p *TcgDeviceOpal1) OpalGetTable(session *TcgSession, table []uint8, startCol, endCol uint16) (*TcgResponse, error) {
	return opalGetTable(session, table, startCol, endCol)
}

func (p *TcgDeviceOpal1) GetDefaultPassword() (string, error) {
	return getDefaultPassword(p)
}

type TcgDeviceOpal2 struct {
	TcgDeviceImpl
}

func (p *TcgDeviceOpal2) IsAnySSC() bool {
	return true
}

func (p *TcgDeviceOpal2) GetDeviceType() TcgDeviceType {
	return OpalV2Device
}

func (p *TcgDeviceOpal2) GetBaseComId() uint16 {
	tcgDh := p.dh.(*TcgDriveHandle)
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcOpalSscV200)]
	if !ok {
		return 0
	}

	feature := (*Discovery0OpalSSCFeatureV200)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.BaseComID)), 2))
}

func (p *TcgDeviceOpal2) GetNumComIds() uint16 {
	tcgDh := p.dh.(*TcgDriveHandle)
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcOpalSscV200)]
	if !ok {
		return 0
	}

	feature := (*Discovery0OpalSSCFeatureV200)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.NumComIDs)), 2))
}

func (p *TcgDeviceOpal2) RevertTPer(password string, isPsid, isAdminSp bool) error {
	return revertTPer(p, password, isPsid)
}

func (p *TcgDeviceOpal2) OpalGetTable(session *TcgSession, table []uint8, startCol, endCol uint16) (*TcgResponse, error) {
	return opalGetTable(session, table, startCol, endCol)
}

func (p *TcgDeviceOpal2) GetDefaultPassword() (string, error) {
	return getDefaultPassword(p)
}

func revertTPer(device TcgDevice, password string, isPsid bool) error {
	sess := NewTcgSession(device)

	uid := SID_UID
	if isPsid {
		sess.SetNoHashPassword(true)
		uid = PSID_UID
	}

	if err := sess.Start(ADMINSP_UID, password, uid); err != nil {
		return err
	}

	cmd := NewTcgCommand()
	cmd.Init(ADMINSP_UID, REVERT)
	cmd.AddToken(STARTLIST)
	cmd.AddToken(ENDLIST)
	cmd.Complete()

	_, err := sess.SendCommand(cmd)
	if err == nil {
		sess.NoAutoClose()
	}

	return err
}

func opalGetTable(session *TcgSession, table []uint8, startCol, endCol uint16) (*TcgResponse, error) {
	cmd := NewTcgCommand()
	cmd.Init(Buf(table), GET)
	cmd.AddToken(STARTLIST)

	cmd.AddToken(STARTLIST)

	cmd.AddToken(STARTNAME)
	cmd.AddToken(STARTCOLUMN)
	cmd.AddNumberToken(uint64(startCol))
	cmd.AddToken(ENDNAME)

	cmd.AddToken(STARTNAME)
	cmd.AddToken(ENDCOLUMN)
	cmd.AddNumberToken(uint64(endCol))
	cmd.AddToken(ENDNAME)

	cmd.AddToken(ENDLIST)

	cmd.AddToken(ENDLIST)
	cmd.Complete()

	return session.SendCommand(cmd)
}

func getDefaultPassword(device TcgDevice) (string, error) {
	session := NewTcgSession(device)
	
	if err := session.Start(ADMINSP_UID, "", UID_HEXFF); err != nil {
		return "", err
	}

	msid := C_PIN_MSID
	table := append([]uint8{uint8(BYTESTRING8)}, msid[:]...)

	resp, err := device.OpalGetTable(session, table, uint16(CREDENTIAL_PIN), uint16(CREDENTIAL_PIN))
	if err != nil {
		return "", err
	}

	passwdToken := resp.GetToken(4)
	if passwdToken == nil {
		return "", ErrIllegalResponse
	}

	return passwdToken.GetString()
}
