package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type Bin struct{}

func (d Bin) getName() string {
	return "bin"
}

func (d Bin) getInfo() string {
	return ""
}

func (d Bin) getSize() int {
	return -1 // Any size
}

func (d Bin) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0b(?P<bin>[0-1]+)$"] = binImportNoSize
	result["^0b<(?P<bits>[0-9]+)>(?P<bin>[0-1]+)$"] = binImportWithSize

	return result
}

func (d Bin) Convert(n *BMNumber) error {
	convertFrom := n.nType.getName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.getName())
	}
	return nil
}

func binImportNoSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	binNum := re.ReplaceAllString(input, "${bin}")
	var nextByte string
	newNumber := BMNumber{}
	newNumber.number = make([]byte, (len(binNum)-1)/8+1)
	newNumber.bits = len(binNum)
	newNumber.nType = Bin{}

	for i := 0; ; i++ {

		if len(binNum) > 8 {
			nextByte = binNum[len(binNum)-8:]
			binNum = binNum[:len(binNum)-8]
		} else if len(binNum) > 0 {
			nextByte = binNum
			binNum = ""
		} else {
			break
		}

		if val, err := strconv.ParseUint(nextByte, 2, 8); err == nil {
			decoded := byte(val)
			newNumber.number[i] = decoded
		} else {
			return nil, errors.New("invalid binary number" + input)
		}
	}
	return &newNumber, nil
}

func binImportWithSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	binNum := re.ReplaceAllString(input, "${bin}")
	binSizeS := re.ReplaceAllString(input, "${bits}")

	var binSize int

	if val, err := strconv.Atoi(binSizeS); err != nil {
		return nil, errors.New("invalid hex number, wrong bits value" + input)
	} else {
		binSize = val
	}

	if len(binNum) > binSize {
		return nil, errors.New("invalid binary number, the specified number if greater than the bits size " + input)
	}

	var nextByte string
	newNumber := BMNumber{}
	newNumber.number = make([]byte, (binSize-1)/8+1)
	newNumber.bits = binSize
	newNumber.nType = Bin{}

	for i := 0; ; i++ {

		if len(binNum) > 8 {
			nextByte = binNum[len(binNum)-8:]
			binNum = binNum[:len(binNum)-8]
		} else if len(binNum) > 0 {
			nextByte = binNum
			binNum = ""
		} else {
			if i*8 < binSize {
				nextByte = "00000000"
			} else {
				break
			}
		}

		if val, err := strconv.ParseUint(nextByte, 2, 8); err == nil {
			decoded := byte(val)
			newNumber.number[i] = decoded
		} else {
			return nil, errors.New("invalid binary number" + input)
		}
	}
	return &newNumber, nil
}

func (d Bin) ExportString(n *BMNumber) (string, error) {
	return n.ExportBinary(true)
}
