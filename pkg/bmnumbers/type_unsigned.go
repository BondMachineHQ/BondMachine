package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type Unsigned struct{}

func (d Unsigned) GetName() string {
	return "unsigned"
}

func (d Unsigned) getInfo() string {
	return ""
}

func (d Unsigned) GetSize() int {
	return -1 // Any size
}

func (d Unsigned) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^(?P<uint>[0-9]+)$"] = unsignedImportNoSize
	result["^0u(?P<uint>[0-9]+)$"] = unsignedImportNoSize
	result["^0d(?P<uint>[0-9]+)$"] = unsignedImportNoSize
	result["^0u(?P<uint>[0-9]+).0+$"] = unsignedImportNoSize
	result["^0d(?P<uint>[0-9]+).0+$"] = unsignedImportNoSize

	return result
}

func (d Unsigned) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func unsignedImportNoSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	// TODO and with size
	uintDec := re.ReplaceAllString(input, "${uint}")
	if s, err := strconv.ParseUint(uintDec, 10, 0); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, 8)

		mask := uint64(255)
		for i := 0; i < 8; i++ {
			newNumber.number[i] = byte(s & mask)
			s = s >> 8
		}

		newNumber.bits = 64
		newNumber.nType = Unsigned{}
		return &newNumber, nil
	} else {
		return nil, errors.New("invalid number" + input)
	}
}

func (d Unsigned) ExportString(n *BMNumber) (string, error) {
	if result, err := n.ExportUint64(); err == nil {
		return strconv.FormatUint(result, 10), nil
	} else {
		return "", err
	}
}

func (d Unsigned) ShowInstructions() map[string]string {
	result := make(map[string]string)
	result["addop"] = "add"
	result["multop"] = "mult"
	result["divop"] = "div"
	return result
}

func (d Unsigned) ShowPrefix() string {
	return "0u"
}
