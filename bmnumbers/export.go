package bmnumbers

import (
	"errors"
	"strconv"
)

func (n *BMNumber) ExportUint64() (uint64, error) {
	if n == nil || n.number == nil || len(n.number) > 8 {
		return 0, errors.New("number cannot be exported as uint64")
	}

	result := uint64(0)

	for i := 0; i < len(n.number); i++ {
		result += uint64(n.number[i]) << (8 * uint(i))
	}

	return result, nil
}

func (n *BMNumber) ExportBinary(withSize bool) (string, error) {
	if n == nil || n.number == nil {
		return "", errors.New("undefined number")
	}

	result := ""

	for _, number := range n.number {
		dataVal := "00000000" + strconv.FormatUint(uint64(number), 2)
		result = dataVal[len(dataVal)-8:] + result
	}

	for len(result) > 0 && result[0] == '0' {
		result = result[1:]
	}

	if withSize {
		result = "0b<" + strconv.Itoa(n.bits) + ">" + result
	}

	return result, nil
}
