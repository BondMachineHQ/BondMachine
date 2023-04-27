package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/x448/float16"
)

type Float16 struct{}

func (d Float16) GetName() string {
	return "float16"
}

func (d Float16) getInfo() string {
	return ""
}

func (d Float16) getSize() int {
	return 16
}

func (d Float16) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0f<16>(?P<number>[^lL].*)$"] = float16Import

	return result
}

func (d Float16) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func float16Import(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")

	if s, err := strconv.ParseFloat(number, 32); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, 2)

		s := uint64(float16.Fromfloat32(float32(s)).Bits())
		mask := uint64(255)
		for i := 0; i < 2; i++ {
			newNumber.number[i] = byte(s & mask)
			s = s >> 8
		}

		newNumber.bits = 16
		newNumber.nType = Float16{}
		return &newNumber, nil

	} else {
		return nil, errors.New("unknown float16 number " + input)
	}
}

func (d Float16) ExportString(n *BMNumber) (string, error) {
	if n.bits != 16 {
		return "", errors.New("cannot export float16 number with " + strconv.Itoa(n.bits) + " bits")
	}

	var s uint16
	for i := 0; i < 2; i++ {
		s = s | (uint16(n.number[i]) << uint16(8*i))
	}

	return "0f<16>" + strconv.FormatFloat(float64(float16.Frombits(s).Float32()), 'f', -1, 32), nil
}

func (d Float16) ShowInstructions() map[string]string {
	result := make(map[string]string)
	result["multop"] = "multf16"
	result["divop"] = "divf16"
	result["addop"] = "addf16"
	return result
}

func (d Float16) ShowPrefix() string {
	return "0f<16>"
}
