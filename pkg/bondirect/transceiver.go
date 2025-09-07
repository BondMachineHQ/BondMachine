package bondirect

import (
	"bytes"
	"fmt"
	"text/template"
)

func (be *BondirectElement) GenerateTransceiver(prefix, nodeName, edgeName, direction string) (string, error) {
	// Implementation for generating a Transceiver

	// Fill Template Data with the request values
	be.TData.Prefix = prefix
	be.TData.NodeName = nodeName
	be.TData.EdgeName = edgeName

	// Define the transceiver template
	trn := bondRx
	if direction == "out" {
		trn = bondTx
	}

	be.PopulateEdgeParams(edgeName)

	fmt.Println(be.DumpTemplateData())

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

func (be *BondirectElement) GetTransceiverSignals(trName string) ([]string, error) {
	signals := make([]string, 0)

	sigs := 0
	for _, tr := range be.Mesh.Transceivers {
		if tr.Name == trName {
			for sName, s := range tr.Signals {
				if s.Type == "clock" {
					// Clock goes first
					pref := make([]string, 0)
					pref = append(pref, sName)
					signals = append(pref, signals...)
					sigs++
					continue
				}
				if s.Type == "data" {
					signals = append(signals, sName)
					sigs++
					continue
				} else {
					return nil, fmt.Errorf("unknown signal type: %s", s.Type)
				}
			}
			break
		}
	}

	if sigs >= 2 {
		return signals, nil
	}
	return nil, fmt.Errorf("not enough signals found for transceiver: %s", trName)
}
