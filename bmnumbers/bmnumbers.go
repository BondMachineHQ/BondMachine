package bmnumbers

import (
	"errors"
	"fmt"
	"regexp"
)

type ImportFunc func(*regexp.Regexp, string) (*BMNumber, error)

type BMNumberType interface {
	getName() string
	getInfo() string
	getSize() int
	importMatchers() map[string]ImportFunc
	Convert(*BMNumber) error
	ExportString(*BMNumber) (string, error)
}

// BMNumber is a binary representation of a number as a slice of bytes
type BMNumber struct {
	number []byte
	bits   int
	nType  BMNumberType
}

var AllTypes []BMNumberType
var AllMatchers map[string]ImportFunc
var AllDynamicalTypes []DynamicalType

func init() {
	AllTypes = make([]BMNumberType, 0)
	AllTypes = append(AllTypes, Unsigned{})
	AllTypes = append(AllTypes, Float32{})
	AllTypes = append(AllTypes, Hex{})
	AllTypes = append(AllTypes, Bin{})

	AllDynamicalTypes = make([]DynamicalType, 0)
	AllDynamicalTypes = append(AllDynamicalTypes, DynFloPoCo{})

	AllMatchers = make(map[string]ImportFunc)
	for _, t := range AllTypes {
		for k, v := range t.importMatchers() {
			AllMatchers[k] = v
		}
	}

	EventuallyCreateType("flpe4f4", nil)

}

func ListTypes() {
	for _, t := range AllTypes {
		fmt.Println(t.getName())
	}
}

func GetType(name string) BMNumberType {
	for _, t := range AllTypes {
		if t.getName() == name {
			return t
		}
	}
	return nil
}

func OverrideType(n *BMNumber, t BMNumberType) error {
	if n == nil || n.number == nil {
		return errors.New("Cannot override type of nil number")
	}

	if t.getSize() != -1 && t.getSize() != n.bits {
		return errors.New("Cannot override number of type " + n.nType.getName() + " with type " + t.getName() + " because they have different sizes")
	}

	n.nType = t
	return nil
}
