package tcg

import (
	// "crypto/rand"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/jc-lab/go-dparm/internal"
)

type TcgSession struct {
	tcgDevice     TcgDevice
	sessionOpened bool

	hostSessionNum, tperSessionNum uint64

	noHashPassword, autoClose bool

	timeout uint32
}

func NewTcgSession(tcgDevice TcgDevice) *TcgSession {
	return &TcgSession{
		tcgDevice: tcgDevice,
		autoClose: true,
		timeout:   60000,
	}
}

func (p *TcgSession) Close() {
	if p.autoClose && p.sessionOpened {
		p.sessionOpened = false

		cmd := &TcgCommand{}
		cmd.AddToken(ENDOFSESSION)
		cmd.Complete(false)

		p.SendCommand(cmd)
	}
}

func (p *TcgSession) IsNoHashPassword() bool {
	return p.noHashPassword
}

func (p *TcgSession) SetNoHashPassword(noHash bool) {
	p.noHashPassword = noHash
}

func (p *TcgSession) NoAutoClose() {
	p.autoClose = false
}

func (p *TcgSession) SetTimeout(timeoutMs uint32) {
	p.timeout = timeoutMs
}

// signAuthority: Buf([]uint8) | OpalUID
func (p *TcgSession) Start(sp UID, hostChallenge string, signAuthority signAuthority) error {
	var buf []uint8

	switch v := signAuthority.(type) {
	case UID:
		buf = append([]uint8{uint8(BYTESTRING8)}, v[:]...)
	case Buf:
		buf = v
	}

	isEnterprise := p.tcgDevice.GetDeviceType() == OpalEnterpriseDevice

	// var hostSessionId uint64
	// rand.Read(unsafe.Slice((*byte)(unsafe.Pointer(&hostSessionId)), 8))

	cmd := NewTcgCommand()
	cmd.Init(SMUID_UID, STARTSESSION)
	cmd.AddToken(STARTLIST)
	// cmd.AddNumberToken(hostSessionId)
	cmd.AddNumberToken(105)
	cmd.AddToken(sp)
	cmd.AddToken(UINT_01)
	if hostChallenge != "" && !isEnterprise {
		cmd.AddToken(STARTNAME)
		cmd.AddToken(UINT_00)
		if !p.noHashPassword {
			hashed := TcgHashPassword(p.tcgDevice, false, hostChallenge)
			cmd.AddStringToken(string(hashed))
		} else {
			cmd.AddStringToken(hostChallenge)
		}
		cmd.AddToken(ENDNAME)

		cmd.AddToken(STARTNAME)
		cmd.AddToken(UINT_03)
		cmd.AddRawToken(buf)
		cmd.AddToken(ENDNAME)
	}

	if isEnterprise {
		text := "SessionTimeout"
		cmd.AddToken(STARTNAME)
		cmd.AddStringToken(text)
		cmd.AddNumberToken(uint64(p.timeout))
		cmd.AddToken(ENDNAME)
	}
	cmd.AddToken(ENDLIST)

	cmd.Complete()

	resp, err := p.SendCommand(cmd)
	if err != nil {
		return err
	}

	hsnToken := resp.GetToken(4)
	tsnToken := resp.GetToken(5)
	if hsnToken == nil || tsnToken == nil {
		return ErrIllegalResponse
	}

	temp, err := hsnToken.GetUint32()
	if err != nil {
		return err
	}
	p.hostSessionNum = uint64(binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Pointer(&temp)), 4)))

	temp, err = tsnToken.GetUint32()
	if err != nil {
		return err
	}
	p.tperSessionNum = uint64(binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Pointer(&temp)), 4)))

	if hostChallenge != "" && isEnterprise {
		return p.Authenticate(buf, hostChallenge)
	}

	return nil
}

func (p *TcgSession) Authenticate(authority []uint8, challenge string) error {
	cmd := NewTcgCommand()

	isEnterprise := p.tcgDevice.GetDeviceType() == OpalEnterpriseDevice

	cmd.Init(THISSP_UID, internal.Ternary(isEnterprise, EAUTHENTICATE, AUTHENTICATE))

	cmd.AddToken(STARTLIST)
	cmd.AddRawToken(authority)
	if challenge != "" {
		cmd.AddToken(STARTNAME)
		if isEnterprise {
			text := "Challenge"
			cmd.AddStringToken(text, len(text))
		} else {
			cmd.AddToken(UINT_00)
		}
		if !p.noHashPassword {
			hashed := TcgHashPassword(p.tcgDevice, false, challenge)
			cmd.AddStringToken(string(hashed), len(hashed))
		} else {
			cmd.AddStringToken(challenge, len(challenge))
		}
		cmd.AddToken(ENDNAME)
	}
	cmd.AddToken(ENDLIST)
	cmd.Complete()

	resp, err := p.SendCommand(cmd)
	if err != nil {
		return err
	}

	tempToken := resp.GetToken(1)
	if tempToken == nil {
		return ErrIllegalResponse
	}

	temp, err := tempToken.GetUint8()
	if err != nil {
		return err
	}

	if temp != 0 {
		return fmt.Errorf("authentication failed: %d", temp)
	}

	return nil
}

func (p *TcgSession) SendCommand(cmd *TcgCommand) (*TcgResponse, error) {
	cmd.SetHSN(uint32(p.hostSessionNum))
	cmd.SetTSN(uint32(p.tperSessionNum))
	cmd.SetComId(p.tcgDevice.GetBaseComId())

	resp, err := p.tcgDevice.Exec(cmd, 0x01)
	if err != nil {
		return nil, err
	}

	respHeader := (*TcgHeader)(unsafe.Pointer(resp.GetRespBuf()))
	if respHeader.Cp.Length == 0 || respHeader.Pkt.Length == 0 || respHeader.Subpkt.Length == 0 {
		// payload is not received
		return resp, ErrIllegalResponse
	}

	tempToken := resp.GetToken(0)
	if tempToken == nil {
		return resp, ErrIllegalResponse
	}
	if tempToken.Type() == ENDOFSESSION {
		return resp, nil
	}

	tokenA := resp.GetToken(resp.GetTokenCount() - 1)
	tokenB := resp.GetToken(resp.GetTokenCount() - 5)
	if tokenA == nil || tokenB == nil {
		return resp, ErrIllegalResponse
	}

	if tokenA.Type() != ENDLIST || tokenB.Type() != STARTLIST {
		// no method status
		return resp, ErrIllegalResponse
	}

	tempToken = resp.GetToken(resp.GetTokenCount() - 4)
	methodStatus, err := tempToken.GetUint8()
	if err != nil {
		return resp, err
	}

	if methodStatus != uint8(SUCCESS) {
		return resp, &TcgError{
			Status: MethodStatus(methodStatus),
		}
	}

	return resp, nil
}
