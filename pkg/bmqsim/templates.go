package bmqsim

import (
	"bytes"
	"text/template"
)

// Modes:
// seq_hardcoded
// seq
// full_hw_hardcoded

var Templates = map[string]string{
	"seq_hardcoded_real": SeqHardcodedReal,
}

type templateData struct {
	funcMap template.FuncMap
}

func (sim *BmQSimulator) createBasicTemplateData() *templateData {
	result := new(templateData)
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
	}
	result.funcMap = funcMap
	return result
}

func (sim *BmQSimulator) ApplyTemplate(mode string) (string, error) {
	d := Templates[mode]
	templateData := sim.createBasicTemplateData()
	t, err := template.New(mode).Funcs(templateData.funcMap).Parse(d)
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
