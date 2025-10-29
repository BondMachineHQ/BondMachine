package fragtester

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

type FragTester struct {
	Valid         bool
	RegisterSize  int
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Ranges        map[string][]float32
	Instances     map[string][]string
	Vars          []string
	OpString      string
	Debug         bool
	Verbose       bool
	Name          string
	NameSuffix    string
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
		Vars:          make([]string, 0),
		OpString:      "",
		Debug:         false,
		Verbose:       false,
		Name:          "",
		NameSuffix:    "",
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
			ft.Vars = append(ft.Vars, param)
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
			ft.Vars = append(ft.Vars, param)
			continue
		}
		re = regexp.MustCompile(`^%fragment\s+(?P<name>\w+).+(?P<resin>resin:[\w:]+).+(?P<resout>resout:[\w:]+).*$`)
		if re.MatchString(line) {
			name := re.ReplaceAllString(line, "${name}")
			resin := re.ReplaceAllString(line, "${resin}")
			resout := re.ReplaceAllString(line, "${resout}")
			ft.Name = name
			ft.Inputs = strings.Split(resin, ":")[1:]
			ft.Outputs = strings.Split(resout, ":")[1:]
			continue
		}
		// Support both resin:... resout:... and resout:... resin:...
		re = regexp.MustCompile(`^%fragment\s+(?P<name>\w+).+(?P<resout>resout:[\w:]+).+(?P<resin>resin:[\w:]+).*$`)
		if re.MatchString(line) {
			name := re.ReplaceAllString(line, "${name}")
			resin := re.ReplaceAllString(line, "${resin}")
			resout := re.ReplaceAllString(line, "${resout}")
			ft.Name = name
			ft.Inputs = strings.Split(resin, ":")[1:]
			ft.Outputs = strings.Split(resout, ":")[1:]
			continue
		}
		// Support only resin:...
		re = regexp.MustCompile(`^%fragment\s+(?P<name>\w+).+(?P<resin>resin:[\w:]+).*$`)
		if re.MatchString(line) {
			name := re.ReplaceAllString(line, "${name}")
			resin := re.ReplaceAllString(line, "${resin}")
			ft.Name = name
			ft.Inputs = strings.Split(resin, ":")[1:]
			continue
		}
		// Support only resout:...
		re = regexp.MustCompile(`^%fragment\s+(?P<name>\w+).+(?P<resout>resout:[\w:]+).*$`)
		if re.MatchString(line) {
			name := re.ReplaceAllString(line, "${name}")
			resout := re.ReplaceAllString(line, "${resout}")
			ft.Name = name
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

func (ft *FragTester) ApplySequence(seq int) {
	rank := len(ft.Ranges) + len(ft.Instances)
	shape := make([]int, rank)
	for i, v := range ft.Vars {
		if _, ok := ft.Ranges[v]; ok {
			shape[i] = len(ft.Ranges[v])
		}
		if _, ok := ft.Instances[v]; ok {
			shape[i] = len(ft.Instances[v])
		}
	}
	// fmt.Println("Shape:", shape)
	pos := make([]int, rank)
	for i := range rank {
		pos[i] = 0
	}
	for i := 0; i < seq; i++ {
		// fmt.Println("Pos:", pos)
		ft.nextPos(&pos, &shape)
	}

	for i, v := range ft.Vars {
		if _, ok := ft.Ranges[v]; ok {
			ft.Params[v] = fmt.Sprintf("%f", ft.Ranges[v][pos[i]])
			ft.OpString += fmt.Sprintf(", %s:%s", v, ft.Params[v])
		}
		if _, ok := ft.Instances[v]; ok {
			ft.Params[v] = ft.Instances[v][pos[i]]
			ft.OpString += fmt.Sprintf(", %s:%s", v, ft.Params[v])
			ft.NameSuffix += fmt.Sprintf("-%s-%s", v, ft.Params[v])
		}
	}
	// fmt.Println("Params:", ft.Params)
}

func (ft *FragTester) nextPos(pos *[]int, shape *[]int) {
	var next int
	for next = 0; next < len(*pos); next++ {
		if (*pos)[next] < (*shape)[next]-1 {
			break
		}
	}
	if next == len(*pos) {
		return
	}
	(*pos)[next]++
	for i := 0; i < next; i++ {
		(*pos)[i] = 0
	}
}

func (ft *FragTester) DescribeFragment() {
	keyColor := yellow
	fmt.Printf(keyColor("Name")+": %s\n", ft.Name)
	fmt.Printf(keyColor("Register size")+": %d\n", ft.RegisterSize)
	fmt.Printf(keyColor("Data type")+": %s\n", ft.DataType)
	fmt.Printf(keyColor("Type prefix")+": %s\n", ft.TypePrefix)
	fmt.Printf(keyColor("Params") + ":\n")
	for param, value := range ft.Params {
		fmt.Printf("  %s: ", param)
		if strings.Contains(value, ",") {
			for i, v := range strings.Split(value, ",") {
				v = strings.TrimSpace(v)
				if i == len(strings.Split(value, ","))-1 {
					fmt.Printf("%s", v)
				} else {
					fmt.Printf("%s, ", v)
				}
			}
		} else {
			fmt.Printf("%s", value)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Sequences: %d\n", ft.Sequences())
	fmt.Println("Vars:", strings.Join(ft.Vars, ", "))
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
	fmt.Print("Sympy:\n")
	for _, line := range strings.Split(ft.Sympy, "\n") {
		if line != "" {
			fmt.Printf("  %s\n", line)
		}
	}
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
	if result, err := ft.ApplySympyTemplate(); err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func (ft *FragTester) WriteApp(flavor string) (string, error) {
	if result, err := ft.ApplyAppTemplate(flavor); err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func (ft *FragTester) WriteStatistics() (string, error) {
	result := "{\"" + ft.Name + ft.NameSuffix + "\": 1}\n"
	return result, nil
}

func (ft *FragTester) WriteSicv2Endpoints() (string, error) {
	if len(ft.Inputs) == 0 || len(ft.Outputs) == 0 {
		return "", fmt.Errorf("cannot create SICv2 endpoints without inputs and outputs")
	}

	return fmt.Sprintf("i0,p0o%d", len(ft.Outputs)-1), nil
}

func (ft *FragTester) CreateMappingFile(filename string) error {
	ioMap := new(bondmachine.IOmap)
	ioMap.Assoc = make(map[string]string)

	for i := range ft.Inputs {
		ioMap.Assoc["i"+strconv.Itoa(i)] = strconv.Itoa(i)
	}
	for i := range ft.Outputs {
		ioMap.Assoc["o"+strconv.Itoa(i)] = strconv.Itoa(i)
	}

	// Write the file
	mapBytes, err := json.Marshal(ioMap)
	if err != nil {
		return fmt.Errorf("failed to marshal I/O map: %v", err)
	}
	if err := os.WriteFile(filename, mapBytes, 0644); err != nil {
		return fmt.Errorf("failed to write I/O map file: %v", err)
	}

	return nil
}
