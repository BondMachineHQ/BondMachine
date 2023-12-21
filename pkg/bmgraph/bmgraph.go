package bmgraph

import (
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
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
