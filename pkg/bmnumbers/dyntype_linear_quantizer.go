package bmnumbers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"strconv"
)

type LinearQuantizer struct {
	linearQuantizerName string
	s                   int
	t                   int
}

func (d LinearQuantizer) GetName() string {
	return d.linearQuantizerName
}

func (d LinearQuantizer) getInfo() string {
	return ""
}

func (d LinearQuantizer) getSize() int {
	return d.s
}

func (d LinearQuantizer) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0lq<(?P<s>[0-9]+)\\.(?P<t>[0-9]+)>(?P<number>.+)$"] = linearQuantizerImport

	return result
}

func (d LinearQuantizer) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func linearQuantizerImport(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")
	ss := re.ReplaceAllString(input, "${s}")
	ts := re.ReplaceAllString(input, "${t}")
	s, _ := strconv.Atoi(ss)
	t, _ := strconv.Atoi(ts)

	if s < 1 || s > 32 {
		return nil, errors.New("invalid s value for linear quantizer")
	}

	// Get the linear quantizer ranges struct
	var lqRanges *map[int]LinearDataRange
	for _, t := range AllDynamicalTypes {
		if t.GetName() == "dyn_linear_quantizer" {
			lqRanges = t.(DynLinearQuantizer).Ranges
		}
	}

	if lqRanges != nil {
		if _, ok := (*lqRanges)[t]; !ok {
			return nil, errors.New("invalid type value for linear quantizer")
		}
	} else {
		return nil, errors.New("linear quantizer ranges not found")
	}

	EventuallyCreateType("lqs"+ss+"t"+ts, nil)

	bandNum := float64(int(1) << uint(s-1))
	bandSize := ((*lqRanges)[t].Max) / bandNum

	if numberNum, err := strconv.ParseFloat(number, 64); err != nil {
		return nil, errors.New("invalid number for linear quantizer")
	} else {
		band := int64(numberNum / bandSize)
		if band >= int64(bandNum) || band <= -int64(bandNum) {
			return nil, errors.New("number out of range for linear quantizer")
		}

		// Create a temporary buffer to write the number to
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, band)

		// Compute how many bytes we need to copy according to the number of bits (s)
		toCopy := (s-1)/8 + 1

		newNumber := BMNumber{}
		newNumber.nType = LinearQuantizer{linearQuantizerName: "lqs" + ss + "t" + ts, s: s, t: t}
		newNumber.number = make([]byte, toCopy)
		copy(newNumber.number, buf.Bytes()[0:toCopy])
		newNumber.bits = s

		// Use a mask to clear the eventually unused bits on the last byte
		mask := byte(0xFF >> uint(8-(s-1)%8-1))
		newNumber.number[toCopy-1] = newNumber.number[toCopy-1] & mask

		return &newNumber, nil
	}
}

func (d LinearQuantizer) ExportString(n *BMNumber) (string, error) {
	s := d.s
	t := d.t
	ss := strconv.Itoa(s)
	ts := strconv.Itoa(t)

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

	// Get the linear quantizer ranges struct
	var lqRanges *map[int]LinearDataRange
	for _, t := range AllDynamicalTypes {
		if t.GetName() == "dyn_linear_quantizer" {
			lqRanges = t.(DynLinearQuantizer).Ranges
		}
	}

	bandNum := float64(int(1) << uint(s-1))
	bandSize := (*lqRanges)[t].Max / bandNum

	numberF := float64(number) * bandSize

	result := "0lq<" + ss + "." + ts + ">" + strconv.FormatFloat(numberF, 'f', -1, 64)
	return result, nil
}

func (d LinearQuantizer) ShowInstructions() map[string]string {
	result := make(map[string]string)
	result["addop"] = "addp"
	result["multop"] = "multp"
	result["divop"] = "divp"
	return result
}

func (d LinearQuantizer) ShowPrefix() string {
	return "0lq<" + strconv.Itoa(d.s) + "." + strconv.Itoa(d.t) + ">"
}
