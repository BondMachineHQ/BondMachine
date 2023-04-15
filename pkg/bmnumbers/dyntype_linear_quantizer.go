package bmnumbers

import (
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

	var binNum string

	bandNum := float64(int(1) << uint(s))
	bandSize := ((*lqRanges)[t].Max - (*lqRanges)[t].Min) / bandNum

	if numberNum, err := strconv.ParseFloat(number, 64); err != nil {
		return nil, errors.New("invalid number for linear quantizer")
	} else {
		band := int((numberNum - (*lqRanges)[t].Min) / bandSize)
		if band < 0 {
			band = 0
		} else if band >= int(bandNum) {
			band = int(bandNum) - 1
		}
		binNum = "00000000000000000000000000000000" + strconv.FormatInt(int64(band), 2)
		binNum = binNum[len(binNum)-s:]
	}

	var nextByte string
	newNumber := BMNumber{}
	newNumber.number = make([]byte, (len(binNum)-1)/8+1)
	newNumber.bits = len(binNum)
	newNumber.nType = LinearQuantizer{linearQuantizerName: "lqs" + ss + "t" + ts, s: s, t: t}

	for i := 0; ; i++ {

		if len(binNum) > 8 {
			nextByte = binNum[len(binNum)-8:]
			binNum = binNum[:len(binNum)-8]
		} else if len(binNum) > 0 {
			nextByte = binNum
			binNum = ""
		} else {
			break
		}

		if val, err := strconv.ParseUint(nextByte, 2, 8); err == nil {
			decoded := byte(val)
			newNumber.number[i] = decoded
		} else {
			return nil, errors.New("invalid binary number" + input)
		}
	}
	return &newNumber, nil
}

func (d LinearQuantizer) ExportString(n *BMNumber) (string, error) {
	s := d.s
	t := d.t
	ss := strconv.Itoa(s)
	ts := strconv.Itoa(t)
	number, _ := n.ExportUint64()

	// Get the linear quantizer ranges struct
	var lqRanges *map[int]LinearDataRange
	for _, t := range AllDynamicalTypes {
		if t.GetName() == "dyn_linear_quantizer" {
			lqRanges = t.(DynLinearQuantizer).Ranges
		}
	}

	bandNum := float64(int(1) << uint(s))
	bandSize := ((*lqRanges)[t].Max - (*lqRanges)[t].Min) / bandNum

	numberF := float64(number)*bandSize + (*lqRanges)[t].Min

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
