package bmnumbers

import "regexp"

type ImportFunc func(*regexp.Regexp, string) (*BMNumber, error)

type BMNumberType interface {
	getName() string
	importMatchers() map[string]ImportFunc
	convert(*BMNumber) error
}

// BMNumber is a binary representation of a number as a slice of bytes
type BMNumber struct {
	number []byte
	bits   int
	nType  BMNumberType
}

var AllTypes []BMNumberType
var AllMatchers map[string]ImportFunc

func init() {
	AllTypes = make([]BMNumberType, 0)
	AllTypes = append(AllTypes, Unsigned{})
	AllTypes = append(AllTypes, Float32{})
	AllTypes = append(AllTypes, Hex{})
	AllTypes = append(AllTypes, Bin{})

	AllMatchers = make(map[string]ImportFunc)
	for _, t := range AllTypes {
		for k, v := range t.importMatchers() {
			AllMatchers[k] = v
		}
	}
}
