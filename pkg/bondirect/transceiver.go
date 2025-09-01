package bondirect

import (
	"bytes"
	"text/template"
)

func (be *BondirectElement) GenerateTransceiver(prefix, nodeName, edgeName, direction string) (string, error) {
	// Implementation for generating a Transceiver

	// Fill Template Data with the request values
	be.TData.NodeName = nodeName
	be.TData.EdgeName = edgeName

	// Define the transceiver template
	trn := bondRx
	if direction == "out" {
		trn = bondTx
	}

	var f bytes.Buffer

	t, err := template.New("transceiver").Funcs(funcMap).Parse(trn)
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
