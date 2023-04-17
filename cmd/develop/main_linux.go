//go:build linux
// +build linux

package main

import (
	"github.com/jc-lab/go-dparm"
	"log"
)

func main() {
	factory := go_dparm.NewSystemDriveFactory()
	handle, err := factory.OpenByPath("/dev/sda")
	if err != nil {
		log.Println(err)
	} else {
		_ = handle
	}

	drives, err := factory.EnumDrives()
	if err != nil {
		log.Println(err)
	} else {
		for i, drive := range drives {
			log.Printf("DRIVE[%d]: %s %s %s %s", i, drive.Model, drive.Serial, drive.FirmwareRevision, drive.VendorId)
		}
	}

	volumes, err := factory.EnumVolumes()
	if err != nil {
		log.Println(err)
	} else {
		for i, volume := range volumes.GetList() {
			log.Printf("VOLUME[%d]: %s", i, volume)
		}
	}
}
