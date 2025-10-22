package bondirect

import (
	"fmt"
	"testing"
)

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
