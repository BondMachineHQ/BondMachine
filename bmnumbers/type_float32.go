package bmnumbers

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

type Float32 struct{}

func (d Float32) getName() string {
	return "unsigned"
}

func (d Float32) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0f(?P<number>[^lL].+)$"] = float32Import

	return result
}

func (d Float32) convert(n *BMNumber) error {
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
