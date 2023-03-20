package main

import (
	"github.com/jc-lab/go-dparm"
	"github.com/jc-lab/go-dparm/windows"
	"log"
)

func main() {
	factory := go_dparm.NewSystemDriveFactory()

	scsiDriver := windows.NewScsiDriver()
	handle, err := scsiDriver.OpenByPath("\\\\.\\PhysicalDrive0")
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
