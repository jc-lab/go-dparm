package tcg

import (
	"fmt"
)

type TcgTokenVO struct {
	buf []uint8
}

func (p *TcgTokenVO) Buffer() []uint8 {
	return p.buf
}

func (p *TcgTokenVO) Type() Token {
	typeVal := p.buf[0]
	switch {
	case typeVal & 0x80 == 0:
		// tiny atom
		if typeVal & 0x40 == 0 {
			return DTA_TOKENID_SINT
		} else {
			return DTA_TOKENID_UINT
		}
	case typeVal & 0x40 == 0:
		// short atom
		if typeVal & 0x20 != 0 {
			return DTA_TOKENID_BYTESTRING
		} else if typeVal & 0x10 != 0 {
			return DTA_TOKENID_SINT
		} else {
			return DTA_TOKENID_UINT
		}
	case typeVal & 0x20 == 0:
		// medium atom
		if typeVal & 0x10 != 0 {
			return DTA_TOKENID_BYTESTRING
		} else if typeVal & 0x08 != 0 {
			return DTA_TOKENID_SINT
		} else {
			return DTA_TOKENID_UINT
		}
	case typeVal & 0x10 == 0:
		// long atom
		if typeVal & 0x02 != 0 {
			return DTA_TOKENID_BYTESTRING
		} else if typeVal & 0x01 != 0 {
			return DTA_TOKENID_SINT
		} else {
			return DTA_TOKENID_UINT
		}
	default:
		return Token(typeVal)
	}
}

func (p *TcgTokenVO) Length() int {
	return len(p.buf)
}

func (p *TcgTokenVO) GetUint64() (uint64, error) {
	typeVal := p.buf[0]
	switch {
	case typeVal & 0x80 == 0:
		// tiny atom
		if typeVal & 0x40 != 0 {
			// signed atom
			return 0, fmt.Errorf("illegal data")
		} else {
			return uint64(typeVal & 0x3f), nil
		}
	case typeVal & 0x40 == 0:
		// short atom
		if typeVal & 0x10 != 0 {
			return 0, fmt.Errorf("illegal data")
		} else {
			var v uint64
			if len(p.buf) > 9 {
				// error?
			}
			for i, b := uint32(len(p.buf)) - 1, 0; i > 0; i-- {
				v |= uint64(p.buf[i]) << (8 * b)
				b++
			}
			return v, nil
		}
	case typeVal & 0x20 == 0, typeVal& 0x10 == 0:
		// medium atom, long atom
		fallthrough
	default:
		return 0, fmt.Errorf("illegal data")
	}
}

func (p *TcgTokenVO) GetUint32() (uint32, error) {
	dres, err := p.GetUint64()
	return uint32(dres), err
}

func (p *TcgTokenVO) GetUint16() (uint16, error) {
	dres, err := p.GetUint64()
	return uint16(dres), err
}

func (p *TcgTokenVO) GetUint8() (uint8, error) {
	dres, err := p.GetUint64()
	return uint8(dres), err
}

func (p *TcgTokenVO) GetString() (string, error) {
	typeVal := p.buf[0]
	offset := 0

	switch {
	case typeVal & 0x80 == 0:
		// tiny atom
		return "", fmt.Errorf("illegal data")
	case typeVal & 0x40 == 0:
		// short atom
		offset = 1
	case typeVal & 0x20 == 0:
		// medium atom
		offset = 2
	case typeVal & 0x10 == 0:
		// long atom
		offset = 4
	default:
		// non string token
		return "", nil
	}

	return string(p.buf[offset:]), nil
}

func (p *TcgTokenVO) GetBytes() ([]byte, error) {
	typeVal := p.buf[0]
	offset := 0

	switch {
	case typeVal & 0x80 == 0:
		// tiny atom
		return nil, fmt.Errorf("illegal data")
	case typeVal & 0x40 == 0:
		// short atom
		offset = 1
	case typeVal & 0x20 == 0:
		// medium atom
		offset = 2
	case typeVal & 0x10 == 0:
		// long atom
		offset = 4
	default:
		// non bytestring atom
		return nil, nil
	}

	return p.buf[offset:], nil
}
