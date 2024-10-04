package tcg

import (
	"errors"
	"unsafe"

	"github.com/jc-lab/go-dparm/internal"
	"github.com/lunixbochs/struc"
)

var (
	ErrInvalidParamType = errors.New("invalid param type")
)

type TcgCommand struct {
	CmdBuf *internal.AlignedBuffer
	Header *OpalHeader
	CmdPtr *uint8
}

type InvokingUid interface {
	[]uint8 | OpalUID
}

func NewTcgCommand() *TcgCommand {
	newCmd := &TcgCommand{
		CmdBuf: internal.NewAlignedBuffer(IO_BUFFER_ALIGNMENT, MAX_BUFFER_LENGTH),
	}
	newCmd.CmdPtr = newCmd.CmdBuf.GetPointer()
	newCmd.Header = (*OpalHeader)(unsafe.Pointer(newCmd.CmdPtr))

	return newCmd
}

func (cmd *TcgCommand) addByteToken(data uint8) error {
	_, err := cmd.CmdBuf.Write([]byte{data})
	return err
}

func (cmd *TcgCommand) Reset() {
	cmd.CmdBuf.ResetWrite()
}

// invokingUid: OpalUID | Buf([]byte)
func (cmd *TcgCommand) Init(invokingUid invokingUID, method OpalMethod) error {
	cmd.Reset()
	if _, err := cmd.CmdBuf.Write([]byte{byte(CALL)}); err != nil {
		return err
	}

	switch uid := invokingUid.(type) {
	case OpalUID:
		cmd.AddToken(uid)
	case Buf:
		cmd.AddRawToken(uid)
	default:
		return ErrInvalidParamType
	}
	cmd.AddToken(method)

	return nil
}

func (cmd *TcgCommand) AddRawToken(data []uint8) error {
	_, err := cmd.CmdBuf.Write(data)

	return err
}

// token: OpalUID | OpalMethod | OpalToken | OpalTinyAtom | OpalShortAtom | OpalLockingState
func (cmd *TcgCommand) AddToken(token token) error {
	switch v := token.(type) {
	case OpalUID:
		return cmd.AddStringToken(string(v[:]), len(v))
	case OpalMethod:
		return cmd.AddStringToken(string(v[:]), len(v))
	case OpalToken:
		return cmd.addByteToken(uint8(v))
	case OpalTinyAtom:
		return cmd.addByteToken(uint8(v))
	case OpalShortAtom:
		return cmd.addByteToken(uint8(v))
	case OpalLockingState:
		return cmd.addByteToken(uint8(v))
	default:
		return ErrInvalidParamType
	}
}

func (cmd *TcgCommand) AddStringToken(text string, length ...int) error {
	var lengthVar int

	if len(length) != 0 {
		lengthVar = length[0]
	} else {
		lengthVar = len(text)
	}

	if lengthVar < 0 {
		lengthVar = len(text)
	}

	if lengthVar == 0 {
		// null token
		_, err := cmd.CmdBuf.Write([]byte{0xa1, 0x00})
		if err != nil {
			return err
		}
	} else if lengthVar < 16 {
		// tiny atom
		_, err := cmd.CmdBuf.Write(append([]uint8{uint8(lengthVar | 0xa0)}, text...))

		if err != nil {
			return err
		}
	} else if lengthVar < 2048 {
		// medium atom
		_, err := cmd.CmdBuf.Write(append([]byte{uint8(0xd0 | (lengthVar>>8)&0x07), uint8(lengthVar)}, text...))

		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *TcgCommand) AddNumberToken(value uint64) error {
	if value < 64 {
		if err := cmd.CmdBuf.WriteByte(uint8(value & 0x3f)); err != nil {
			return err
		}
	} else {
		var startat int
		if value < 0x100 {
			if err := cmd.CmdBuf.WriteByte(0x81); err != nil {
				return err
			}
			startat = 0
		} else if value < 0x10000 {
			if err := cmd.CmdBuf.WriteByte(0x82); err != nil {
				return err
			}
			startat = 1
		} else if value < 0x100000000 {
			if err := cmd.CmdBuf.WriteByte(0x84); err != nil {
				return err
			}
			startat = 3
		} else {
			if err := cmd.CmdBuf.WriteByte(0x88); err != nil {
				return err
			}
			startat = 7
		}
		temp := make([]byte, startat+1)
		for i := startat; i >= 0; i-- {
			temp[startat-i] = uint8(value >> (i * 8))
		}
		_, err := cmd.CmdBuf.Write(temp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *TcgCommand) Complete(eod ...bool) error {
	eodVar := true
	if len(eod) != 0 {
		eodVar = eod[0]
	}

	if eodVar {
		if _, err := cmd.CmdBuf.Write([]uint8{uint8(ENDOFDATA), uint8(STARTLIST), 0x00, 0x00, 0x00, uint8(ENDLIST)}); err != nil {
			return err
		}
	}

	cmd.CmdBuf.ResetRead()
	header := &OpalHeader{}
	if err := struc.Unpack(cmd.CmdBuf, header); err != nil {
		return err
	}

	header.Subpkt.Length = uint32(cmd.CmdBuf.GetPos()) - uint32(unsafe.Sizeof(*header))

	for cmd.CmdBuf.GetPos()&3 != 0 {
		if err := cmd.CmdBuf.WriteByte(0x00); err != nil {
			return err
		}
	}
	header.Pkt.Length = uint32(cmd.CmdBuf.GetPos()) - uint32(unsafe.Sizeof(header.Cp)) - uint32(unsafe.Sizeof(header.Pkt))
	header.Cp.Length = uint32(cmd.CmdBuf.GetPos()) - uint32(unsafe.Sizeof(header.Cp))

	cmd.CmdBuf.ResetWrite()
	if err := struc.Pack(cmd.CmdBuf, header); err != nil {
		return err
	}

	return nil
}

func (cmd *TcgCommand) SetComId(comId uint16) {
	copy(cmd.Header.Cp.ExtendedComID[:], []uint8{uint8(comId >> 8), uint8(comId), 0x00, 0x00})
}

func (cmd *TcgCommand) SetTSN(tsn uint32) {
	cmd.Header.Pkt.Tsn = tsn
}

func (cmd *TcgCommand) SetHSN(hsn uint32) {
	cmd.Header.Pkt.Hsn = hsn
}

func (cmd *TcgCommand) GetCmdPtr() *uint8 {
	return cmd.CmdPtr
}

func (cmd *TcgCommand) GetCmdSize() uint32 {
	x := cmd.CmdBuf.GetPos() & 511
	if x != 0 {
		return uint32(512 - x + cmd.CmdBuf.GetPos())
	}
	return uint32(cmd.CmdBuf.GetPos())
}
