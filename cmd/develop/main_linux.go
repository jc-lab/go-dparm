//go:build linux
// +build linux

package main

import (
	"encoding/json"
	"log"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	go_dparm "github.com/jc-lab/go-dparm"
	"github.com/jc-lab/go-dparm/nvme"
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
			raw, _ := json.Marshal(volume)
			log.Printf("VOLUME[%d]: %s", i, string(raw))
		}
	}

	// -- NVMe test code --

	handle, err = factory.OpenByPath("/dev/nvme0n1")
	if err != nil {
		log.Println(err)
	} else {
		_ = handle
	}

	info := handle.GetDriveInfo()
	spew.Dump(info)

	logPage, err := handle.NvmeGetLogPage(0xffffffff, uint32(nvme.NVME_GET_LOG_PAGE_SMART), false, int(unsafe.Sizeof(nvme.SmartLogPage{})))
	if err != nil {
		log.Println(err)
	}
	spew.Dump(logPage)
}
