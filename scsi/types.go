package scsi

type SENSE_DATA struct {
	B00                          uint8    `struc:"uint8"`
	SegmentNumber                uint8    `struc:"uint8"`
	B02                          uint8    `struc:"uint8"`
	Information                  [4]uint8 `struc:"[4]uint8"`
	AdditionalSenseLength        uint8    `struc:"uint8"`
	CommandSpecificInformation   [4]uint8 `struc:"[4]uint8"`
	AdditionalSenseCode          uint8    `struc:"uint8"`
	AdditionalSenseCodeQualifier uint8    `struc:"uint8"`
	FieldReplaceableUnitCode     uint8    `struc:"uint8"`
	SenseKeySpecific             [3]uint8 `struc:"[3]uint8"`
}

func (s *SENSE_DATA) GetErrorCode() uint8 {
	return s.B00 & 0x7F
}

func (s *SENSE_DATA) IsValid() bool {
	return (s.B00 & 0x80) != 0
}

func (s *SENSE_DATA) GetSenseKey() uint8 {
	return s.B02 & 0xF0
}

func (s *SENSE_DATA) IsIncorrectLength() bool {
	return (s.B02 & 0x20) != 0
}

func (s *SENSE_DATA) IsEndOfMedia() bool {
	return (s.B02 & 0x60) != 0
}

func (s *SENSE_DATA) IsFileMark() bool {
	return (s.B02 & 0x80) != 0
}
