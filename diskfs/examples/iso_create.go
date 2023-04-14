package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/jc-lab/go-dparm/diskfs"
	"github.com/jc-lab/go-dparm/diskfs/disk"
	"github.com/jc-lab/go-dparm/diskfs/filesystem"
	"github.com/jc-lab/go-dparm/diskfs/filesystem/iso9660"
)

func CreateIso(diskImg string) {
	if diskImg == "" {
		log.Fatal("must have a valid path for diskImg")
	}
	var diskSize int64 = 10 * 1024 * 1024 // 10 MB
	mydisk, err := diskfs.Create(diskImg, diskSize, diskfs.Raw, diskfs.SectorSizeDefault)
	check(err)

	// the following line is required for an ISO, which may have logical block sizes
	// only of 2048, 4096, 8192
	mydisk.LogicalBlocksize = 2048
	fspec := disk.FilesystemSpec{Partition: 0, FSType: filesystem.TypeISO9660, VolumeLabel: "label"}
	fs, err := mydisk.CreateFilesystem(fspec)
	check(err)
	rw, err := fs.OpenFile("demo.txt", os.O_CREATE|os.O_RDWR)
	content := []byte("demo")
	_, err = rw.Write(content)
	check(err)
	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		check(fmt.Errorf("not an iso9660 filesystem"))
	}
	err = iso.Finalize(iso9660.FinalizeOptions{})
	check(err)
}
