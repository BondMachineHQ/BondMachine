package bmnumbers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"regexp"
	"strconv"
)

type FloPoCoFixedPoint struct {
	FixedPointName string
	s              int
	f              int
	instructions   map[string]string
}

func (d FloPoCoFixedPoint) GetName() string {
	return d.FixedPointName
}

func (d FloPoCoFixedPoint) getInfo() string {
	return ""
}

func (d FloPoCoFixedPoint) GetSize() int {
	return d.s
}

func (d FloPoCoFixedPoint) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0flpfps<(?P<s>[0-9]+)\\.(?P<f>[0-9]+)>(?P<number>.+)$"] = floPoCofixedPointImport

	return result
}

func (d FloPoCoFixedPoint) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

// fixedPointImport converts a string like "12.75fps16f6" (captured by re)
// into an internal BMNumber in two’s-complement fixed-point.
// func floPoCofixedPointImport(re *regexp.Regexp, input string) (*BMNumber, error) {
// 	// 1. extract the named submatches using ReplaceAllString
// 	ss := re.ReplaceAllString(input, "${s}")
// 	fs := re.ReplaceAllString(input, "${f}")
// 	numStr := re.ReplaceAllString(input, "${number}")
// 	if ss == "" || fs == "" || numStr == "" {
// 		return nil, fmt.Errorf("input %q does not match fixed-point pattern", input)
// 	}

// 	// 2. parse s (total bits) and f (fractional bits)
// 	s, err := strconv.Atoi(ss)
// 	if err != nil || s < 2 || s > 64 {
// 		return nil, fmt.Errorf("invalid s (%s): need 2–64", ss)
// 	}
// 	f, err := strconv.Atoi(fs)
// 	if err != nil || f < 0 || f >= s {
// 		return nil, fmt.Errorf("invalid f (%s): need 0≤f<s", fs)
// 	}

// 	// 3. parse the numeric literal and quantise
// 	val, err := strconv.ParseFloat(numStr, 64)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid number %q", numStr)
// 	}
// 	scale := math.Pow(2, float64(f))
// 	fp := int64(math.Round(val * scale))

// 	// 4. range-check against signed s-bit
// 	min := -(int64(1) << (s - 1))
// 	max := (int64(1) << (s - 1)) - 1
// 	if fp < min || fp > max {
// 		return nil, fmt.Errorf("%.6g does not fit in Q%d.%d", val, s, f)
// 	}

// 	// 5. serialise to a little-endian byte slice
// 	byteLen := (s + 7) / 8
// 	raw := make([]byte, byteLen)
// 	for i := 0; i < byteLen; i++ {
// 		raw[i] = byte(fp >> (8 * i))
// 	}
// 	// clear the unused high-order bits in the last byte
// 	unused := byteLen*8 - s
// 	if unused > 0 {
// 		mask := byte(0xFF >> unused)
// 		raw[byteLen-1] &= mask
// 	}

// 	// 6. build the fixed-point instruction names (if you need them)
// 	instr := map[string]string{
// 		"multop": "multflpfps" + ss + "f" + fs,
// 		"addop":  "addflpfps" + ss + "f" + fs,
// 		"divop":  "divflpfps" + ss + "f" + fs,
// 	}

// 	// 7. register the new type
// 	EventuallyCreateType("flpfps"+ss+"f"+fs, nil)

// 	log.Printf("number is %d, byteLen is %d, raw is %v", fp, byteLen, raw)

// 	// 8. wrap up in a BMNumber and return
// 	return &BMNumber{
// 		number: raw,
// 		bits:   s,
// 		nType: FloPoCoFixedPoint{
// 			FixedPointName: "flpfps" + ss + "f" + fs,
// 			s:              s,
// 			f:              f,
// 			instructions:   instr,
// 		},
// 	}, nil
// }

func floPoCofixedPointImport(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")
	ss := re.ReplaceAllString(input, "${s}")
	fs := re.ReplaceAllString(input, "${f}")
	s, _ := strconv.Atoi(ss)
	f, _ := strconv.Atoi(fs)

	if s < 1 || s > 32 {
		return nil, errors.New("invalid s value for fixed point")
	}

	EventuallyCreateType("fps"+ss+"f"+fs, nil)

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
		newNumber.nType = FixedPoint{FixedPointName: "fps" + ss + "f" + fs, s: s, f: f}
		newNumber.number = make([]byte, toCopy)
		copy(newNumber.number, buf.Bytes()[0:toCopy])
		newNumber.bits = s

		// Use a mask to clear the eventually unused bits on the last byte
		mask := byte(0xFF >> uint(8-(s-1)%8-1))
		newNumber.number[toCopy-1] = newNumber.number[toCopy-1] & mask

		return &newNumber, nil
	}
}

func (d FloPoCoFixedPoint) ExportString(n *BMNumber) (string, error) {
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

	result := "0flpfps<" + ss + "." + fs + ">" + strconv.FormatFloat(numberF, 'f', -1, 64)
	return result, nil
}

func (d FloPoCoFixedPoint) ShowInstructions() map[string]string {
	return d.instructions
}

func (d FloPoCoFixedPoint) ShowPrefix() string {
	return "0flpfps<" + strconv.Itoa(d.s) + "." + strconv.Itoa(d.f) + ">"
}
