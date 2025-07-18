package bondirect

import (
	"fmt"
	"testing"
)

func TestImportCluster(t *testing.T) {

	config := new(Config)
	var mycluster Cluster

	if cluster, err := UnmarshalCluster(config, "bondirect_test_cluster.json"); err != nil {
		panic(err)
	} else {
		mycluster = *cluster
	}
	fmt.Println(mycluster)
}

func TestImportMesh(t *testing.T) {

	config := new(Config)
	var mymesh Mesh

	if mesh, err := UnmarshalMesh(config, "bondirect_test_mesh.json"); err != nil {
		panic(err)
	} else {
		mymesh = *mesh
	}
	fmt.Println(mymesh)
}
