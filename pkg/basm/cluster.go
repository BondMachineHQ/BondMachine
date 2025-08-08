package basm

import (
	"fmt"
	"strconv"
)

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

	// Device for each cp
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
				meta := "%meta cpdef " + cp.GetValue()
				for key, value := range cp.LoopMeta() {
					if key != "templated" && key != "device" {
						meta += " " + key + ":" + value + ","
					}
				}
				bi.clusteredBondMachines[edgeId] += meta[:len(meta)-1] + "\n"
			}
		}
	}

	linkDone := make(map[string]struct{})

	// Create a structure that holds the I/O counter for every node
	iCounter := make([]int, len(bi.clusteredBondMachines))
	oCounter := make([]int, len(bi.clusteredBondMachines))

	for i := 0; i < len(bi.clusteredBondMachines); i++ {
		iCounter[i] = 0
		oCounter[i] = 0
	}

	linkIdx := 0

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
					other := l2CpName
					devName := cpDev[other]
					devId := bi.clusteredNames[devName]
					meta := "%meta ioatt " + l2.GetValue()
					for key, value := range l2.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"

					meta = "%meta ioatt " + l1.GetValue()
					for key, value := range l1.LoopMeta() {
						if key == "index" {
							if l1.GetMeta("type") == "input" {
								meta += " " + key + ":" + strconv.Itoa(iCounter[devId]) + ","
								iCounter[devId]++
							} else {
								meta += " " + key + ":" + strconv.Itoa(oCounter[devId]) + ","
								oCounter[devId]++
							}
						} else {
							meta += " " + key + ":" + value + ","
						}
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"
					linkDone[l1Name] = struct{}{}
					break
				}

				if l2CpName == "bm" {
					other := l1CpName
					devName := cpDev[other]
					devId := bi.clusteredNames[devName]
					meta := "%meta ioatt " + l1.GetValue()
					for key, value := range l1.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"

					meta = "%meta ioatt " + l2.GetValue()
					for key, value := range l2.LoopMeta() {
						if key == "index" {
							if l2.GetMeta("type") == "input" {
								meta += " " + key + ":" + strconv.Itoa(iCounter[devId]) + ","
								iCounter[devId]++
							} else {
								meta += " " + key + ":" + strconv.Itoa(oCounter[devId]) + ","
								oCounter[devId]++
							}
						} else {
							meta += " " + key + ":" + value + ","
						}
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"
					linkDone[l2Name] = struct{}{}
					break
				}

				l1DevName := cpDev[l1CpName]
				l2DevName := cpDev[l2CpName]

				fmt.Println("l1DevName:", l1DevName, "l2DevName:", l2DevName)
				if l1DevName == l2DevName {
					// The bonds are within the same device
					devId := bi.clusteredNames[l1DevName]
					meta := "%meta ioatt " + l1.GetValue()
					for key, value := range l1.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"

					meta = "%meta ioatt " + l2.GetValue()
					for key, value := range l2.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[devId] += meta[:len(meta)-1] + "\n"
					linkDone[l1Name] = struct{}{}
					break
				} else {
					dev1Id := bi.clusteredNames[l1DevName]
					meta := "%meta ioatt " + l1.GetValue()
					for key, value := range l1.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[dev1Id] += meta[:len(meta)-1] + "\n"

					// Get the peer index
					peerIdx := 0
					for i, p := range bi.cluster.Peers {
						if p.PeerName == l1DevName {
							peerIdx = i
							break
						}
					}

					meta = "%meta ioatt " + l1.GetValue()
					for key, value := range l1.LoopMeta() {
						switch key {
						case "index":
							if l1.GetMeta("type") == "input" {
								meta += " " + key + ":" + strconv.Itoa(iCounter[dev1Id]) + ","
								bi.cluster.Peers[peerIdx].Inputs = append(bi.cluster.Peers[peerIdx].Inputs, uint32(linkIdx))
								iCounter[dev1Id]++
							} else {
								meta += " " + key + ":" + strconv.Itoa(oCounter[dev1Id]) + ","
								bi.cluster.Peers[peerIdx].Outputs = append(bi.cluster.Peers[peerIdx].Outputs, uint32(linkIdx))
								oCounter[dev1Id]++
							}
						case "cp":
							meta += " " + key + ":bm,"
						default:
							meta += " " + key + ":" + value + ","
						}
					}
					bi.clusteredBondMachines[dev1Id] += meta[:len(meta)-1] + "\n"

					dev2Id := bi.clusteredNames[l2DevName]
					meta = "%meta ioatt " + l2.GetValue()
					for key, value := range l2.LoopMeta() {
						meta += " " + key + ":" + value + ","
					}
					bi.clusteredBondMachines[dev2Id] += meta[:len(meta)-1] + "\n"

					// Get the peer index
					for i, p := range bi.cluster.Peers {
						if p.PeerName == l2DevName {
							peerIdx = i
							break
						}
					}

					meta = "%meta ioatt" + l2.GetValue()
					for key, value := range l2.LoopMeta() {
						switch key {
						case "index":
							if l2.GetMeta("type") == "input" {
								meta += " " + key + ":" + strconv.Itoa(iCounter[dev2Id]) + ","
								bi.cluster.Peers[peerIdx].Inputs = append(bi.cluster.Peers[peerIdx].Inputs, uint32(linkIdx))
								iCounter[dev2Id]++
							} else {
								meta += " " + key + ":" + strconv.Itoa(oCounter[dev2Id]) + ","
								bi.cluster.Peers[peerIdx].Outputs = append(bi.cluster.Peers[peerIdx].Outputs, uint32(linkIdx))
								oCounter[dev2Id]++
							}
						case "cp":
							meta += " " + key + ":bm,"
						default:
							meta += " " + key + ":" + value + ","
						}
					}
					bi.clusteredBondMachines[dev2Id] += meta[:len(meta)-1] + "\n"

					linkDone[l1Name] = struct{}{}
					linkIdx++
				}

			}
		}
	}
	return nil
}
