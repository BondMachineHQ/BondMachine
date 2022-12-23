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
	Rsize       uint8
	Buswidth    string
	Inputs      []string
	InputsBins  []string
	Outputs     []string
	OutputsBins []string
	SendSM      stateMachine
	funcmap     template.FuncMap
	PackageName string
	ModuleName  string
}

func (bmach *Bondmachine) createBasicTemplateData() *templateData {
	result := new(templateData)
	result.Rsize = bmach.Rsize
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
	}
	result.funcmap = funcMap
	return result
}
