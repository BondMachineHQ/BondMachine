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
	Ranges        map[string][]float32
	Instances     map[string][]string
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
		Ranges:        make(map[string][]float32),
		Instances:     make(map[string][]string),
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
		}
		re = regexp.MustCompile(`^;fragtester\s+range\s+(?P<param>\w+)\s+arange\((?P<from>[-0-9.]+),(?P<to>[-0-9.]+),(?P<step>[-0-9.]+)\)$`)
		if re.MatchString(line) {
			param := re.ReplaceAllString(line, "${param}")
			from := re.ReplaceAllString(line, "${from}")
			to := re.ReplaceAllString(line, "${to}")
			step := re.ReplaceAllString(line, "${step}")
			fromF, _ := strconv.ParseFloat(from, 32)
			toF, _ := strconv.ParseFloat(to, 32)
			stepF, _ := strconv.ParseFloat(step, 32)
			ft.Ranges[param] = make([]float32, 0)
			for i := fromF; i < toF; i += stepF {
				ft.Ranges[param] = append(ft.Ranges[param], float32(i))
			}
			continue
		}
		re = regexp.MustCompile(`^;fragtester\s+instance\s+(?P<param>\w+)\s+(?P<seq>\S+)$`)
		if re.MatchString(line) {
			param := re.ReplaceAllString(line, "${param}")
			seq := re.ReplaceAllString(line, "${seq}")
			ft.Instances[param] = make([]string, 0)
			for _, v := range strings.Split(seq, ",") {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				ft.Instances[param] = append(ft.Instances[param], v)
			}
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

func (ft *FragTester) Sequences() int {
	seq := 1
	for _, rng := range ft.Ranges {
		seq *= len(rng)
	}
	for _, inst := range ft.Instances {
		seq *= len(inst)
	}
	return seq
}

func (ft *FragTester) DescribeFragment() {
	fmt.Printf("Name: %s\n", ft.Name)
	fmt.Printf("Sequences: %d\n", ft.Sequences())
	fmt.Printf("Ranges:\n")
	for param, rng := range ft.Ranges {
		fmt.Printf("  %s: ", param)
		for i, v := range rng {
			if i == len(rng)-1 {
				fmt.Printf("%f", v)
			} else {
				fmt.Printf("%f, ", v)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Instances:\n")
	for param, inst := range ft.Instances {
		fmt.Printf("  %s: ", param)
		for i, v := range inst {
			if i == len(inst)-1 {
				fmt.Printf("%s", v)
			} else {
				fmt.Printf("%s, ", v)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Inputs: %s\n", strings.Join(ft.Inputs, ", "))
	fmt.Printf("Outputs: %s\n", strings.Join(ft.Outputs, ", "))
	fmt.Printf("Sympy:\n", ft.Sympy, "\n")
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
