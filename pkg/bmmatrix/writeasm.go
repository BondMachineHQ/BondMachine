package bmmatrix

import (
	"errors"
	"fmt"
)

func (mo *MatrixOpertions) WriteBasm() (string, error) {
	if mo == nil {
		return "", errors.New("Graph is nil")
	}

	// config := mo.Config
	result := ""

	regSize := mo.RegisterSize
	result += fmt.Sprintf("%%meta bmdef     global registersize:%d\n", regSize)
	switch mo.IOMode {
	case ASYNC:
		result += fmt.Sprintf("%%meta bmdef     global iomode:async\n")
	case SYNC:
		result += fmt.Sprintf("%%meta bmdef     global iomode:sync\n")
	}
	// Find out all the vertices
	// vertices := make(map[string]*cgraph.Node)
	// verticesInputs := make(map[string]int)
	// verticesOutputs := make(map[string]int)
	// for n := g.FirstNode(); n != nil; n = g.NextNode(n) {
	// 	vertices[n.Name()] = n
	// 	verticesInputs[n.Name()] = 0
	// 	verticesOutputs[n.Name()] = 0
	// 	// result += "%meta fidef " + n.Name() + " fragment:" + n.Name() + "\n"
	// }

	// for _, v := range vertices {
	// 	for e := g.FirstEdge(v); e != nil; e = g.NextEdge(e, v) {
	// 		if e.Name() != "" {
	// 			dest := e.Node().Name()
	// 			src := v.Name()
	// 			result += "%meta filinkdef " + src + "_" + dest + " type:fl\n"
	// 			destIdx := verticesInputs[dest]
	// 			srcIdx := verticesOutputs[src]
	// 			result += "%meta filinkatt " + src + "_" + dest + " fi:" + src + ", type:output, index:" + fmt.Sprintf("%d", srcIdx) + "\n"
	// 			result += "%meta filinkatt " + src + "_" + dest + " fi:" + dest + ", type:input, index:" + fmt.Sprintf("%d", destIdx) + "\n"
	// 			verticesInputs[dest]++
	// 			verticesOutputs[src]++

	// 		}
	// 	}
	// }

	// for _, v := range vertices {
	// 	name := v.Name()
	// 	if config.UseNilNeuron {
	// 		name = "nil"
	// 	}
	// 	result += "%meta fidef " + v.Name() + " fragment:" + name + ", inputs:" + fmt.Sprintf("%d", verticesInputs[v.Name()]) + ", outputs:" + fmt.Sprintf("%d", verticesOutputs[v.Name()]) + "\n"
	// }

	// for _, v := range vertices {
	// 	result += "%meta cpdef " + v.Name() + " fragcollapse:" + v.Name() + "\n"
	// }

	return result, nil
}
