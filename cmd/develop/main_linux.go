//go:build linux
// +build linux

package main

import (
	"github.com/jc-lab/go-dparm/plat_linux"
	"log"
)

func main() {
	var disk plat_linux.SgDriverHandle
	disk.D = plat_linux.NewSgDriver()
	_, err := disk.D.OpenByPath("/dev/sda")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	log.Printf("Disk info: %v\n", string(disk.Identity[:]))
}
