package basm

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"golang.org/x/exp/maps"
)

type templateData struct {
	Params  map[string]string
	funcMap template.FuncMap
}

func createBasicTemplateData() *templateData {
	result := new(templateData)
	result.Params = make(map[string]string)
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"incs": func(i string) string {
			j, _ := strconv.Atoi(i)
			return strconv.Itoa(j + 1)
		},
		"dec": func(i int) int {
			return i - 1
		},
		"decs": func(i string) string {
			j, _ := strconv.Atoi(i)
			return strconv.Itoa(j - 1)
		},
		"add": func(i, j int) int {
			return i + j
		},
		"adds": func(i, j string) string {
			k, _ := strconv.Atoi(i)
			l, _ := strconv.Atoi(j)
			return strconv.Itoa(k + l)
		},
		"intRange": func(startS, endS string) []int {
			start, _ := strconv.Atoi(startS)
			end, _ := strconv.Atoi(endS)

			n := end - start
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = start + i
			}
			return result
		},
		"atoi": func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		},
	}

	result.funcMap = funcMap
	return result
}

func (bi *BasmInstance) templateAutoMark() error {
	// Set the template meta to true for all the sections and fragments that contains the template mark
	for _, section := range bi.sections {
		if section.isTemplate() {
			section.sectionBody.BasmMeta = section.sectionBody.SetMeta("template", "true")
		} else {
			section.sectionBody.BasmMeta = section.sectionBody.SetMeta("template", "false")
		}
	}

	for _, fragment := range bi.fragments {
		if fragment.isTemplate() {
			fragment.fragmentBody.BasmMeta = fragment.fragmentBody.SetMeta("template", "true")
		} else {
			fragment.fragmentBody.BasmMeta = fragment.fragmentBody.SetMeta("template", "false")
		}
	}
	return nil
}

func isTemplate(s string) bool {
	pattern := `{{[^{}]*}}`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllString(s, -1)
	return len(matches) > 0
}

func (f *BasmFragment) isTemplate() bool {
	return isTemplate(f.fragmentBody.Flat())
}

func (s *BasmSection) isTemplate() bool {
	return isTemplate(s.sectionBody.Flat())
}

func applyTemplate(e *bmline.BasmElement, params map[string]string) error {
	// Search for default values
	for k, v := range params {
		if strings.HasPrefix(k, "default_") {
			key := strings.TrimPrefix(k, "default_")
			if _, ok := params[key]; !ok {
				params[key] = v
			}
		}
	}

	td := createBasicTemplateData()

	for key, value := range params {
		td.Params[key] = value
	}

	var f bytes.Buffer

	t, err := template.New("template").Funcs(td.funcMap).Parse(e.GetValue())
	if err != nil {
		return err
	}

	err = t.Execute(&f, *td)
	if err != nil {
		return err
	}

	e.SetValue(f.String())

	// fmt.Println(params)

	return nil
}

func bodyTemplateResolver(body *bmline.BasmBody, params map[string]string) error {

	// TODO The parsing with {{ end }} construct is not working. It should be fixed

	bodyParams := make(map[string]string)
	maps.Copy(bodyParams, params)
	maps.Copy(bodyParams, body.LoopMeta())

	for _, line := range body.Lines {

		if isTemplate(line.Flat()) {
			lineParams := make(map[string]string)
			maps.Copy(lineParams, bodyParams)
			maps.Copy(lineParams, line.LoopMeta())

			operation := line.Operation
			operationParams := make(map[string]string)
			maps.Copy(operationParams, lineParams)
			maps.Copy(operationParams, operation.LoopMeta())
			if isTemplate(operation.Flat()) {
				if err := applyTemplate(operation, operationParams); err != nil {
					return err
				}
			}

			for _, arg := range line.Elements {
				argParams := make(map[string]string)
				maps.Copy(argParams, operationParams)
				maps.Copy(argParams, arg.LoopMeta())
				if isTemplate(arg.Flat()) {
					if err := applyTemplate(arg, argParams); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
