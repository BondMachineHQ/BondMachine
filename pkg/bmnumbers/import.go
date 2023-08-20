package bmnumbers

import (
	"bufio"
	"errors"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// InitBMNumber creates a new BMNumber for a string
func ImportString(input string) (*BMNumber, error) {

	for k, v := range AllMatchers {
		re := regexp.MustCompile(k)
		if re.MatchString(input) {
			return v(re, input)
		}
	}

	return nil, errors.New("unknown number format " + input)
}

func ImportBytes(input []byte, bits int) (*BMNumber, error) {
	result := new(BMNumber)
	result.nType = Unsigned{}
	result.bits = bits
	result.number = make([]byte, len(input))
	copy(result.number, input)
	slices.Reverse(result.number)
	return result, nil
}

func ImportUint(input interface{}, optionalBits int) (*BMNumber, error) {
	result := new(BMNumber)
	result.nType = Unsigned{}
	switch input.(type) {
	case uint8:
		result.bits = 8
		result.number = make([]byte, 1)
		result.number[0] = input.(uint8)
	case uint16:
		result.bits = 16
		result.number = make([]byte, 2)
		result.number[1] = uint8(input.(uint16) >> 8 & 0xFF)
		result.number[0] = uint8(input.(uint16) & 0xFF)
	case uint32:
		result.bits = 32
		result.number = make([]byte, 4)
		result.number[3] = uint8(input.(uint32) >> 24 & 0xFF)
		result.number[2] = uint8(input.(uint32) >> 16 & 0xFF)
		result.number[1] = uint8(input.(uint32) >> 8 & 0xFF)
		result.number[0] = uint8(input.(uint32) & 0xFF)
	case uint64:
		result.bits = 64
		result.number = make([]byte, 8)
		result.number[7] = uint8(input.(uint64) >> 56 & 0xFF)
		result.number[6] = uint8(input.(uint64) >> 48 & 0xFF)
		result.number[5] = uint8(input.(uint64) >> 40 & 0xFF)
		result.number[4] = uint8(input.(uint64) >> 32 & 0xFF)
		result.number[3] = uint8(input.(uint64) >> 24 & 0xFF)
		result.number[2] = uint8(input.(uint64) >> 16 & 0xFF)
		result.number[1] = uint8(input.(uint64) >> 8 & 0xFF)
		result.number[0] = uint8(input.(uint64) & 0xFF)
	default:
		return nil, errors.New("unknown uint type")
	}
	if optionalBits > 0 {
		// TODO: Finish this
		result.bits = optionalBits
	}
	return result, nil
}

func LoadLinearDataRangesFromFile(filename string) error {

	// Get the linear quantizer ranges struct
	var lqRanges *map[int]LinearDataRange
	for _, t := range AllDynamicalTypes {
		if t.GetName() == "dyn_linear_quantizer" {
			lqRanges = t.(DynLinearQuantizer).Ranges
		}
	}

	splitted := strings.Split(filename, ",")
	if len(splitted)%2 != 0 {
		return errors.New("Error: Invalid linear data range files")
	}

	// Load a file for each index
	for i := 0; i < len(splitted); i += 2 {
		index, err := strconv.Atoi(splitted[i])
		if err != nil {
			return err
		}

		if index == 0 {
			return errors.New("index cannot be 0 (reserved)")
		}

		// Check if the index is already present
		if _, ok := (*lqRanges)[index]; ok {
			return errors.New("index already present")
		}

		filename := splitted[i+1]

		// Read all the lines of the file
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		// Read all the lines of the file
		var max float64
		scanner := bufio.NewScanner(f)
		first := true
		for scanner.Scan() {
			line := scanner.Text()

			// Parse the max values
			if val, err := strconv.ParseFloat(line, 64); err == nil {
				val = math.Abs(val)
				if first {
					max = val
					first = false
				}

				if val > max {
					max = val
				}
			}
		}

		// Add the range to the map
		(*lqRanges)[index] = LinearDataRange{Max: max}
		f.Close()
	}
	return nil
}
