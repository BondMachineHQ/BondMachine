package etherbond

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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
