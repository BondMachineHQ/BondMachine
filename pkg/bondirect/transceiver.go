package bondirect

import (
	"bytes"
	"text/template"
)

func (be *BondirectElement) GenerateTransceiver(prefix, nodeName, edgeName, direction string) (string, error) {
	// Implementation for generating a Transceiver

	trn := bondRx
	if direction == "out" {
		trn = bondTx
	}

	var f bytes.Buffer

	t, err := template.New("transceiver").Funcs(funcMap).Parse(trn)
	if err != nil {
		return "", err
	}

	err = t.Execute(&f, be.TData)
	if err != nil {
		return "", err
	}

	return f.String(), nil
}
