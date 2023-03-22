package internal

import (
	"encoding/binary"
	"github.com/lunixbochs/struc"
)

const IS_64_BIT = uint64(^uintptr(0)) == ^uint64(0)

func Ternary[T any](statement bool, a T, b T) T {
	if statement {
		return a
	}
	return b
}

func GetStrucOptions() *struc.Options {
	return &struc.Options{
		Order:   binary.LittleEndian,
		PtrSize: Ternary(IS_64_BIT, 64, 32),
	}
}

func GetStrucOptionsWithBigEndian() *struc.Options {
	return &struc.Options{
		Order:   binary.BigEndian,
		PtrSize: Ternary(IS_64_BIT, 64, 32),
	}
}

func AtaSwapWordEndian(input []byte) {
	for i := 0; i < len(input); i += 2 {
		a := input[i]
		input[i] = input[i+1]
		input[i+1] = a
	}
}

func AtaSwap16(x uint16) uint16 {
	return ((x & 0x00ff) << 8) | ((x & 0xff00) >> 8)
}

func AtaSwap32(x uint32) uint32 {
	return (((x) & 0x000000ff) << 24) | (((x) & 0x0000ff00) << 8) | (((x) & 0x00ff0000) >> 8) | (((x) & 0xff000000) >> 24)
}
