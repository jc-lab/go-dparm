package ata

const lba28_limit = uint64(1)<<28 - 1

func IsNeedsLba48(op OpCode, lba uint64, nsect uint) bool {
	switch op {
	case ATA_OP_DSM:
		fallthrough
	case ATA_OP_READ_PIO_EXT:
		fallthrough
	case ATA_OP_READ_DMA_EXT:
		fallthrough
	case ATA_OP_WRITE_PIO_EXT:
		fallthrough
	case ATA_OP_WRITE_DMA_EXT:
		fallthrough
	case ATA_OP_READ_VERIFY_EXT:
		fallthrough
	case ATA_OP_WRITE_UNC_EXT:
		fallthrough
	case ATA_OP_READ_NATIVE_MAX_EXT:
		fallthrough
	case ATA_OP_SET_MAX_EXT:
		fallthrough
	case ATA_OP_FLUSHCACHE_EXT:
		return true
	case ATA_OP_SECURITY_ERASE_PREPARE:
		fallthrough
	case ATA_OP_SECURITY_ERASE_UNIT:
		fallthrough
	case ATA_OP_VENDOR_SPECIFIC_0x80:
		fallthrough
	case ATA_OP_SMART:
		return false
	}
	if lba >= lba28_limit {
		return true
	}
	if nsect > 0 {
		if nsect > 0xff {
			return true
		}
		if (lba + uint64(nsect) - 1) >= lba28_limit {
			return true
		}
	}
	return false
}

func IsDma(op OpCode) bool {
	switch op {
	case ATA_OP_DSM:
		fallthrough
	case ATA_OP_READ_DMA_EXT:
		fallthrough
	case ATA_OP_READ_FPDMA:
		fallthrough
	case ATA_OP_WRITE_DMA_EXT:
		fallthrough
	case ATA_OP_WRITE_FPDMA:
		fallthrough
	case ATA_OP_READ_DMA:
		fallthrough
	case ATA_OP_WRITE_DMA:
		return true /* SG_DMA */
	}
	return false
}

func TfInit(tf Tf, op OpCode, lba uint64, nsect uint) {
	//memset(tf, 0, sizeof(*tf))
	tf.Command = op
	tf.Dev = ATA_USING_LBA
	tf.Lob.Lbal = uint8(lba)
	tf.Lob.Lbam = uint8(lba >> 8)
	tf.Lob.Lbah = uint8(lba >> 16)
	tf.Lob.Nsect = uint8(nsect)
	if IsNeedsLba48(op, lba, nsect) {
		tf.IsLba48 = 1
		tf.Hob.Nsect = uint8(nsect >> 8)
		tf.Hob.Lbal = uint8(lba >> 24)
		tf.Hob.Lbam = uint8(lba >> 32)
		tf.Hob.Lbah = uint8(lba >> 40)
	} else {
		tf.Dev |= uint8((lba >> 24) & 0x0f)
	}
}

func FixAtaStringOrder(data []byte, trimRight bool) []byte {
	out := make([]byte, len(data))
	outLen := 0

	for i := 0; i < len(data); i += 2 {
		out[i] = data[i+1]
		out[i+1] = data[i]
		if out[i] != 0 {
			outLen++
		} else {
			break
		}
		if out[i+1] != 0 {
			outLen++
		} else {
			break
		}
	}

	if trimRight {
		for (outLen > 0) && (out[outLen-1] == 0x20) {
			outLen--
		}
	}

	return out[:outLen]
}
