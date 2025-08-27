package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondirect"
)

var debug = flag.Bool("d", false, "Verbose")

// Bondirect mesh file and cluster spec file
var bondirectMesh = flag.String("bondirect-mesh", "", "Bondirect mesh File ")
var clusterSpec = flag.String("cluster-spec", "", "Cluster Spec File ")

// Operations
var showMessages = flag.Bool("show-messages", false, "Show messages")
var showPaths = flag.Bool("show-paths", false, "Show paths")

//TODO var listNodes = flag.Bool("list-nodes", false, "List nodes")
//TODO var listEdges = flag.Bool("list-edges", false, "List edges")

// The bondirect components are:
// - Transceiver: Handles the communication (in or out) on one end of a wire, they can be recv or send.
// - Wire or Edge: Connects two Transceivers, send+recv on both ends, also is a logic edge on the mesh
// - Node: Represents a logical endpoint in the mesh
// - Cluster: Represents a group of Nodes with the messages among them
// - Path: Represents a sequence of Nodes and Wires connecting them
// - Mesh: Represents the entire network of Nodes and Wires
// - Line: Represents a couple of Transceivers of a wire in a node.
// - Endpoint: If the elements that connects BM with all the wires in the mash

// So every BM has 1 element connected to the BM. It has as many lines as the wires
// Going to others BMs from that BM. Every line has 2 transceivers.

// Objects specify
var prefix = flag.String("prefix", "", "Prefix for all the generated names")

var nodeName = flag.String("node", "", "Node name")
var edgeName = flag.String("edge", "", "Edge name")
var direction = flag.String("direction", "", "Direction (in(recv)/out(send))")

// Generation
var generateTransceiver = flag.Bool("generate-transceiver", false, "Generate Transceiver")
var generateLine = flag.Bool("generate-line", false, "Generate Line")
var generateEndpoint = flag.Bool("generate-endpoint", false, "Generate Endpoint")

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
	} else {
		log.Fatal("Bondirect Mesh must be provided")
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
	} else {
		log.Fatal("Cluster Spec must be provided")
	}

	be := new(bondirect.BondirectElement)
	be.Config = c
	be.Mesh = myMesh
	be.Cluster = myCluster

	if c.Debug {
		fmt.Println("Cluster Spec:", myCluster)
	}

	if *showMessages {
		if myMesh == nil || myCluster == nil {
			fmt.Println("Both Bondirect Mesh and Cluster Spec must be provided to show messages.")
		} else {
			be.ShowMessages()
		}
	}

	if *showPaths {
		if myMesh == nil {
			fmt.Println("Bondirect Mesh must be provided to show paths.")
		} else {
			be.ShowPaths()
		}
	}

	if *emitMeshDot {
		if myMesh == nil {
			fmt.Println("Bondirect Mesh must be provided to emit Graphviz DOT.")
		} else {
			dot, err := be.EmitMeshDot()
			if err != nil {
				panic(err)
			} else {
				fmt.Println(dot)
			}
		}
	}
}
