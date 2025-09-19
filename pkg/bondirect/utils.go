package bondirect

import (
	"fmt"
	"strconv"
)

func (be *BondirectElement) NamesConsistency() error {
	// Check that all nodes in the cluster have a corresponding node in the mesh
	// Eventually with different names
	// Also check that there are no duplicate names in either the cluster or the mesh
	// and also that names in the cluster and mesh do not overlap with different peerids
	usedNames := make(map[string]struct{})

	for _, cNode := range be.Cluster.Peers {
		found := false
		nodeName := cNode.PeerName
		for meshName, mNode := range be.Mesh.Nodes {
			if cNode.PeerId == mNode.PeerId {
				// Check for duplicate names
				if _, exists := usedNames[meshName]; exists {
					return fmt.Errorf("duplicate node name %s found in mesh", meshName)
				}
				if _, exists := usedNames[nodeName]; exists {
					return fmt.Errorf("duplicate node name %s found in cluster", nodeName)
				}
				usedNames[nodeName] = struct{}{}
				usedNames[meshName] = struct{}{}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("node %s with PeerId %s in cluster not found in mesh", nodeName, cNode.PeerId)
		}
	}

	return nil
}

func (be *BondirectElement) CheckClusterNodeName(nodeName string) bool {
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

func (be *BondirectElement) AnyNameToClusterName(nodeName string) (string, error) {
	if be.CheckClusterNodeName(nodeName) {
		return nodeName, nil
	}

	for n, mnode := range be.Mesh.Nodes {
		if n == nodeName {
			for _, cnode := range be.Cluster.Peers {
				if cnode.PeerId == mnode.PeerId {
					return cnode.PeerName, nil
				}
			}
			return "", fmt.Errorf("node %s found in mesh but not in cluster", nodeName)
		}
	}

	return "", fmt.Errorf("node %s not found in either cluster or mesh", nodeName)
}

func (be *BondirectElement) AnyNameToMeshName(nodeName string) (string, error) {
	if be.CheckMeshNodeName(nodeName) {
		return nodeName, nil
	}

	for _, cNode := range be.Cluster.Peers {
		if cNode.PeerName == nodeName {
			for n, mNode := range be.Mesh.Nodes {
				if cNode.PeerId == mNode.PeerId {
					return n, nil
				}
			}
			return "", fmt.Errorf("node %s found in cluster but not in mesh", nodeName)
		}
	}

	return "", fmt.Errorf("node %s not found in either cluster or mesh", nodeName)
}

func zerosPrefix(num int, value string) string {
	result := value
	for i := 0; i < num-len(value); i++ {
		result = "0" + result
	}
	return result
}

func getBinary(i int) string {
	result := strconv.FormatInt(int64(i), 2)
	return result
}
