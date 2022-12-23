package bondirect

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func UnmarshallMesh(c *Config, mesh_file string) (*Mesh, error) {

	clus := new(Mesh)

	if _, err := os.Stat(mesh_file); err == nil {
		if jsonfile, err := ioutil.ReadFile(mesh_file); err == nil {
			if err := json.Unmarshal([]byte(jsonfile), clus); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return clus, nil
}

func UnmarshallCluster(c *Config, cluster_file string) (*Cluster, error) {

	clus := new(Cluster)

	if _, err := os.Stat(cluster_file); err == nil {
		if jsonfile, err := ioutil.ReadFile(cluster_file); err == nil {
			if err := json.Unmarshal([]byte(jsonfile), clus); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return clus, nil
}
