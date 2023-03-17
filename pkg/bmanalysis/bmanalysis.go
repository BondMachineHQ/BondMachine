package bmanalysis

import (
	"bytes"
	"text/template"
)

type BmAnalysis struct {
	ProjectsList    []string
	PivotRun		int
	funcMap template.FuncMap
}

func CreateAnalysisTemplate() *BmAnalysis {
	result := new(BmAnalysis)
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	result.funcMap = funcMap
	return result
}

func (s *BmAnalysis) WritePython() (string, error) {

	var f bytes.Buffer

	t, err := template.New("analysis").Funcs(s.funcMap).Parse(notebook)
	
	if err != nil {
		return "", err
	}

	err = t.Execute(&f, *s)
	if err != nil {
		return "", err
	}

	return f.String(), nil

}