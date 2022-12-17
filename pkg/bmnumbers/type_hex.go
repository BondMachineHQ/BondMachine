package bmnumbers

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
)

type Hex struct{}

func (d Hex) getName() string {
	return "hex"
}

func (d Hex) getInfo() string {
	return ""
}

func (d Hex) getSize() int {
	return -1 // Any size
}

func (d Hex) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0x(?P<hex>[0-9a-fA-F]+)$"] = hexImportNoSize
	result["^0x<(?P<bits>[0-9]+)>(?P<hex>[0-9a-fA-F]+)$"] = hexImportWithSize

	return result
}

func (d Hex) Convert(n *BMNumber) error {
	convertFrom := n.nType.getName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.getName())
	}
	return nil
}

func hexImportNoSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	hexNum := re.ReplaceAllString(input, "${hex}")
	if len(hexNum)%2 != 0 {
		hexNum = "0" + hexNum
	}

	if decoded, err := hex.DecodeString(hexNum); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, len(decoded))

		for i := 0; i < len(decoded); i++ {
			newNumber.number[i] = decoded[len(decoded)-1-i]
		}

		newNumber.bits = len(decoded) * 8
		newNumber.nType = Hex{}
		return &newNumber, nil
	} else {
		return nil, errors.New("invalid hex number" + input)
	}
}

func hexImportWithSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	hexNum := re.ReplaceAllString(input, "${hex}")
	hexSizeS := re.ReplaceAllString(input, "${bits}")

	var hexSize int

	if val, err := strconv.Atoi(hexSizeS); err != nil {
		return nil, errors.New("invalid hex number, wrong bits value" + input)
	} else {
		if val%8 != 0 {
			return nil, errors.New("invalid hex number, the number of bits as to be a multiple of 8" + input)
		}
		hexSize = val
	}

	if len(hexNum)%2 != 0 {
		hexNum = "0" + hexNum
	}

	if decoded, err := hex.DecodeString(hexNum); err == nil {
		if len(decoded)*8 > hexSize {
			return nil, errors.New("invalid hex number, the specified number if greater than the bits size " + input)
		}

		newNumber := BMNumber{}
		newNumber.number = make([]byte, hexSize)

		for i := 0; i < len(decoded); i++ {
			newNumber.number[i] = decoded[len(decoded)-1-i]
		}

		for i := len(decoded); i < hexSize; i++ {
			newNumber.number[i] = 0
		}

		newNumber.bits = hexSize
		newNumber.nType = Hex{}
		return &newNumber, nil
	} else {
		return nil, errors.New("invalid hex number" + input)
	}
}

func (b Hex) ExportString(n *BMNumber) (string, error) {
	bitS := strconv.Itoa(n.bits)
	result := ""
	for i := len(n.number) - 1; i >= 0; i-- {
		result += hex.EncodeToString([]byte{n.number[i]})
	}

	for len(result) > 1 && result[0] == '0' {
		result = result[1:]
	}

	return "0x<" + bitS + ">" + result, nil
}
