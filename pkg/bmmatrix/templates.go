package bmmatrix

import (
	"bytes"
	"errors"
	"math"
	"text/template"
)

type templateData1M struct {
	Mtx     [][]float32
	funcMap template.FuncMap
}

type templateData2M struct {
	Mtx1    [][]float32
	Mtx2    [][]float32
	funcMap template.FuncMap
}

func (exp *BasmExporter) createBasicTemplateData1M() *templateData1M {
	result := new(templateData1M)
	// result.Mtx = make([][]float32, result.NumGates)
	// for g, m := range sim.Mtx {
	// 	result.MtxReal[g] = make([][]float32, m.N)
	// 	result.MtxImag[g] = make([][]float32, m.N)
	// 	for i := 0; i < m.N; i++ {
	// 		result.MtxReal[g][i] = make([]float32, m.N)
	// 		result.MtxImag[g][i] = make([]float32, m.N)
	// 		for j := 0; j < m.N; j++ {
	// 			result.MtxReal[g][i][j] = m.Data[i][j].Real
	// 			result.MtxImag[g][i][j] = m.Data[i][j].Imag
	// 		}
	// 	}
	// }

	result.funcMap = getFuncMap()
	return result
}

func (exp *BasmExporter) createBasicTemplateData2M() *templateData2M {
	result := new(templateData2M)
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

func (exp *BasmExporter) ApplyTemplate1M() (string, error) {
	templateData := exp.createBasicTemplateData1M()
	t, err := template.New("mult").Funcs(templateData.funcMap).Parse(templateMult)
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

func (exp *BasmExporter) ApplyTemplate2M() (string, error) {
	templateData := exp.createBasicTemplateData2M()
	t, err := template.New("mult").Funcs(templateData.funcMap).Parse(templateMult)
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
