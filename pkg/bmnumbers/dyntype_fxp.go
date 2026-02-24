package bmnumbers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"strconv"
)

type FXP struct {
	FixedPointName string
	s              int
	f              int
	instructions   map[string]string
}

func (d FXP) GetName() string {
	return d.FixedPointName
}

func (d FXP) getInfo() string {
	return ""
}

func (d FXP) GetSize() int {
	return d.s
}

func (d FXP) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0fxp<(?P<s>[0-9]+)\\.(?P<f>[0-9]+)>(?P<number>.+)$"] = fxpImport

	return result
}

func (d FXP) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func fxpImport(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")
	ss := re.ReplaceAllString(input, "${s}")
	fs := re.ReplaceAllString(input, "${f}")
	s, _ := strconv.Atoi(ss)
	f, _ := strconv.Atoi(fs)

	if s < 1 || s > 32 {
		return nil, errors.New("invalid s value for fixed point")
	}

	EventuallyCreateType("fxps"+ss+"f"+fs, nil)

	if numberNum, err := strconv.ParseFloat(number, 64); err != nil {
		return nil, errors.New("invalid number for fixed point")
	} else {
		scale := float64(int(1) << uint(f))
		fpNumber := int64(numberNum * scale)

		// Create a temporary buffer to write the number to
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, fpNumber)

		// Compute how many bytes we need to copy according to the number of bits (s)
		toCopy := (s-1)/8 + 1

		newNumber := BMNumber{}
		newNumber.nType = FXP{FixedPointName: "fxps" + ss + "f" + fs, s: s, f: f}
		newNumber.number = make([]byte, toCopy)
		copy(newNumber.number, buf.Bytes()[0:toCopy])
		newNumber.bits = s

		// Use a mask to clear the eventually unused bits on the last byte
		mask := byte(0xFF >> uint(8-(s-1)%8-1))
		newNumber.number[toCopy-1] = newNumber.number[toCopy-1] & mask

		return &newNumber, nil
	}
}

func (d FXP) ExportString(n *BMNumber) (string, error) {
	s := d.s
	f := d.f
	ss := strconv.Itoa(s)
	fs := strconv.Itoa(f)

	// Find out if the number is negative
	var isNegative bool
	if n.number[len(n.number)-1]&(uint8(128)>>uint(8-(s-1)%8-1)) != 0 {
		isNegative = true
	}

	copied := make([]byte, len(n.number))
	copy(copied, n.number)

	if isNegative {
		// Use a mask to clear the eventually unused bits on the last byte
		lastByte := (s-1)/8 + 1
		mask := byte(0xFF >> uint(8-(s-1)%8-1))
		copied[lastByte-1] = copied[lastByte-1] | ^mask

		for i := len(copied); i < 8; i++ {
			copied = append(copied, 0xFF)
		}
	} else {
		for i := len(copied); i < 8; i++ {
			copied = append(copied, 0x00)
		}
	}
	var number int64
	buf := bytes.NewReader(copied)

	if err := binary.Read(buf, binary.LittleEndian, &number); err != nil {
		return "", err
	}

	scale := float64(int(1) << uint(f))

	numberF := float64(number) / scale

	result := "0fxp<" + ss + "." + fs + ">" + strconv.FormatFloat(numberF, 'f', -1, 64)
	return result, nil
}

func (d FXP) ShowInstructions() map[string]string {
	return d.instructions
}

func (d FXP) ShowPrefix() string {
	return "0fxp<" + strconv.Itoa(d.s) + "." + strconv.Itoa(d.f) + ">"
}
