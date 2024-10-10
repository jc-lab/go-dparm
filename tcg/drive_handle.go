package tcg

import ()

// information returned from tcg level 0 discovery
type TcgLevel0Info struct {
	TcgTper              bool
	TcgLocking           bool
	TcgGeometryReporting bool
	TcgOpalSscV100       bool
	TcgOpalSscV200       bool
	TcgEnterprise        bool
	TcgSingleUser        bool
	TcgDataStore         bool

	TcgRawFeatures map[uint16][]byte
}

type DriveCommandHandler interface {
	SecurityCommand (rw bool, dma bool, protocol uint8, comId uint16, buffer []byte, timeoutSecs int) error
	GetTcgLevel0InfoAndSerial() (TcgLevel0Info, string)
}

type TcgDriveHandle struct {
	DriveCommandHandler
	TcgLevel0Info
	serial string
}

func NewTcgDriveHandle(dc DriveCommandHandler) *TcgDriveHandle {
	h := &TcgDriveHandle{
		DriveCommandHandler: dc,
	}
	h.TcgLevel0Info, h.serial = dc.GetTcgLevel0InfoAndSerial()

	return h
}
