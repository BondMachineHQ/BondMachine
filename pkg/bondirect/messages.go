package bondirect

import (
	"fmt"
	"strconv"
)

// ShowMessages displays the message flow between nodes in the cluster
func (be *BondirectElement) ShowMessages() {
	rawMess := be.Cluster.GetMessages()

	for _, rawMes := range rawMess {
		from := rawMes.From
		to := rawMes.To

		bmFrom := strconv.Itoa(from.BmId)
		idxFrom := strconv.Itoa(from.Index)

		bmTo := strconv.Itoa(to.BmId)
		idxTo := strconv.Itoa(to.Index)

		fmt.Printf("bm%sidx%stobm%sidx%sdata\n", bmFrom, idxFrom, bmTo, idxTo)
		fmt.Printf("bm%sidx%stobm%sidx%svalid\n", bmFrom, idxFrom, bmTo, idxTo)
		fmt.Printf("bm%sidx%stobm%sidx%srecv\n", bmTo, idxTo, bmFrom, idxFrom)
	}
}

func (be *BondirectElement) SolveMessages() (map[string]NodeMessages, error) {

	cluster := be.Cluster
	mesh := be.Mesh

	rawMess := cluster.GetMessages()

	nodeMessages := make(map[string]NodeMessages)

	for nodeName := range mesh.Nodes {
		nodeMessages[nodeName] = NodeMessages{
			PeerId:            mesh.Nodes[nodeName].PeerId,
			Origins:           &[]string{},
			OriginsNextHop:    &[]string{},
			Destinations:      &[]string{},
			Routes:            &[]string{},
			RoutesNextHop:     &[]string{},
			OriginsNextHopVia: &[]string{},
			RoutesNextHopVia:  &[]string{},
		}
	}

	for _, rawMes := range rawMess {
		from := rawMes.From
		to := rawMes.To

		peerIdFrom := uint32(from.BmId)
		bmFrom := strconv.Itoa(from.BmId)
		idxFrom := strconv.Itoa(from.Index)

		peerIdTo := uint32(to.BmId)
		bmTo := strconv.Itoa(to.BmId)
		idxTo := strconv.Itoa(to.Index)

		formName := ""
		toName := ""

		for nodeName := range mesh.Nodes {
			if mesh.Nodes[nodeName].PeerId == peerIdFrom {
				formName = nodeName
			}
			if mesh.Nodes[nodeName].PeerId == peerIdTo {
				toName = nodeName
			}
		}

		if formName == "" || toName == "" {
			return nil, fmt.Errorf("unable to find node names for peer IDs %d and %d", peerIdFrom, peerIdTo)
		}

		path, err := be.GetPath(formName, toName)
		if err != nil {
			return nil, fmt.Errorf("failed to get path from %s to %s: %w", formName, toName, err)
		}

		messDataName := fmt.Sprintf("bm%sidx%stobm%sidx%sdata", bmFrom, idxFrom, bmTo, idxTo)
		messValidName := fmt.Sprintf("bm%sidx%stobm%sidx%svalid", bmFrom, idxFrom, bmTo, idxTo)
		messRecvName := fmt.Sprintf("bm%sidx%stobm%sidx%srecv", bmFrom, idxFrom, bmTo, idxTo)

		// Process the path
		for i, step := range path.Nodes {
			fmt.Println(path.Nodes, path.Via)
			if i == 0 {
				// First step
				*nodeMessages[step].Origins = append(*nodeMessages[step].Origins, messDataName)
				*nodeMessages[step].OriginsNextHop = append(*nodeMessages[step].OriginsNextHop, path.Nodes[1])
				*nodeMessages[step].OriginsNextHopVia = append(*nodeMessages[step].OriginsNextHopVia, path.Via[1])
				*nodeMessages[step].Origins = append(*nodeMessages[step].Origins, messValidName)
				*nodeMessages[step].OriginsNextHop = append(*nodeMessages[step].OriginsNextHop, path.Nodes[1])
				*nodeMessages[step].OriginsNextHopVia = append(*nodeMessages[step].OriginsNextHopVia, path.Via[1])
				*nodeMessages[step].Destinations = append(*nodeMessages[step].Destinations, messRecvName)
			} else if i == len(path.Nodes)-1 {
				// Last step
				*nodeMessages[step].Destinations = append(*nodeMessages[step].Destinations, messDataName)
				*nodeMessages[step].Destinations = append(*nodeMessages[step].Destinations, messValidName)
				*nodeMessages[step].Origins = append(*nodeMessages[step].Origins, messRecvName)
				*nodeMessages[step].OriginsNextHop = append(*nodeMessages[step].OriginsNextHop, path.Nodes[i-1])
				*nodeMessages[step].OriginsNextHopVia = append(*nodeMessages[step].OriginsNextHopVia, path.Via[i])
			} else {
				// Intermediate step
				*nodeMessages[step].Routes = append(*nodeMessages[step].Routes, messDataName)
				*nodeMessages[step].RoutesNextHop = append(*nodeMessages[step].RoutesNextHop, path.Nodes[i+1])
				*nodeMessages[step].RoutesNextHopVia = append(*nodeMessages[step].RoutesNextHopVia, path.Via[i+1])
				*nodeMessages[step].Routes = append(*nodeMessages[step].Routes, messValidName)
				*nodeMessages[step].RoutesNextHop = append(*nodeMessages[step].RoutesNextHop, path.Nodes[i+1])
				*nodeMessages[step].RoutesNextHopVia = append(*nodeMessages[step].RoutesNextHopVia, path.Via[i+1])
				*nodeMessages[step].Routes = append(*nodeMessages[step].Routes, messRecvName)
				*nodeMessages[step].RoutesNextHop = append(*nodeMessages[step].RoutesNextHop, path.Nodes[i-1])
				*nodeMessages[step].RoutesNextHopVia = append(*nodeMessages[step].RoutesNextHopVia, path.Via[i])
			}
		}
	}
	return nodeMessages, nil
}
