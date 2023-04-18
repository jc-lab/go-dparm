//go:build linux
// +build linux

package main

import (
	"log"

	"github.com/jc-lab/go-dparm"
	// "github.com/jc-lab/go-dparm/plat_linux"
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

	// ---- testing nvme passthrough command - only if nvme device exists ----

	/* nvmeDrive := plat_linux.NewLinuxNvmeDriver()
	nvmeFd, err := plat_linux.OpenDevice("/dev/nvme0n1")
	if err != nil {
		log.Println(err)
	}

	nvmeHandle, err := nvmeDrive.OpenByFd(nvmeFd)
	if err != nil {
		log.Fatalln(err)
	}

	nvmeInfo := nvmeHandle.(*plat_linux.LinuxNvmeDriverHandle).GetIdentity()
	log.Printf("Incoming data: %v\n", nvmeInfo) */
}
