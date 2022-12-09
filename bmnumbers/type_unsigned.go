package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type Unsigned struct{}

func (d Unsigned) getName() string {
	return "unsigned"
}

func (d Unsigned) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^(?P<uint>[0-9]+)$"] = unsignedImportNoSize
	result["^0u(?P<uint>[0-9]+)$"] = unsignedImportNoSize
	result["^0d(?P<uint>[0-9]+)$"] = unsignedImportNoSize

	return result
}

func (d Unsigned) convert(n *BMNumber) error {
	return nil
}

func unsignedImportNoSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	uintDec := re.ReplaceAllString(input, "${uint}")
	if s, err := strconv.ParseUint(uintDec, 10, 0); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, 8)

		mask := uint64(255)
		for i := 0; i < 8; i++ {
			newNumber.number[i] = byte(s & mask)
			s = s >> 8
		}

		newNumber.bits = 64
		newNumber.nType = Unsigned{}
		return &newNumber, nil
	} else {
		return nil, errors.New("invalid number" + input)
	}
}