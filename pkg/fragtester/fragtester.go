package fragtester

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FragTester struct {
	Valid         bool
	RegisterSize  int
	DataType      string
	TypePrefix    string
	Params        map[string]string
	OpString      string
	Debug         bool
	Verbose       bool
	Name          string
	Inputs        []string
	Outputs       []string
	Sympy         string
	NeuronLibPath string
}

func NewFragTester() *FragTester {
	return &FragTester{
		Valid:         false,
		DataType:      "float32",
		TypePrefix:    "0f",
		Params:        make(map[string]string),
		OpString:      "",
		Debug:         false,
		Verbose:       false,
		Name:          "",
		Inputs:        make([]string, 0),
		Outputs:       make([]string, 0),
		Sympy:         "",
		NeuronLibPath: "",
	}
}

func (ft *FragTester) AnalyzeFragment(fragment string) error {
	// Read the string line by line
	lines := strings.Split(fragment, "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`^;fragtester.*$`)
		if re.MatchString(line) {
			ft.Valid = true
			continue
		}
		re = regexp.MustCompile(`^%fragment\s+(?P<name>\w+).+(?P<resin>resin:[\w:]+).+(?P<resout>resout:[\w:]+)$`)
		if re.MatchString(line) {
			name := re.ReplaceAllString(line, "${name}")
			resin := re.ReplaceAllString(line, "${resin}")
			resout := re.ReplaceAllString(line, "${resout}")
			ft.Name = name
			ft.Inputs = strings.Split(resin, ":")[1:]
			ft.Outputs = strings.Split(resout, ":")[1:]
			continue
		}
		re = regexp.MustCompile(`^;sympy\s(?P<sympy>.+)$`)
		if re.MatchString(line) {
			sympy := re.ReplaceAllString(line, "${sympy}")
			ft.Sympy += sympy + "\n"
		}
	}

	return nil
}

func (ft *FragTester) WriteBasm() (string, error) {
	result := fmt.Sprintf("%%meta bmdef     global registersize:%d\n", ft.RegisterSize)
	result += fmt.Sprintf("%%meta bmdef     global iomode: sync\n")
	result += fmt.Sprintf("%%meta fidef     cpu fragment: %s %s\n", ft.Name, ft.OpString)
	result += fmt.Sprintf("%%meta cpdef     cpu fragcollapse:cpu\n")
	for i, input := range ft.Inputs {
		result += fmt.Sprintf("%%meta filinkatt linki" + input + " fi:ext, type: input, index: " + strconv.Itoa(i) + "\n")
		result += fmt.Sprintf("%%meta filinkatt linki" + input + " fi:cpu, type: input, index: " + strconv.Itoa(i) + "\n")
	}
	for i, output := range ft.Outputs {
		result += fmt.Sprintf("%%meta filinkatt linko" + output + " fi:ext, type: output, index: " + strconv.Itoa(i) + "\n")
		result += fmt.Sprintf("%%meta filinkatt linko" + output + " fi:cpu, type: output, index: " + strconv.Itoa(i) + "\n")
	}
	return result, nil
}

func (ft *FragTester) WriteSympy() (string, error) {
	result := ""
	result += ft.Sympy
	return result, nil
}
