package bmcluster

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
