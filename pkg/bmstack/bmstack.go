package bmstack

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"text/template"
)

type Push struct {
	Agent string
	Tick  uint64
	Value string
}

type Pop struct {
	Agent string
	Tick  uint64
}

type TestBenchData struct {
	Pops         []Pop
	Pushes       []Push
	TestSequence []string // Pushes and pops in order
}

type BmStack struct {
	ModuleName string
	DataSize   int
	Depth      int
	Senders    []string
	Receivers  []string
	MemType    string
	funcMap    template.FuncMap

	// TestBench data
	TestBenchData
}

func CreateBasicStack() *BmStack {
	result := new(BmStack)
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"next": func(i int, max int) int {
			if i < max-1 {
				return i + 1
			} else {
				return 0
			}
		},
		"bits": func(i int) int {
			return NeededBits(i)
		},
	}
	result.funcMap = funcMap

	// Default values
	result.ModuleName = "bmstack"

	return result
}

func (s *BmStack) WriteHDL() (string, error) {

	var f bytes.Buffer

	t, err := template.New("stack").Funcs(s.funcMap).Parse(stack)
	if err != nil {
		return "", err
	}

	err = t.Execute(&f, *s)
	if err != nil {
		return "", err
	}

	return f.String(), nil

}

func (s *BmStack) WriteTestBench() (string, error) {

	var f bytes.Buffer

	// Sort the pushes and pops by tick
	sort.Slice(s.Pushes, func(i, j int) bool {
		return s.Pushes[i].Tick < s.Pushes[j].Tick
	})
	sort.Slice(s.Pops, func(i, j int) bool {
		return s.Pops[i].Tick < s.Pops[j].Tick
	})

	s.TestSequence = make([]string, 0)

	var subSeq []string
	var absTick uint64
	var oldAbsTick uint64

	for i, j := 0, 0; i < len(s.Pushes) || j < len(s.Pops); {
		if i == len(s.Pushes) {
			absTick = s.Pops[j].Tick
			if absTick != oldAbsTick {
				relTick := absTick - oldAbsTick
				s.TestSequence = append(s.TestSequence, subSeq...)
				s.TestSequence = append(s.TestSequence, "")
				s.TestSequence = append(s.TestSequence, "#"+strconv.Itoa(int(relTick))+";")
				subSeq = make([]string, 0)
				oldAbsTick += relTick
			}
			agent := s.Pops[j].Agent
			tick := strconv.Itoa(int(absTick))
			subSeq = append(subSeq, "// Pop agent "+agent+" at tick "+tick)
			subSeq = append(subSeq, agent+"Impulse=1;")
			subSeq = append(subSeq, "#5;")
			subSeq = append(subSeq, agent+"Impulse=0;")
			j++
		} else if j == len(s.Pops) {
			absTick = s.Pushes[i].Tick
			if absTick != oldAbsTick {
				relTick := absTick - oldAbsTick
				s.TestSequence = append(s.TestSequence, subSeq...)
				s.TestSequence = append(s.TestSequence, "")
				s.TestSequence = append(s.TestSequence, "#"+strconv.Itoa(int(relTick))+";")
				subSeq = make([]string, 0)
				oldAbsTick += relTick
			}
			agent := s.Pushes[i].Agent
			value := s.Pushes[i].Value
			tick := strconv.Itoa(int(absTick))
			subSeq = append(subSeq, "// Push agent "+agent+" at tick "+tick+" with value "+value)
			subSeq = append(subSeq, agent+"Data="+value+";")
			subSeq = append(subSeq, agent+"Impulse=1;")
			subSeq = append(subSeq, "#5;")
			subSeq = append(subSeq, agent+"Impulse=0;")
			i++
		} else {
			if s.Pushes[i].Tick <= s.Pops[j].Tick {
				absTick = s.Pushes[i].Tick
				if absTick != oldAbsTick {
					relTick := absTick - oldAbsTick
					s.TestSequence = append(s.TestSequence, subSeq...)
					s.TestSequence = append(s.TestSequence, "")
					s.TestSequence = append(s.TestSequence, "#"+strconv.Itoa(int(relTick))+";")
					subSeq = make([]string, 0)
					oldAbsTick += relTick
				}
				agent := s.Pushes[i].Agent
				value := s.Pushes[i].Value
				tick := strconv.Itoa(int(absTick))
				subSeq = append(subSeq, "// Push agent "+agent+" at tick "+tick+" with value "+value)
				subSeq = append(subSeq, agent+"Data="+value+";")
				subSeq = append(subSeq, agent+"Impulse=1;")
				subSeq = append(subSeq, "#5;")
				subSeq = append(subSeq, agent+"Impulse=0;")
				i++
			} else {
				absTick = s.Pops[j].Tick
				if absTick != oldAbsTick {
					relTick := absTick - oldAbsTick
					s.TestSequence = append(s.TestSequence, subSeq...)
					s.TestSequence = append(s.TestSequence, "")
					s.TestSequence = append(s.TestSequence, "#"+strconv.Itoa(int(relTick))+";")
					subSeq = make([]string, 0)
					oldAbsTick += relTick
				}
				agent := s.Pops[j].Agent
				tick := strconv.Itoa(int(absTick))
				subSeq = append(subSeq, "// Pop agent "+agent+" at tick "+tick)
				subSeq = append(subSeq, agent+"Impulse=1;")
				subSeq = append(subSeq, "#5;")
				subSeq = append(subSeq, agent+"Impulse=0;")
				j++
			}
		}
	}
	s.TestSequence = append(s.TestSequence, subSeq...)

	t, err := template.New("testbench").Funcs(s.funcMap).Parse(testbench)
	if err != nil {
		return "", err
	}

	err = t.Execute(&f, *s)
	if err != nil {
		return "", err
	}

	return f.String(), nil

}

func NeededBits(num int) int {
	if num > 0 {
		for bits := 1; true; bits++ {
			if 1<<uint8(bits) >= num {
				return bits
			}
		}
	}
	return 0
}

func (td *TestBenchData) SaveJSON(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(td)
}

func (td *TestBenchData) LoadJSON(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(td)
}
