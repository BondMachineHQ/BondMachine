package bmnumbers

import (
	"errors"
	"math"
	"strconv"
)

// InitBMNumber creates a new BMNumber for a string
func ImportNumberFromString(input string) (*BMNumber, error) {
	if len(input) > 2 {
		if input[:2] == "0x" {
			// TODO Hex
		} else if input[:2] == "0b" {
			// TODO Binary
		} else if input[:2] == "0d" {
			// Decimal (also the default)
			if s, err := strconv.Atoi(input[2:]); err == nil {
				newNumber := BMNumber{}
				newNumber.number = make([]uint64, 1)
				newNumber.number[0] = uint64(s)
				return &newNumber, nil
			} else {
				return nil, errors.New("invalid number" + input)
			}
		} else if input[:2] == "0f" {
			// Float32
			if s, err := strconv.ParseFloat(input[2:], 32); err == nil {
				newNumber := BMNumber{}
				newNumber.number = make([]uint64, 1)
				newNumber.number[0] = uint64(math.Float32bits(float32(s)))
				return &newNumber, nil
			} else {
				return nil, errors.New("unknown float32 number " + input)
			}
		} else {
			// Decimal (also the default)
			if s, err := strconv.Atoi(input); err == nil {
				newNumber := BMNumber{}
				newNumber.number = make([]uint64, 1)
				newNumber.number[0] = uint64(s)
				return &newNumber, nil
			} else {
				return nil, errors.New("invalid number" + input)
			}

		}
	} else {
		// Decimal (also the default)
		if s, err := strconv.Atoi(input); err == nil {
			newNumber := BMNumber{}
			newNumber.number = make([]uint64, 1)
			newNumber.number[0] = uint64(s)
			return &newNumber, nil
		} else {
			return nil, errors.New("invalid number" + input)
		}
	}

	return nil, errors.New("unknown number format " + input)
}
