package neuralbond

import (
	"fmt"
	"sort"
)

type TrainedNet struct {
	Nodes   []Node
	Weights []Weight
}

type Weight struct {
	Layer        int
	PosCurrLayer int
	PosPrevLayer int
	RelPosDown   int
	RelPosUp     int
	Value        float32
}

type Node struct {
	Layer   int
	Pos     int
	Type    string
	Bias    float32
	Inputs  int
	Outputs int
}

func (n *TrainedNet) Normalize() {
	for w, weight := range n.Weights {
		downL, downP := weight.Layer-1, weight.PosPrevLayer
		upL, upP := weight.Layer, weight.PosCurrLayer

		for i, node := range n.Nodes {
			if node.Layer == downL && node.Pos == downP {
				n.Nodes[i].Outputs++
			}
			if node.Layer == upL && node.Pos == upP {
				n.Nodes[i].Inputs++
			}
		}

		sameDown := make([]int, 0)
		sameUp := make([]int, 0)

		for _, chWeight := range n.Weights {
			if chWeight.Layer-1 == downL && chWeight.PosPrevLayer == downP {
				sameDown = append(sameDown, chWeight.PosCurrLayer)
			}
			if chWeight.Layer == upL && chWeight.PosCurrLayer == upP {
				sameUp = append(sameUp, chWeight.PosPrevLayer)
			}
		}

		// Sort the sameDown and sameUp arrays
		sort.Ints(sameDown)
		sort.Ints(sameUp)

		// fmt.Println(w, sameDown, sameUp)

		for i, v := range sameDown {
			if weight.PosCurrLayer == v {
				n.Weights[w].RelPosDown = i
				break
			}
		}

		for i, v := range sameUp {
			if weight.PosPrevLayer == v {
				n.Weights[w].RelPosUp = i
				break
			}
		}

	}
}

func (n *TrainedNet) WriteBasm() (string, error) {
	result := "%meta bmdef     global registersize:32\n"
	for _, node := range n.Nodes {
		if node.Type == "input" {
			result += fmt.Sprintf("%%meta cpdef node_0_%d romcode:terminal\n", node.Pos)
			result += fmt.Sprintf("%%meta iodef input_%d type:io\n", node.Pos)
			result += fmt.Sprintf("%%meta ioatt input_%d cp:node_0_%d, type:input, index:0\n", node.Pos, node.Pos)
			result += fmt.Sprintf("%%meta ioatt input_%d cp:bm, type:input, index:%d\n", node.Pos, node.Pos)
		} else if node.Type == "linear" {
			result += fmt.Sprintf("%%meta cpdef node_%d_%d romcode:linear, inputs:%d\n", node.Layer, node.Pos, node.Inputs)
		} else if node.Type == "softmax" {
			result += fmt.Sprintf("%%meta cpdef node_%d_%d romcode:softmax, inputs:%d\n", node.Layer, node.Pos, node.Inputs)
		} else if node.Type == "output" {
			result += fmt.Sprintf("%%meta cpdef node_%d_%d romcode:terminal\n", node.Layer, node.Pos)
			result += fmt.Sprintf("%%meta iodef output_%d type:io\n", node.Pos)
			result += fmt.Sprintf("%%meta ioatt output_%d cp:node_%d_%d, type:output, index:0\n", node.Pos, node.Layer, node.Pos)
			result += fmt.Sprintf("%%meta ioatt output_%d cp:bm, type:output, index:%d\n", node.Pos, node.Pos)
		} else {
			return "", fmt.Errorf("Unknown node type: %s", node.Type)
		}
	}

	for _, weight := range n.Weights {
		weightCP := fmt.Sprintf("weightcp_%d_%d__%d_%d", weight.Layer-1, weight.PosPrevLayer, weight.Layer, weight.PosCurrLayer)
		downNode := fmt.Sprintf("node_%d_%d", weight.Layer-1, weight.PosPrevLayer)
		upNode := fmt.Sprintf("node_%d_%d", weight.Layer, weight.PosCurrLayer)
		// result += fmt.Sprintf("%%meta cpdef %s romcode:weight, weight:%f\n", weightCP, weight.Value)
		result += fmt.Sprintf("%%meta cpdef %s romcode:weight\n", weightCP)
		result += fmt.Sprintf("%%meta iodef up%s type:io\n", weightCP)
		result += fmt.Sprintf("%%meta iodef down%s type:io\n", weightCP)
		result += fmt.Sprintf("%%meta ioatt down%s cp:%s, type:input, index:0\n", weightCP, weightCP)
		result += fmt.Sprintf("%%meta ioatt down%s cp:%s, type:output, index:0\n", weightCP, downNode)
		result += fmt.Sprintf("%%meta ioatt up%s cp:%s, type:input, index:%d\n", weightCP, upNode, weight.RelPosUp)
		result += fmt.Sprintf("%%meta ioatt up%s cp:%s, type:output, index:0\n", weightCP, weightCP)
	}

	return result, nil
}
