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
