package bmmatrix

import (
	"bytes"
	"errors"
	"math"
	"text/template"
)

type templateData1M struct {
	Mtx     [][]float32
	Iomode  string
	funcMap template.FuncMap
}

type templateData2M struct {
	Mtx1    [][]string
	Mtx2    [][]string
	Iomode  string
	funcMap template.FuncMap
}

func (exp *BasmExporter) createBasicTemplateData1M() *templateData1M {
	result := new(templateData1M)
	result.funcMap = getFuncMap()
	return result
}

func (exp *BasmExporter) createBasicTemplateData2M() *templateData2M {
	result := new(templateData2M)
	result.funcMap = getFuncMap()
	return result
}

func getFuncMap() template.FuncMap {

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
		"rows": func(a [][]string) int {
			return len(a)
		},
		"cols": func(a [][]string) int {
			return len(a[0])
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
	return funcMap
}

func (exp *BasmExporter) ApplyTemplate1M(templateData *templateData1M, templateName string, templateObj string) (string, error) {
	t, err := template.New(templateName).Funcs(templateData.funcMap).Parse(templateObj)
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

func (exp *BasmExporter) ApplyTemplate2M(templateData *templateData2M, templateName string, templateObj string) (string, error) {
	t, err := template.New(templateName).Funcs(templateData.funcMap).Parse(templateObj)
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
