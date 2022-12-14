package bmnumbers

import (
	"errors"
	"regexp"
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
