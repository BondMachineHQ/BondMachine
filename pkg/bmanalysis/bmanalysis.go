package bmanalysis

import (
	"bytes"
	"text/template"
)

type BmAnalysis struct {
	ProjectLists    []string
}

func CreateAnalysisTemplate() *BmAnalysis {
	result := new(BmAnalysis)
	return result
}

func (s *BmAnalysis) WritePython() (string, error) {

	var f bytes.Buffer

	t, err := template.New("analysis").Parse(notebook)
	
	if err != nil {
		return "", err
	}

	err = t.Execute(&f, *s)
	if err != nil {
		return "", err
	}

	return f.String(), nil

}