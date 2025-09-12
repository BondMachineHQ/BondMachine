package bondirect

import (
	"fmt"
)

func (be *BondirectElement) ShowPaths() {

	paths, err := be.GetPaths()
	if err != nil {
		fmt.Println("Error getting paths:", err)
		return
	}

	fmt.Println("Paths:")
	for _, path := range paths {
		fmt.Println(" -", path)
	}

	messPaths, _ := be.SolveMessages()
	fmt.Println("Messages paths:")
	for _, mp := range messPaths {
		fmt.Println(" -", mp.PeerId, mp.Origins, mp.Destinations, mp.Routes, mp.OriginsNextHop, mp.OriginsNextHopVia, mp.RoutesNextHop, mp.RoutesNextHopVia)
	}
}

func (be *BondirectElement) GetPaths() ([]Path, error) {

	mesh := be.Mesh

	paths := make([]Path, 0)

	for nodeIName := range mesh.Nodes {
		for nodeJName := range mesh.Nodes {
			path, err := be.GetPath(nodeIName, nodeJName)
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func (be *BondirectElement) GetPath(nodeA string, nodeB string) (Path, error) {
	if nodeA == nodeB {
		return Path{
			NodeA: nodeA,
			NodeB: nodeB,
			Nodes: []string{nodeA},
			Via:   []string{"self"},
		}, nil
	}

	possiblePaths := make([][]string, 1)
	viaPaths := make([][]string, 1)
	possiblePaths[0] = []string{nodeA}
	viaPaths[0] = []string{"self"}

	for len(possiblePaths) > 0 {
		newPossiblePaths := make([][]string, 0)
		newViaPaths := make([][]string, 0)
		for j, p := range possiblePaths {
			lastNode := p[len(p)-1]

			// Get neighbors of the last node
			neighbors, via, err := be.GetNeighbors(lastNode)
			if err != nil {
				return Path{}, err
			}

			for i := range neighbors {
				neighbor := neighbors[i]
				viaP := via[i]

				// Avoid cycles
				alreadyInPath := false
				for _, n := range p {
					if n == neighbor {
						alreadyInPath = true
						break
					}
				}
				if alreadyInPath {
					continue
				}

				newPath := make([]string, len(p))
				newViaPath := make([]string, len(p))
				copy(newPath, p)
				copy(newViaPath, viaPaths[j])
				newPath = append(newPath, neighbor)
				newViaPath = append(newViaPath, viaP)

				if neighbor == nodeB {
					return Path{
						NodeA: nodeA,
						NodeB: nodeB,
						Nodes: newPath,
						Via:   newViaPath,
					}, nil
				}

				newPossiblePaths = append(newPossiblePaths, newPath)
				newViaPaths = append(newViaPaths, newViaPath)
			}
		}

		possiblePaths = newPossiblePaths
		viaPaths = newViaPaths
	}

	return Path{}, nil
}

func (be *BondirectElement) GetNeighbors(nodeName string) ([]string, []string, error) {
	neighbors := make([]string, 0)
	via := make([]string, 0)
	mesh := be.Mesh

	_, ok := mesh.Nodes[nodeName]
	if !ok {
		return nil, nil, fmt.Errorf("Node not found: %s", nodeName)
	}

	for edgeName, edge := range mesh.Edges {
		if edge.NodeA == nodeName {
			neighbors = append(neighbors, edge.NodeB)
			via = append(via, edgeName)
		}
		if edge.NodeB == nodeName {
			neighbors = append(neighbors, edge.NodeA)
			via = append(via, edgeName)
		}
	}

	return neighbors, via, nil
}
