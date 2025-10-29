package bondmachine

import (
	//"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
)

type peerlist map[int]int

type Redeployer struct {
	Cluster      *bmcluster.Cluster
	Bondmachines map[int]*Bondmachine
	Maps         map[int]*IOmap
}

func (rd *Redeployer) Init() {
	rd.Cluster = nil
	rd.Bondmachines = make(map[int]*Bondmachine)
	rd.Maps = make(map[int]*IOmap)
}

func (rd *Redeployer) Validate() bool {
	return false
}

func (rd *Redeployer) Dot(conf *Config) string {
	cl := rd.Cluster

	result := ""
	result += "digraph callgraph {\n"
	result += "\tbgcolor=\"#a2a6a7\";\n"
	result += "\tcompound=true;\n"
	result += "\tnode [fontname=\"verdana\"];\n"
	result += "\tfontname=\"Verdana\";\n"

	allinputs := make(map[uint32]peerlist, 0)
	alloutputs := make(map[uint32]peerlist, 0)
	allchannels := make(map[uint32]peerlist, 0)

	for i, peer := range cl.Peers {

		peerid := int(peer.PeerId)
		peerName := peer.PeerName

		result += "\tsubgraph cluster_b" + strconv.Itoa(int(i)) + " {\n"
		if GV_config(GVPEER) != "" {
			result += "\t" + GV_config(GVPEER) + ";\n"
		}
		result += "\t\tlabel=\"BM " + strconv.Itoa(peerid) + ": " + peerName + "\";\n"
		inps := int(len(peer.Inputs))
		outs := int(len(peer.Outputs))
		chs := int(len(peer.Channels))
		result += "\t\tsubgraph cluster_b" + strconv.Itoa(int(i)) + "_inputs {\n"
		if GV_config(GVCLUSININPEER) != "" {
			result += "\t\t" + GV_config(GVCLUSININPEER) + ";\n"
		}
		result += "\t\t\tlabel=\"Inputs\";\n"
		for j := 0; j < inps; j++ {
			if _, ok := allinputs[peer.Inputs[j]]; !ok {
				allinputs[peer.Inputs[j]] = map[int]int{int(i): j}
			} else {
				allinputs[peer.Inputs[j]][int(i)] = j
			}
			result += "\t\t\tnode [label=\"" + strconv.Itoa(int(peer.Inputs[j])) + "\" " + GV_config(GVNODEININPEER) + "] b" + strconv.Itoa(int(i)) + "i" + strconv.Itoa(j) + ";\n"
		}
		result += "\t\t}\n"
		result += "\t\tsubgraph cluster_b" + strconv.Itoa(int(i)) + "_outputs {\n"
		if GV_config(GVCLUSOUTINPEER) != "" {
			result += "\t\t" + GV_config(GVCLUSOUTINPEER) + ";\n"
		}
		result += "\t\t\tlabel=\"Outputs\";\n"
		for j := 0; j < outs; j++ {
			if _, ok := alloutputs[peer.Outputs[j]]; !ok {
				alloutputs[peer.Outputs[j]] = map[int]int{int(i): j}
			} else {
				alloutputs[peer.Outputs[j]][int(i)] = j
			}
			result += "\t\tnode [label=\"" + strconv.Itoa(int(peer.Outputs[j])) + "\" " + GV_config(GVNODEOUTINPEER) + "] b" + strconv.Itoa(int(i)) + "o" + strconv.Itoa(j) + ";\n"
		}
		result += "\t\t}\n"
		result += "\t\tsubgraph cluster_b" + strconv.Itoa(int(i)) + "_channels {\n"
		if GV_config(GVCLUSCHINPEER) != "" {
			result += "\t\t" + GV_config(GVCLUSCHINPEER) + ";\n"
		}
		result += "\t\t\tlabel=\"Channels\";\n"
		for j := 0; j < chs; j++ {
			if _, ok := allchannels[peer.Channels[j]]; !ok {
				allchannels[peer.Channels[j]] = map[int]int{int(i): j}
			} else {
				allchannels[peer.Channels[j]][int(i)] = j
			}
			result += "\t\tnode [label=\"" + strconv.Itoa(int(peer.Channels[j])) + "\" " + GV_config(GVNODECHINPEER) + "] b" + strconv.Itoa(int(i)) + "ch" + strconv.Itoa(j) + ";\n"
		}
		result += "\t\t}\n"

		if bmach, ok := rd.Bondmachines[peerid]; ok {
			result += bmach.Dot(conf, "id"+strconv.Itoa(int(i)), nil, nil)
		}

		if bmap, ok := rd.Maps[peerid]; ok {
			for objname, objid := range bmap.Assoc {
				for j := 0; j < inps; j++ {
					if strconv.Itoa(int(peer.Inputs[j])) == objid {
						result += "\t\tid" + strconv.Itoa(int(i)) + objname + " -> b" + strconv.Itoa(int(i)) + "i" + strconv.Itoa(j) + "[arrowhead=none];\n"
					}
				}
				for j := 0; j < outs; j++ {
					if strconv.Itoa(int(peer.Outputs[j])) == objid {
						result += "\t\tid" + strconv.Itoa(int(i)) + objname + " -> b" + strconv.Itoa(int(i)) + "o" + strconv.Itoa(j) + "[arrowhead=none];\n"
					}
				}
			}
		}

		result += "\t}\n"

	}

	for j, _ := range allchannels {
		result += "\tnode [label=\"ch" + strconv.Itoa(int(j)) + "\" " + GV_config(GVNODECHINPEER) + "] " + "ch" + strconv.Itoa(int(j)) + ";\n"
	}

	for i, outmap := range alloutputs {
		if inmap, ok := allinputs[i]; ok {
			for outpeer, outpos := range outmap {
				for inpeer, inpos := range inmap {
					result += "\tb" + strconv.Itoa(outpeer) + "o" + strconv.Itoa(outpos) + " -> b" + strconv.Itoa(inpeer) + "i" + strconv.Itoa(inpos) + ";\n"
				}
			}
		}
	}

	for j, chanmap := range allchannels {
		for peer, peerpos := range chanmap {
			result += "\tch" + strconv.Itoa(int(j)) + " -> b" + strconv.Itoa(peer) + "ch" + strconv.Itoa(peerpos) + "[arrowhead=none];\n"
		}
	}

	result += "}\n"

	return result
}
