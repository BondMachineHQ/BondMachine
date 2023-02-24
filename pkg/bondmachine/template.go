package bondmachine

import (
	"strconv"
	"text/template"
)

type stateMachine struct {
	Nums     int
	Bits     int
	Buswidth string
	Names    []string
	Binary   []string
}

type templateData struct {
	Rsize       int
	Buswidth    string
	Inputs      []string
	InputNum    int
	InputsBins  []string
	Outputs     []string
	OutputNum   int
	OutputsBins []string
	SendSM      stateMachine
	funcmap     template.FuncMap
	PackageName string
	ModuleName  string
}

func (bmach *Bondmachine) createBasicTemplateData() *templateData {
	result := new(templateData)
	result.Rsize = int(bmach.Rsize)
	result.Buswidth = "[" + strconv.Itoa(int(bmach.Rsize)-1) + ":0]"
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
		"bits": func(i int) int {
			return Needed_bits(i)
		},
	}
	result.funcmap = funcMap
	return result
}
