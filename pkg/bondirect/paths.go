package bondirect

import "fmt"

func ShowPaths(c *Config, mesh *Mesh) {
	paths, err := GetPaths(c, mesh)
	if err != nil {
		fmt.Println("Error getting paths:", err)
		return
	}

	fmt.Println("Paths:")
	for _, path := range paths {
		fmt.Println(" -", path)
	}
}

func GetPaths(c *Config, mesh *Mesh) ([]Path, error) {
	paths := make([]Path, 0)

	for nodeIName := range mesh.Nodes {
		for nodeJName := range mesh.Nodes {
			path, err := GetPath(c, mesh, nodeIName, nodeJName)
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func GetPath(c *Config, mesh *Mesh, nodeA string, nodeB string) (Path, error) {
	if nodeA == nodeB {
		return Path{
			NodeA: nodeA,
			NodeB: nodeB,
			Nodes: []string{nodeA},
		}, nil
	}

	possiblePaths := make([][]string, 1)
	possiblePaths[0] = []string{nodeA}

	for len(possiblePaths) > 0 {
		newPossiblePaths := make([][]string, 0)
		for _, p := range possiblePaths {
			lastNode := p[len(p)-1]

			// Get neighbors of the last node
			neighbors, err := GetNeighbors(c, mesh, lastNode)
			if err != nil {
				return Path{}, err
			}

			for _, neighbor := range neighbors {
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
				copy(newPath, p)
				newPath = append(newPath, neighbor)

				if neighbor == nodeB {
					return Path{
						NodeA: nodeA,
						NodeB: nodeB,
						Nodes: newPath,
					}, nil
				}

				newPossiblePaths = append(newPossiblePaths, newPath)
			}
		}

		possiblePaths = newPossiblePaths
	}

	return Path{}, nil
}

func GetNeighbors(c *Config, mesh *Mesh, nodeName string) ([]string, error) {
	neighbors := make([]string, 0)

	_, ok := mesh.Nodes[nodeName]
	if !ok {
		return nil, fmt.Errorf("Node not found: %s", nodeName)
	}

	for _, edge := range mesh.Edges {
		if edge.NodeA == nodeName {
			neighbors = append(neighbors, edge.NodeB)
		}
		if edge.NodeB == nodeName {
			neighbors = append(neighbors, edge.NodeA)
		}
	}

	return neighbors, nil
}
