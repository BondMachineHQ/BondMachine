package fragtester

import (
	"bytes"
	"errors"
	"math"
	"text/template"
)

type templateData struct {
	Params  map[string]string
	Inputs  []string
	Outputs []string
	funcMap template.FuncMap
}

func (ft *FragTester) createBasicTemplateData() *templateData {
	result := new(templateData)

	result.Params = make(map[string]string)
	for k, v := range ft.Params {
		result.Params[k] = v
	}

	result.Inputs = make([]string, len(ft.Inputs))
	copy(result.Inputs, ft.Inputs)

	result.Outputs = make([]string, len(ft.Outputs))
	copy(result.Outputs, ft.Outputs)

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"n": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
		"ns": func(start, end, step int) []int {
			var result []int
			for i := start; i < end; i += step {
				result = append(result, i)
			}
			return result
		},
		"sum": func(a, b int) int {
			return a + b
		},
		"div": func(a, b int) int {
			return a / b
		},
		"mult": func(a, b int) int {
			return a * b
		},
		"pow": func(a, b int) int {
			return int(math.Pow(float64(a), float64(b)))
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
	result.funcMap = funcMap
	return result
}

func (ft *FragTester) ApplySympyTemplate() (string, error) {
	data := ft.Sympy
	templateData := ft.createBasicTemplateData()
	t, err := template.New("sympy").Funcs(templateData.funcMap).Parse(data)
	if err != nil {
		return "", err
	}
	var f bytes.Buffer
	err = t.Execute(&f, *templateData)
	if err != nil {
		return "", err
	}
	return f.String(), nil
}

func (ft *FragTester) ApplyAppTemplate(flavor string) (string, error) {
	var data string
	switch flavor {
	case "cpynqapi":
		data = CPynqApi
	default:
		return "", errors.New("unknown app template flavor: " + flavor)
	}
	templateData := ft.createBasicTemplateData()
	t, err := template.New("app").Funcs(templateData.funcMap).Parse(data)
	if err != nil {
		return "", err
	}
	var f bytes.Buffer
	err = t.Execute(&f, *templateData)
	if err != nil {
		return "", err
	}
	return f.String(), nil
}
