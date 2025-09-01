package bondirect

import (
	"bytes"
	"text/template"
)

func (be *BondirectElement) GenerateLine(prefix, nodeName, edgeName string) (string, error) {
	// Implementation for generating a Line

	// Fill Template Data with the request values
	be.TData.NodeName = nodeName
	be.TData.EdgeName = edgeName

	// Define the line template
	ln := bdLine
	var f bytes.Buffer

	t, err := template.New("line").Funcs(funcMap).Parse(ln)
	if err != nil {
		return "", err
	}

	// Execute the template with the filled data
	err = t.Execute(&f, be.TData)
	if err != nil {
		return "", err
	}

	return f.String(), nil
}
