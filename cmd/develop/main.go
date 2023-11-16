//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"github.com/jc-lab/go-dparm"
	"log"
)

func main() {
	factory := go_dparm.NewSystemDriveFactory()
	handle, err := factory.OpenByPath("\\\\.\\PhysicalDrive1")
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
			raw, _ := json.Marshal(drive)
			log.Printf("DRIVE[%d]: %s: %s", i, drive.Model, string(raw))
		}
	}

	volumes, err := factory.EnumVolumes()
	if err != nil {
		log.Println(err)
	} else {
		for i, drive := range volumes.GetList() {
			raw, _ := json.Marshal(drive)
			log.Printf("VOLUME[%d]: %s", i, string(raw))
		}
	}
}
