package tcg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"github.com/jc-lab/go-dparm/internal"
)

var (
	ErrIllegalResponse = errors.New("illegal response")
)

type TcgResponse struct {
	buf    *internal.AlignedBuffer
	header *OpalHeader
	ptr    *uint8
	tokens []*TcgTokenVO
}

func NewTcgResponse() *TcgResponse {
	newResp :=  &TcgResponse{
		buf: internal.NewAlignedBuffer(IO_BUFFER_ALIGNMENT, MIN_BUFFER_LENGTH),
	}
	newResp.ptr = newResp.buf.GetPointer()
	newResp.header = (*OpalHeader)(unsafe.Pointer(newResp.ptr))

	return newResp
}

func (p *TcgResponse) Reset() {
	p.buf.Reset()
}

func (p *TcgResponse) GetRespBuf() *uint8 {
	return p.ptr
}

func (p *TcgResponse) GetRespBufSize() uint32 {
	return MIN_BUFFER_LENGTH
}

func (p *TcgResponse) Commit() error {
	var respTokens []*TcgTokenVO

	subpktLen := binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Pointer(&p.header.Subpkt.Length)), 4))
	cur := (*uint8)(unsafe.Add(unsafe.Pointer(p.ptr), unsafe.Sizeof(*p.header)))
	end := (*uint8)(unsafe.Add(unsafe.Pointer(cur), subpktLen))

	var curTokenLen uint32
	var tempTokenBuf []uint8

	if uintptr(unsafe.Pointer(end)) > uintptr(unsafe.Pointer(p.ptr)) + MIN_BUFFER_LENGTH {
		return fmt.Errorf("illegal data")
	}

	for uintptr(unsafe.Pointer(cur)) < uintptr(unsafe.Pointer(end)) {
		switch {
		case *cur & 0x80 == 0:
			// tiny atom
			curTokenLen = 1
		case *cur & 0x40 == 0:
			// short atom
			curTokenLen = uint32(*cur & 0x0f) + 1
		case *cur & 0x20 == 0:
			// medium atom
			curTokenLen = (((uint32(*cur) & 0x07) << 8) | (uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(cur), 1))))) + 2
		case *cur & 0x10 == 0:
			// long atom
			curTokenLen = ((uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(cur), 1))) << 16) | (uint32((*(*uint8)(unsafe.Add(unsafe.Pointer(cur), 2)))) << 8) | uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(cur), 3)))) + 4
		default:
			// token
			curTokenLen = 1
		}

		tempTokenBuf = unsafe.Slice(cur, curTokenLen)
		cur = (*uint8)(unsafe.Add(unsafe.Pointer(cur), curTokenLen))

		if len(tempTokenBuf) != 1 || tempTokenBuf[0] != uint8(EMPTYATOM) {
			respTokens = append(respTokens, &TcgTokenVO{
				buf: tempTokenBuf,
			})
		}
	}

	p.tokens = respTokens

	return nil
}

func (p *TcgResponse) GetTokenCount() int {
	return len(p.tokens)
}

func (p *TcgResponse) GetToken(index int) *TcgTokenVO {
	if index >= len(p.tokens) {
		return nil
	}
	return p.tokens[index]
}
