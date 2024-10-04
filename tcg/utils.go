package tcg

import (
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

func TcgHashPassword(device TcgDevice, noHashPassword bool, password string) []uint8 {
	var outHash []uint8

	if len(password) > 32 {
		password = password[:32]
	}

	if noHashPassword {
		outHash = append([]uint8{0xd0, uint8(len(password))}, []uint8(password)...)
	} else {
		driveInfo := device.GetDriveHandle().GetDriveInfo()

		derived := pbkdf2.Key([]byte(password), driveInfo.RawSerial[:], 75000, 32, sha1.New)
		outHash = append([]uint8{0xd0, uint8(len(derived))}, derived...)
	}

	return outHash
}
