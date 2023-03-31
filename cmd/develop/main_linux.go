//go:build linux
// +build linux

package main

import (
	"github.com/jc-lab/go-dparm"
	"log"
)

func main() {
	factory := go_dparm.NewSystemDriveFactory()
	handle, err := factory.OpenByPath("dev/sda")
	if err != nil {
		log.Println(err)
	} else {
		_ = handle
	}
}
