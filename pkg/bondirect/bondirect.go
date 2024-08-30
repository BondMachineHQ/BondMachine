package bondirect

// Config struct
type Config struct {
	Rsize uint8
	Debug bool
}

type Peer struct {
	PeerId   uint32
	Channels []uint32
	Inputs   []uint32
	Outputs  []uint32
}

type Cluster struct {
	ClusterId uint32
	Peers     []Peer
}

// Mesh description

type EdgesList []string

type NodesParams struct {
	Data map[string]string
}

type EdgesParams struct {
	From  string
	To    string
	Wires uint8
	Clock uint8
	Data  map[string]string
}

type Mesh struct {
	Adjacency map[string]EdgesList
	Nodes     map[string]NodesParams
	Edges     map[string]EdgesParams
}

//

type Ips struct {
	Assoc map[string]string
}
