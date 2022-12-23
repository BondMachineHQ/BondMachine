package etherbond

import (
	"errors"
)

// Encoding

func packCommon(rsize uint8, id uint8) ([]byte, error) {
	var size int
	switch id {
	case ADV_CLU_FR:
		size = 23
	case ADV_CH_FR:
		size = 27
	case ADV_IN_FR:
		size = 27
	case ADV_OUT_FR:
		size = 27
	case IO_TR_FR:
		served := 1
		for bytes := 1; bytes < 16; bytes++ {
			if served<<uint8(bytes) >= int(rsize) {
				size = 31 + int(bytes)
				break
			}
		}
	case ACK_FR:
		size = 19
	}
	if size < 14 {
		return nil, errors.New("Packet to short")
	}

	buf := make([]byte, size)
	return buf, nil
}

func pcat(lbuf []byte, buf []byte) []byte {
	size := len(lbuf)
	for i, val := range lbuf {
		buf[i] = val
	}

	return buf[size:]
}

func pint8(val uint8, buf []byte) []byte {
	buf[0] = val
	return buf[1:]
}

func pint16(val uint16, buf []byte) []byte {
	buf[1] = uint8(val)
	buf[0] = uint8(val >> 8)
	return buf[2:]
}

func pint32(val uint32, buf []byte) []byte {
	buf[3] = uint8(val)
	buf[2] = uint8(val >> 8)
	buf[1] = uint8(val >> 16)
	buf[0] = uint8(val >> 24)
	return buf[4:]
}

func pint64(val uint64, buf []byte) []byte {
	buf[7] = uint8(val)
	buf[6] = uint8(val >> 8)
	buf[5] = uint8(val >> 16)
	buf[4] = uint8(val >> 24)
	buf[3] = uint8(val >> 32)
	buf[2] = uint8(val >> 40)
	buf[1] = uint8(val >> 48)
	buf[0] = uint8(val >> 56)
	return buf[8:]
}

// Decoding

func gcat(n int, buf []byte) ([]byte, []byte) {
	return buf[:n], buf[n:]
}

func gint8(buf []byte) (uint8, []byte) { return buf[0], buf[1:] }

func gint16(buf []byte) (uint16, []byte) {
	return uint16(buf[0]) << 8 | (uint16(buf[1])), buf[2:]
}

func gint32(buf []byte) (uint32, []byte) {
	return uint32(buf[0]) << 24 | (uint32(buf[1]) << 16) | (uint32(buf[2]) << 8) |
			(uint32(buf[3])),
		buf[4:]
}
