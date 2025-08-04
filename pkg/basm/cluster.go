package basm

import "fmt"

type Peer struct {
	PeerId   uint32
	PeerName string
	Channels []uint32
	Inputs   []uint32
	Outputs  []uint32
}

type Cluster struct {
	ClusterId uint32
	Peers     []Peer
}

func (bi *BasmInstance) IsClustered() bool {
	if bi == nil {
		return false
	}
	return bi.isClustered
}

func (bi *BasmInstance) GetClusteredBondMachines() []string {
	if bi == nil {
		return nil
	}
	return bi.clusteredBondMachines
}

func (bi *BasmInstance) GetClusteredName() map[string]int {
	if bi == nil {
		return nil
	}
	return bi.clusteredNames
}

func (bi *BasmInstance) Assembler2Cluster() error {
	if bi == nil {
		return nil
	}
	if !bi.isClustered {
		return nil
	}

	for _, cp := range bi.cps {
		cpName := cp.GetValue()
		fmt.Println(cpName, cp)
		// TODO: Implement the logic to convert assembler to cluster
		// This is a placeholder for the actual conversion logic
	}

	// Loop through the clustered bond machines
	for edgeName, edgeId := range bi.clusteredNames {
		if bi.debug {
			fmt.Println(green("\tProcessing BM:"), red(edgeName), green("id"), blue(edgeId))
		}

	}

	return nil
}
