package bmcluster

import (
	"encoding/json"
	"os"
)

type Peer struct {
	PeerId   uint32
	PeerName string
	Channels []uint32
	Inputs   []uint32
	Outputs  []uint32
}

type Cluster struct {
	ClusterId uint32
	Peers     []Peer
}

func UnmarshalCluster(clusterFile string) (*Cluster, error) {

	cluster := new(Cluster)

	if _, err := os.Stat(clusterFile); err == nil {
		if jsonFile, err := os.ReadFile(clusterFile); err == nil {
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

func MarshalCluster(cluster *Cluster, clusterFile string) error {
	if jsonData, err := json.MarshalIndent(cluster, "", "  "); err != nil {
		return err
	} else {
		if err := os.WriteFile(clusterFile, jsonData, 0644); err != nil {
			return err
		}
	}
	return nil
}
