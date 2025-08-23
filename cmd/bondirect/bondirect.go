package main

import (
	"flag"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondirect"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

// Bondirect mesh file and cluster spec file
var bondirectMesh = flag.String("bondirect-mesh", "", "Bondirect mesh File ")
var clusterSpec = flag.String("cluster-spec", "", "Cluster Spec File ")

// Operations
var showMessages = flag.Bool("show-messages", false, "Show messages")

// Graphviz
var emitMeshDot = flag.Bool("emit-mesh-dot", false, "Emit Graphviz DOT for the mesh")

func init() {
	flag.Parse()
}

func main() {

	c := new(bondirect.Config)

	var myMesh *bondirect.Mesh
	var myCluster *bmcluster.Cluster

	if *debug {
		c.Debug = true
	}

	if *bondirectMesh != "" {
		if mesh, err := bondirect.UnmarshalMesh(c, *bondirectMesh); err != nil {
			panic(err)
		} else {
			myMesh = mesh
		}
	}

	if c.Debug {
		fmt.Println("Bondirect Mesh:", myMesh)
	}

	if *clusterSpec != "" {
		if cluster, err := bmcluster.UnmarshalCluster(*clusterSpec); err != nil {
			panic(err)
		} else {
			myCluster = cluster
		}
	}

	if c.Debug {
		fmt.Println("Cluster Spec:", myCluster)
	}

	if *showMessages {
		if myMesh == nil || myCluster == nil {
			fmt.Println("Both Bondirect Mesh and Cluster Spec must be provided to show messages.")
		} else {
			bondirect.ShowMessages(c, myMesh, myCluster)
		}
	}

	if *emitMeshDot {
		if myMesh == nil {
			fmt.Println("Bondirect Mesh must be provided to emit Graphviz DOT.")
		} else {
			dot, err := bondirect.EmitMeshDot(c, myMesh)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(dot)
			}
		}
	}
}
