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
	if name, err := be.GetTransceiverName(nodeName, edgeName, direction); err == nil {
		be.TData.TransName = name
	} else {
		return "", err
	}

	if err := be.PopulateIOData(nodeName); err != nil {
		return "", err
	}
	if err := be.PopulateWireData(nodeName); err != nil {
		return "", err
	}

	be.PopulateTransParams(be.TData.TransName)

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

func (be *BondirectElement) GetTransceiverSignals(trName string) ([]string, []string, error) {
	signals := make([]string, 0)
	ports := make([]string, 0)

	sigs := 0
	for _, tr := range be.Mesh.Transceivers {
		if tr.Name == trName {
			for sName, s := range tr.Signals {
				if s.Type == "clock" {
					// Clock goes first
					pref := make([]string, 0)
					pref = append(pref, sName)
					signals = append(pref, signals...)
					pref2 := make([]string, 0)
					pref2 = append(pref2, s.Name)
					ports = append(pref2, ports...)
					sigs++
					continue
				}
				if s.Type == "data" {
					signals = append(signals, sName)
					ports = append(ports, s.Name)
					sigs++
					continue
				} else {
					return nil, nil, fmt.Errorf("unknown signal type: %s", s.Type)
				}
			}
			break
		}
	}

	if sigs >= 2 {
		return signals, ports, nil
	}
	return nil, nil, fmt.Errorf("not enough signals found for transceiver: %s", trName)
}

func (be *BondirectElement) GetTransceiverName(nodeName, lineName, direction string) (string, error) {
	// Using cluster names to find the mesh node name (that can be different)
	if meshNodeName, err := be.GetMeshNodeName(nodeName); err == nil {
		nodeName = meshNodeName
	} else {
		return "", fmt.Errorf("failed to get mesh node name: %v", err)
	}

	if line, exists := be.Mesh.Edges[lineName]; exists {
		if direction == "in" {
			if line.NodeA == nodeName {
				return line.FromBtoA.ATransceiver, nil
			} else if line.NodeB == nodeName {
				return line.FromAtoB.BTransceiver, nil
			}
		} else if direction == "out" {
			if line.NodeA == nodeName {
				return line.FromAtoB.ATransceiver, nil
			} else if line.NodeB == nodeName {
				return line.FromBtoA.BTransceiver, nil
			}
		} else {
			return "", fmt.Errorf("invalid direction: %s", direction)
		}
		return "", fmt.Errorf("node %s is not connected to line %s", nodeName, lineName)
	}
	return "", fmt.Errorf("line %s not found", lineName)
}
