package etherbond

import (
	"strconv"
)

const (
	GVPEER = uint8(0) + iota

	GVNODEININPEER
	GVNODEOUTINPEER
	GVNODECHINPEER

	GVCLUSININPEER
	GVCLUSOUTINPEER
	GVCLUSCHINPEER
)

func GV_config(element uint8) string {
	result := ""
	switch element {
	case GVPEER:
		result += "style=\"filled, rounded\" fillcolor=coral color=grey30"
	case GVNODEININPEER:
		result += "style=\"filled\" shape=box fillcolor=lightskyblue color=black"
	case GVNODEOUTINPEER:
		result += "style=\"filled\" shape=box fillcolor=indianred3 color=black"
	case GVNODECHINPEER:
		result += "style=\"filled\" shape=box fillcolor=red color=black"
	case GVCLUSININPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSOUTINPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey50"
	case GVCLUSCHINPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey30"
	}
	return result
}

//const (
//	GVNODE = uint8(0) + iota
//	GVCLUS
//	GVCLUSPROC
//	GVCLUSIN
//	GVCLUSOUT
//	GVNODEIN
//	GVNODEOUT
//
//	GVNODEINPROC
//
//	GVEDGE
//)
//
//func GV_config(element uint8) string {
//	result := ""
//	switch element {
//	case GVNODEIN:
//		result += "style=filled fillcolor=greenyellow color=black"
//	case GVNODEOUT:
//		result += "style=filled fillcolor=lightcoral color=black"
//	case GVCLUSIN:
//		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey80"
//	case GVCLUSOUT:
//		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey50"
//	case GVCLUSPROC:
//		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=aquamarine3"
//	}
//	return result
//}

func (cl *Cluster) Dot() string {

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
		result += "\tsubgraph cluster_b" + strconv.Itoa(int(i)) + " {\n"
		if GV_config(GVPEER) != "" {
			result += "\t" + GV_config(GVPEER) + ";\n"
		}
		result += "\t\tlabel=\"BM " + strconv.Itoa(int(i)) + "\";\n"
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
