package bondirect

import "github.com/BondMachineHQ/BondMachine/pkg/bmcluster"

type BondirectElement struct {
	*Config
	*bmcluster.Cluster
	*Mesh
}

// Config struct
type Config struct {
	Rsize uint8
	Debug bool
}

// Mesh JSON structures
type Mesh struct {
	Transceivers []Transceiver   `json:"Transceivers"`
	Nodes        map[string]Node `json:"Nodes"`
	Edges        map[string]Edge `json:"Edges"`
}

type Transceiver struct {
	Type    string                 `json:"Type"`
	Name    string                 `json:"Name"`
	Signals map[string]Signal      `json:"Signals"`
	Data    map[string]interface{} `json:"Data"`
}

type Signal struct {
	Type string `json:"Type"`
	Name string `json:"Name"`
}

type Node struct {
	PeerId uint32                 `json:"PeerId"`
	Data   map[string]interface{} `json:"Data"`
}

type Edge struct {
	NodeA    string                 `json:"NodeA"`
	NodeB    string                 `json:"NodeB"`
	FromAtoB EdgeDirection          `json:"FromAtoB"`
	FromBtoA EdgeDirection          `json:"FromBtoA"`
	Data     map[string]interface{} `json:"Data"`
}

type EdgeDirection struct {
	ATransceiver string                 `json:"ATransceiver"`
	BTransceiver string                 `json:"BTransceiver"`
	Data         map[string]interface{} `json:"Data"`
}

type Path struct {
	NodeA string   `json:"NodeA"`
	NodeB string   `json:"NodeB"`
	Nodes []string `json:"Nodes"`
}

type NodeMessages struct {
	PeerId             uint32
	Origins            *[]string
	OriginDestinations *[]string
	Destinations       *[]string
	Routes             *[]string
	RouteDestinations  *[]string
}
