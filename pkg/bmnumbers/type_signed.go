package bmnumbers

import (
	"errors"
	"regexp"
	"strconv"
)

type Signed struct{}

func (d Signed) GetName() string {
	return "signed"
}

func (d Signed) getInfo() string {
	return ""
}

func (d Signed) getSize() int {
	return -1 // Any size
}

func (d Signed) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0s(?P<int>-?[0-9]+)$"] = signedImportNoSize
	result["^0sd(?P<int>-?[0-9]+)$"] = signedImportNoSize

	return result
}

func (d Signed) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func signedImportNoSize(re *regexp.Regexp, input string) (*BMNumber, error) {
	// TODO and with size
	intDec := re.ReplaceAllString(input, "${int}")
	if s, err := strconv.ParseInt(intDec, 10, 0); err == nil {
		newNumber := BMNumber{}
		newNumber.number = make([]byte, 8)

		mask := int64(255)
		for i := 0; i < 8; i++ {
			newNumber.number[i] = byte(s & mask)
			s = s >> 8
		}

		newNumber.bits = 64
		newNumber.nType = Signed{}
		return &newNumber, nil
	} else {
		return nil, errors.New("invalid number" + input)
	}
}

func (d Signed) ExportString(n *BMNumber) (string, error) {
	return "", errors.New("not implemented")
}

func (d Signed) ShowInstructions() map[string]string {
	result := make(map[string]string)
	result["addop"] = "addsp"
	result["multop"] = "multsp"
	result["divop"] = "divsp"
	return result
}

func (d Signed) ShowPrefix() string {
	return "0s"
}
