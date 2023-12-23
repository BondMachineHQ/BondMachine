package bmgraph

import (
	"errors"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/goccy/go-graphviz/cgraph"
)

const (
	ASYNC = uint8(0) + iota
	SYNC
)

type Group []string
type Config struct {
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Pruned        []string
	Collapsed     []Group
	Debug         bool
	Verbose       bool
	NeuronLibPath string
	*bminfo.BMinfo
}

type Neuron struct {
	Params []string
}

type Graph struct {
	*cgraph.Graph
}

func (g *Graph) WriteBasm() (string, error) {
	if g == nil {
		return "", errors.New("Graph is nil")
	}

	// Find out all the vertices
	vertices := make(map[string]*cgraph.Node)
	for n := g.FirstNode(); n != nil; n = g.NextNode(n) {
		vertices[n.Name()] = n
		fmt.Println(n.Name())
	}

	for vn, v := range vertices {
		fmt.Println(vn)
		for e := g.FirstEdge(v); e != nil; e = g.NextEdge(e, v) {
			dest := e.Node()
			fmt.Println(v.Name(), dest.Name(), e.Name())
		}
	}

	return "", nil
}
