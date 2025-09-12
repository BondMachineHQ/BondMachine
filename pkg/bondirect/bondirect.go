package bondirect

import "github.com/BondMachineHQ/BondMachine/pkg/bmcluster"

type BondirectElement struct {
	*Config            // Configuration settings
	*bmcluster.Cluster // Reference to the cluster
	*Mesh              // Reference to the mesh
	*TData             // Metadata
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
	Type    string            `json:"Type"`
	Name    string            `json:"Name"`
	Signals map[string]Signal `json:"Signals"`
	Data    map[string]string `json:"Data"`
}

type Signal struct {
	Type string `json:"Type"`
	Name string `json:"Name"`
}

type Node struct {
	PeerId uint32            `json:"PeerId"`
	Data   map[string]string `json:"Data"`
}

type Edge struct {
	NodeA    string            `json:"NodeA"`
	NodeB    string            `json:"NodeB"`
	FromAtoB EdgeDirection     `json:"FromAtoB"`
	FromBtoA EdgeDirection     `json:"FromBtoA"`
	Data     map[string]string `json:"Data"`
}

type EdgeDirection struct {
	ATransceiver string            `json:"ATransceiver"`
	BTransceiver string            `json:"BTransceiver"`
	Data         map[string]string `json:"Data"`
}

type Path struct {
	NodeA string   `json:"NodeA"`
	NodeB string   `json:"NodeB"`
	Nodes []string `json:"Nodes"`
	Via   []string `json:"Via"`
}

type NodeMessages struct {
	PeerId            uint32
	Origins           *[]string
	OriginsNextHop    *[]string
	OriginsNextHopVia *[]string
	Destinations      *[]string
	Routes            *[]string
	RoutesNextHop     *[]string
	RoutesNextHopVia  *[]string
}
