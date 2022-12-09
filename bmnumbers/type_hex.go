package bmnumbers

import (
	"encoding/hex"
	"errors"
	"regexp"
)

type Hex struct{}

func (d Hex) getName() string {
	return "hex"
}

func (d Hex) getInfo() string {
	return ""
}

func (d Hex) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0x(?P<hex>[0-9a-fA-F]+)$"] = hexImportNoSize
	// result["^0x<(?P<bits>[0-9]+)>(?P<hex>[0-9a-fA-F]+)$"] = hexImportWithSize

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
