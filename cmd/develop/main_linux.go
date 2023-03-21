//go:build linux
// +build linux

package main

import (
	"github.com/jc-lab/go-dparm"
	"log"
)

// name conflict
func test() {
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
			log.Printf("DRIVE[%d]: %s", i, drive)
		}
	}

	volumes, err := factory.EnumVolumes()
	if err != nil {
		log.Println(err)
	} else {
		for i, drive := range volumes {
			log.Printf("VOLUME[%d]: %s", i, drive)
		}
	}
}
