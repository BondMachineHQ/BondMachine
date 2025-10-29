package bondirect

import "fmt"

func (be *BondirectElement) DumpNodeMetaData(nodeName string) (string, error) {
	if node, exists := be.Nodes[nodeName]; exists {
		result := ""
		for key, value := range node.Data {
			result += fmt.Sprintf("%s=%s\n", key, value)
		}
		return result, nil
	}
	return "", fmt.Errorf("Node %s not found", nodeName)
}

func (be *BondirectElement) DumpEdgeMetaData(edgeName string) (string, error) {
	if edge, exists := be.Edges[edgeName]; exists {
		result := ""
		for key, value := range edge.Data {
			result += fmt.Sprintf("%s=%s\n", key, value)
		}
		return result, nil
	}
	return "", fmt.Errorf("Edge %s not found", edgeName)
}
