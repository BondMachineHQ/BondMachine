package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var clusterInfoFile = flag.String("clusterinfo-file", "", "Cluster Info JSON file")

var emitDot = flag.Bool("emit-dot", false, "Emit dot file on stdout")
var dotDetail = flag.Int("dot-detail", 1, "Detail of infos on dot file 1-5")

func init() {
	flag.Parse()
}

func main() {
	conf := new(bondmachine.Config)
	conf.Debug = *debug
	conf.Dotdetail = uint8(*dotDetail)

	clMain := new(bmcluster.ClusterInfo)
	if *clusterInfoFile != "" {
		if clJSON, err := os.ReadFile(*clusterInfoFile); err == nil {
			if err := json.Unmarshal([]byte(clJSON), clMain); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}

		// Create the Redeployer data struct starting from the json redeployer file
		rd := new(bondmachine.Redeployer)

		rd.Init()

		if clMain.ClusterFile != "" {
			if cluster, err := bmcluster.UnmarshalCluster(clMain.ClusterFile); err != nil {
				panic(err)
			} else {
				rd.Cluster = cluster
			}
		} else {
			panic("Wrong Cluster file")
		}

		if len(clMain.BMIds) != len(clMain.BMFiles) {
			panic("Wrong number of files")
		}
		if len(clMain.BMIds) != len(clMain.BMMaps) {
			panic("Wrong number of files")
		}

		for i, id := range clMain.BMIds {
			bMach := new(bondmachine.Bondmachine)
			if bondmachineJSON, err := os.ReadFile(clMain.BMFiles[i]); err == nil {
				var bMachJ bondmachine.Bondmachine_json
				if err := json.Unmarshal([]byte(bondmachineJSON), &bMachJ); err == nil {
					bMach = (&bMachJ).Dejsoner()
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}

			rd.Bondmachines[id] = bMach

			ioMap := new(bondmachine.IOmap)
			if mapFileJSON, err := os.ReadFile(clMain.BMMaps[i]); err == nil {
				if err := json.Unmarshal([]byte(mapFileJSON), ioMap); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}

			rd.Maps[id] = ioMap

		}

		if *emitDot {
			fmt.Print(rd.Dot(conf))
		}
	}
}
