package bmnumbers

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

type Float32 struct{}

func (d Float32) getName() string {
	return "float32"
}

func (d Float32) getInfo() string {
	return ""
}

func (d Float32) getSize() int {
	return 32
}

func (d Float32) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0f(?P<number>[^lL].*)$"] = float32Import

	return result
}

func (d Float32) Convert(n *BMNumber) error {
	convertFrom := n.nType.getName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.getName())
	}
	return nil
}

func float32Import(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")

	if s, err := strconv.ParseFloat(number, 32); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, 4)

		s := uint64(math.Float32bits(float32(s)))
		mask := uint64(255)
		for i := 0; i < 4; i++ {
			newNumber.number[i] = byte(s & mask)
			s = s >> 8
		}

		newNumber.bits = 32
		newNumber.nType = Float32{}
		return &newNumber, nil

	} else {
		return nil, errors.New("unknown float32 number " + input)
	}
}

func (d Float32) ExportString(n *BMNumber) (string, error) {
	if n.bits != 32 {
		return "", errors.New("cannot export float32 number with " + strconv.Itoa(n.bits) + " bits")
	}

	var s uint32
	for i := 0; i < 4; i++ {
		s = s | (uint32(n.number[i]) << uint32(8*i))
	}

	return "0f<32>" + strconv.FormatFloat(float64(math.Float32frombits(s)), 'f', -1, 32), nil
}

func (d Float32) ShowInstructions() map[string]string {
	result := make(map[string]string)
	result["multop"] = "multf"
	result["divop"] = "divf"
	result["addop"] = "addf"
	return result
}
