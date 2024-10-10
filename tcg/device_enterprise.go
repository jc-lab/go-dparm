package tcg

import (
	"encoding/binary"
	"unsafe"
)

type TcgDeviceEnterprise struct {
	TcgDeviceImpl
}

func (p *TcgDeviceEnterprise) GetDeviceType() TcgDeviceType {
	return OpalEnterpriseDevice
}

func (p *TcgDeviceEnterprise) GetBaseComId() uint16 {
	tcgDh := p.dh
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcEnterprise)]
	if !ok {
		return 0
	}

	feature := (*Discovery0EnterpriseSSCFeature)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.BaseComID)), 2))
}

func (p *TcgDeviceEnterprise) GetNumComIds() uint16 {
	tcgDh := p.dh
	rawBuf, ok := tcgDh.TcgRawFeatures[uint16(FcEnterprise)]
	if !ok {
		return 0
	}

	feature := (*Discovery0EnterpriseSSCFeature)(unsafe.Pointer(&rawBuf[0]))

	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Pointer(&feature.NumberComIDs)), 2))
}

func (p *TcgDeviceEnterprise) RevertTPer(password string, isPsid, isAdminSp bool) error {
	sess := NewTcgSession(p)
	cmd := NewTcgCommand()

	uid := SID_UID
	if isPsid {
		sess.SetNoHashPassword(true)
		uid = PSID_UID
	}

	if err := sess.Start(ADMINSP_UID, password, uid); err != nil {
		return err
	}

	if isAdminSp {
		cmd.Init(ADMINSP_UID, REVERT)
	} else {
		cmd.Init(THISSP_UID, REVERTSP)
	}
	cmd.AddToken(STARTLIST)
	cmd.AddToken(ENDLIST)
	cmd.Complete()

	sess.NoAutoClose()

	_, err := sess.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (p *TcgDeviceEnterprise) EnterpriseGetTable(session *TcgSession, table []uint8, startCol, endCol []uint8) (*TcgResponse, error) {
	cmd := NewTcgCommand()
	cmd.Init(Buf(table), GET)
	cmd.AddToken(STARTLIST)

	cmd.AddToken(STARTLIST)

	cmd.AddToken(STARTNAME)
	cmd.AddStringToken("startColumn")
	cmd.AddStringToken(string(startCol))
	cmd.AddToken(ENDNAME)

	cmd.AddToken(STARTNAME)
	cmd.AddStringToken("endColumn")
	cmd.AddStringToken(string(endCol))
	cmd.AddToken(ENDNAME)

	cmd.AddToken(ENDLIST)

	cmd.AddToken(ENDLIST)
	cmd.Complete()

	return session.SendCommand(cmd)
}

func (p *TcgDeviceEnterprise) GetDefaultPassword() (string, error) {
	session := NewTcgSession(p)

	if err := session.Start(ADMINSP_UID, "", UID_HEXFF); err != nil {
		return "", err
	}

	table := append([]uint8{uint8(BYTESTRING8)}, C_PIN_MSID[:]...)

	resp, err := p.EnterpriseGetTable(session, table, []uint8("PIN"), []uint8("PIN"))
	if err != nil {
		return "", err
	}

	passwdToken := resp.GetToken(4)
	if passwdToken == nil {
		return "", ErrIllegalResponse
	}

	return passwdToken.GetString()
}
