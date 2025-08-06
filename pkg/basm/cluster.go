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

func (bi *BasmInstance) GetCluster() *Cluster {
	if bi == nil {
		return nil
	}
	return bi.cluster
}

func (bi *BasmInstance) Assembler2Cluster() error {
	if bi == nil {
		return nil
	}
	if !bi.isClustered {
		return nil
	}

	// Loop through the clustered bond machines to create peers first
	if bi.debug {
		fmt.Println(green("Creating cluster and peers"))
	}

	for edgeName, edgeId := range bi.clusteredNames {
		if bi.debug {
			fmt.Println(green("\tProcessing BM:"), red(edgeName), green("id"), blue(edgeId))
		}

		// Create a new peer for the cluster
		peer := Peer{
			PeerId:   uint32(edgeId),
			PeerName: edgeName,
			Channels: make([]uint32, 0),
			Inputs:   make([]uint32, 0),
			Outputs:  make([]uint32, 0),
		}

		// Add the peer to the cluster
		if bi.debug {
			fmt.Println(green("\tAdding BM:"), red(peer.PeerName), green("with id"), blue(peer.PeerId), green("to cluster"))
		}
		bi.cluster.Peers = append(bi.cluster.Peers, peer)

	}

	// Device for each clustered bond machine
	cpDev := make(map[string]string)

	// Now loop through the clustered bond machines to create metadata
	if bi.debug {
		fmt.Println(green("Creating metadata for clustered bond machines"))
	}
	for edgeName, edgeId := range bi.clusteredNames {
		if bi.debug {
			fmt.Println(green("\tProcessing BM:"), red(edgeName), green("id"), blue(edgeId))
		}

		// Write the node cps
		for _, cp := range bi.cps {
			devName := "default"
			if cp.GetMeta("device") != "" {
				devName = cp.GetMeta("device")
			}

			if devName == edgeName {
				cpDev[cp.GetValue()] = devName
				if bi.debug {
					fmt.Println(green("\t\tAdding cp:"), red(cp.GetValue()), green("to BM:"), red(edgeName))
				}
				bi.clusteredBondMachines[edgeId] += "%meta cpdef " + cp.GetValue()
				for key, value := range cp.LoopMeta() {
					if key != "templated" && key != "device" {
						bi.clusteredBondMachines[edgeId] += " " + key + ":" + value
					}
				}
				bi.clusteredBondMachines[edgeId] += "\n"
			}
		}
	}

	linkDone := make(map[string]struct{})

	// Write the node ioAttach
	for l1Idx, l1 := range bi.ioAttach {
		l1Name := l1.GetValue()

		// Skip if the link is already processed
		if _, ok := linkDone[l1Name]; ok {
			continue
		}

		// Find the other end of the link
		for l2Idx, l2 := range bi.ioAttach {
			l2Name := l2.GetValue()

			if l1Idx != l2Idx && l1Name == l2Name {

				l1CpName := l1.GetMeta("cp")
				l2CpName := l2.GetMeta("cp")

				if l1CpName == "bm" {

				}

				if l2CpName == "bm" {

				}
				l1DevName := "default"
				if d := l1.GetMeta("device"); d != "" {
					l1DevName = d
				}
				l2DevName := "default"
				if d := l2.GetMeta("device"); d != "" {
					l2DevName = d
				}

				// TODO Finish it
			}
		}
	}

	return nil
}
