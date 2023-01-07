package bmstack

import (
	"bytes"
	"text/template"
)

type stateMachine struct {
	Nums     int
	Bits     int
	Buswidth string
	Names    []string
	Binary   []string
}

type BmStack struct {
	ModuleName string
	DataSize   int
	Senders    []string
	Receivers  []string

	Rsize       uint8
	Buswidth    string
	Inputs      []string
	InputsBins  []string
	Outputs     []string
	OutputsBins []string
	SendSM      stateMachine
	funcMap     template.FuncMap
	PackageName string
}

func CreateBasicStack() *BmStack {
	result := new(BmStack)
	//result.Rsize = bmach.Rsize
	//result.Buswidth = "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0]"
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"next": func(i int) int {
			if i < result.SendSM.Nums-1 {
				return i + 1
			} else {
				return 0
			}
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
