package bondirect

import (
	"bytes"
	"text/template"
)

func (be *BondirectElement) GenerateEndpoint(prefix, nodeName string) (string, error) {
	// TODO Implementation for generating an Endpoint

	// Fill Template Data with the request values
	be.TData.Prefix = prefix
	be.TData.NodeName = nodeName
	if err := be.PopulateIOData(nodeName); err != nil {
		return "", err
	}
	if err := be.PopulateWireData(nodeName); err != nil {
		return "", err
	}

	// fmt.Println(be.DumpTemplateData())

	// Define the endpoint template
	en := bdEndpoint
	var f bytes.Buffer

	t, err := template.New("endpoint").Funcs(funcMap).Parse(en)
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
