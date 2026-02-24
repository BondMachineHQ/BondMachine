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

	// With size specification
	result["^0u<(?P<size>[0-9]+)>(?P<uint>[0-9]+)$"] = unsignedImportWithSize
	result["^0d<(?P<size>[0-9]+)>(?P<uint>[0-9]+)$"] = unsignedImportWithSize

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

// unsignedImportWithSize creates an unsigned number with a specified bit size
// Format: 0u<size>value or 0d<size>value
// Examples: 0u<8>255, 0u<16>1000, 0u<32>12345
func unsignedImportWithSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	sizeStr := re.ReplaceAllString(input, "${size}")
	uintDec := re.ReplaceAllString(input, "${uint}")

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, errors.New("invalid size in " + input)
	}

	if size <= 0 || size > 64 {
		return nil, errors.New("size must be between 1 and 64 bits")
	}

	s, err := strconv.ParseUint(uintDec, 10, 64)
	if err != nil {
		return nil, errors.New("invalid number " + input)
	}

	// Check if value fits in specified size
	if size < 64 {
		maxValue := uint64(1<<uint(size)) - 1
		if s > maxValue {
			return nil, errors.New("value " + uintDec + " exceeds maximum for " + sizeStr + "-bit unsigned")
		}
	}

	newNumber := BMNumber{}
	// Calculate bytes needed (round up to nearest byte)
	bytesNeeded := (size + 7) / 8
	newNumber.number = make([]byte, bytesNeeded)

	mask := uint64(255)
	for i := 0; i < bytesNeeded; i++ {
		newNumber.number[i] = byte(s & mask)
		s = s >> 8
	}

	newNumber.bits = size
	newNumber.nType = Unsigned{}
	return &newNumber, nil
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
