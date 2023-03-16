package bmanalysis

import (
	"bytes"
	"text/template"
)

type BmStack struct {
	ModuleName string
	DataSize   int
	Depth      int
	Senders    []string
	Receivers  []string
	MemType    string
	funcMap    template.FuncMap
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

	t, err := template.New("stack").Funcs(s.funcMap).Parse(notebook)
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
