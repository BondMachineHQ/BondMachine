package bmanalysis

import (
	"bytes"
	"errors"
	"text/template"
)

type BmAnalysis struct {
	ProjectsList   []string
	PivotRun       int
	funcMap        template.FuncMap
	BmAnalysisType string
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
	var t *template.Template
	var err error

	switch s.BmAnalysisType {
	case "ml":
		t, err = template.New("analysis").Funcs(s.funcMap).Parse(notebookML)
	case "mlsim":
		t, err = template.New("analysis").Funcs(s.funcMap).Parse(notebookMLSim)
	default:
		return "", errors.New("invalid analysis type")
	}

	if err != nil {
		return "", err
	}
	err = t.Execute(&f, *s)
	if err != nil {
		return "", err
	}

	return f.String(), nil

}
