package bondirect

import (
	"encoding/json"
	"os"
)

func UnmarshalMesh(c *Config, mesh_file string) (*Mesh, error) {

	cluster := new(Mesh)

	if _, err := os.Stat(mesh_file); err == nil {
		if jsonFile, err := os.ReadFile(mesh_file); err == nil {
			if err := json.Unmarshal([]byte(jsonFile), cluster); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return cluster, nil
}

func UnmarshalCluster(c *Config, cluster_file string) (*Cluster, error) {

	cluster := new(Cluster)

	if _, err := os.Stat(cluster_file); err == nil {
		if jsonFile, err := os.ReadFile(cluster_file); err == nil {
			if err := json.Unmarshal([]byte(jsonFile), cluster); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return cluster, nil
}
