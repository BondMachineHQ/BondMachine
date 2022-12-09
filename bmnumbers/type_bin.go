package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type Bin struct{}

func (d Bin) getName() string {
	return "unsigned"
}

func (d Bin) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0b(?P<bin>[0-1]+)$"] = binImportNoSize
	// result["^0b<(?P<bits>[0-9]+)>(?P<bin>[0-1]+)$"] = binImportWithSize

	return result
}

func (d Bin) convert(n *BMNumber) error {
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
