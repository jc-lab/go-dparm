package windows

type ATA_PASS_THROUGH_DIRECT struct {
	Length             uint16   `struc:"uint16"`
	AtaFlags           uint16   `struc:"uint16"`
	PathId             uint8    `struc:"uint8"`
	TargetId           uint8    `struc:"uint8"`
	Lun                uint8    `struc:"uint8"`
	ReservedAsUchar    uint8    `struc:"uint8"`
	DataTransferLength uint32   `struc:"uint32"`
	TimeOutValue       uint32   `struc:"uint32"`
	ReservedAsUlong    uint32   `struc:"uint32"`
	DataBuffer         uintptr  `struc:"off_t"`
	PreviousTaskFile   [8]uint8 `struc:"[8]uint8"`
	CurrentTaskFile    [8]uint8 `struc:"[8]uint8"`
}
