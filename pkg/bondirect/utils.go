package bondirect

import (
	"fmt"
)

func (be *BondirectElement) GetMeshNodeName(nodeName string) (string, error) {

	// Using cluster names to find the mesh node name (that can be different)
	found := false
	for _, node := range be.Cluster.Peers {
		if node.PeerName == nodeName {
			for n, mnode := range be.Mesh.Nodes {
				if node.PeerId == mnode.PeerId {
					nodeName = n
					break
				}
			}
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("node %s not found", nodeName)
	}

	return nodeName, nil
}

func (be *BondirectElement) CheckClusterNoneName(nodeName string) bool {
	// Check if the node name exists in the cluster
	for _, node := range be.Cluster.Peers {
		if node.PeerName == nodeName {
			return true
		}
	}
	return false
}

func (be *BondirectElement) CheckMeshNodeName(nodeName string) bool {
	// Check if the node name exists in the mesh and has a peerid also in the cluster
	for n, mnode := range be.Mesh.Nodes {
		if n == nodeName {
			for _, cnode := range be.Cluster.Peers {
				if cnode.PeerId == mnode.PeerId {
					return true
				}
			}
		}
	}
	return false
}
