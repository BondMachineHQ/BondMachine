package neuralbond

import "fmt"

type TrainedNet struct {
	Nodes   []Node
	Weights []Weight
}

type Weight struct {
	Layer        int
	PosCurrLayer int
	PosPrevLayer int
	Value        float32
}

type Node struct {
	Layer int
	Pos   int
	Type  string
	Bias  float32
}

func (n *TrainedNet) WriteBasm() (string, error) {
	result := ""
	for _, node := range n.Nodes {
		if node.Type == "input" {
			result += fmt.Sprintf("%meta cpdef input_%d romcode:input\n", node.Pos)
		} else if node.Type == "hidden" {
			// fmt.Println("hidden")
		} else if node.Type == "output" {
			// fmt.Println("output")
		} else {
			return "", fmt.Errorf("Unknown node type: %s", node.Type)
		}
	}
	return result, nil
}
