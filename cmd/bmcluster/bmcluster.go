package main

import (
	"encoding/json"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"

	//"errors"
	"flag"
	"fmt"

	//"log"
	"os"
	//"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	//"github.com/BondMachineHQ/BondMachine/pkg/simbox"
	//"sort"
	//"strconv"
)

type BMCluster struct {
	clFile string
	bmIds  []int
	bmFile []string
	bmMaps []string
}

type Transformation struct {
	Transformations []string
}

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var clusterFile = flag.String("cluster-file", "", "Cluster JSON file")

var emitDot = flag.Bool("emit-dot", false, "Emit dot file on stdout")
var dotDetail = flag.Int("dot-detail", 1, "Detail of infos on dot file 1-5")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	flag.Parse()
}

func main() {
	conf := new(bondmachine.Config)
	conf.Debug = *debug
	conf.Dotdetail = uint8(*dotDetail)

	clMain := new(BMCluster)
	if *clusterFile != "" {
		if clJSON, err := os.ReadFile(*clusterFile); err == nil {
			if err := json.Unmarshal([]byte(clJSON), clMain); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}

		// Create the Redeployer data struct starting from the json redeployer file
		rd := new(bondmachine.Redeployer)

		rd.Init()

		if clMain.clFile != "" {
			if cluster, err := bmcluster.UnmarshalCluster(clMain.clFile); err != nil {
				panic(err)
			} else {
				rd.Cluster = cluster
			}
		} else {
			panic("Wrong Cluster file")
		}

		if len(clMain.bmIds) != len(clMain.bmFile) {
			panic("Wrong number of files")
		}
		if len(clMain.bmIds) != len(clMain.bmMaps) {
			panic("Wrong number of files")
		}

		for i, id := range clMain.bmIds {
			bMach := new(bondmachine.Bondmachine)
			if bondmachineJSON, err := os.ReadFile(clMain.bmFile[i]); err == nil {
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
			if mapFileJSON, err := os.ReadFile(clMain.bmMaps[i]); err == nil {
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
