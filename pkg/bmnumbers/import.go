package bmnumbers

import (
	"bufio"
	"errors"
	"math"
	"os"
	"regexp"
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
