package bmnumbers

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
)

type ImportFunc func(*regexp.Regexp, string) (*BMNumber, error)

type BMNumberConfig struct {
	OmitPrefix bool
}

type BMNumberType interface {
	GetName() string
	getInfo() string
	GetSize() int
	importMatchers() map[string]ImportFunc
	Convert(*BMNumber) error
	ExportString(*BMNumber) (string, error)
	ShowInstructions() map[string]string
	ShowPrefix() string
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
	AllTypes = append(AllTypes, Signed{})
	AllTypes = append(AllTypes, Float16{})
	AllTypes = append(AllTypes, Float32{})
	AllTypes = append(AllTypes, Hex{})
	AllTypes = append(AllTypes, Bin{})

	AllDynamicalTypes = make([]DynamicalType, 0)
	AllDynamicalTypes = append(AllDynamicalTypes, DynFloPoCo{})
	dynLQ := DynLinearQuantizer{}
	dynLQRages := make(map[int]LinearDataRange)
	dynLQ.Ranges = &dynLQRages
	AllDynamicalTypes = append(AllDynamicalTypes, dynLQ)
	AllDynamicalTypes = append(AllDynamicalTypes, DynFixedPoint{})

	AllMatchers = make(map[string]ImportFunc)
	for _, t := range AllTypes {
		for k, v := range t.importMatchers() {
			AllMatchers[k] = v
		}
	}

	EventuallyCreateType("flpe4f4", nil)
	EventuallyCreateType("lqs8t0", nil)
	EventuallyCreateType("fps8f4", nil)

}

func ListTypes() {
	for _, t := range AllTypes {
		fmt.Println(t.GetName())
	}
}

func GetType(name string) BMNumberType {
	for _, t := range AllTypes {
		if t.GetName() == name {
			return t
		}
	}
	return nil
}

func (n *BMNumber) GetTypeName() string {
	return n.nType.GetName()
}

func (n *BMNumber) GetBytes() []byte {
	rev := make([]byte, len(n.number))
	copy(rev, n.number)
	slices.Reverse(rev)
	return rev
}

func CastType(n *BMNumber, t BMNumberType) error {
	if n == nil || n.number == nil {
		return errors.New("Cannot cast type of nil number")
	}

	if t.GetSize() != -1 && t.GetSize() != n.bits {
		return errors.New("Cannot cast number of type " + n.nType.GetName() + " with type " + t.GetName() + " because they have different sizes")
	}

	n.nType = t
	return nil
}
