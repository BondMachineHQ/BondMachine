package bondirect

import (
	"fmt"
)

func (be *BondirectElement) PopulateNodeParams(nodeName string) {
	params := make(map[string]string)

	// Using cluster names to find the mesh node name (that can be different)
	if meshNodeName, err := be.AnyNameToMeshName(nodeName); err == nil {
		nodeName = meshNodeName
	} else {
		return
	}

	for mNodeName, mNode := range be.Mesh.Nodes {
		if mNodeName == nodeName {
			for pName, pValue := range mNode.Data {
				params[pName] = pValue
			}
			break
		}
	}

	be.TData.NodeParams = params
}

func (be *BondirectElement) PopulateEdgeParams(edgeName string) {
	params := make(map[string]string)

	for mEdgeName, mEdge := range be.Mesh.Edges {
		if mEdgeName == edgeName {
			for pName, pValue := range mEdge.Data {
				params[pName] = pValue
			}
			break
		}
	}

	be.TData.EdgeParams = params
}

func (be *BondirectElement) PopulateTransParams(transName string) {
	params := make(map[string]string)
	params["CountersLen"] = "32" // Default values
	params["ClkGraceWait"] = "10000"
	params["OutClockWait"] = "20000"
	params["ClkTimeout"] = "1000000"
	params["NumWires"] = "1" // Default values

	if numWires, _, err := be.GetTransceiverSignals(transName); err == nil {
		params["NumWires"] = fmt.Sprintf("%d", len(numWires)-1) // Exclude clock
	}

	for _, mTrans := range be.Mesh.Transceivers {

		if mTrans.Name == transName {
			for pName, pValue := range mTrans.Data {
				params[pName] = pValue
			}
			break
		}
	}

	be.TData.TransParams = params
}
