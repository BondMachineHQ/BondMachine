package main

import (
	"encoding/json"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/etherbond"

	//"errors"
	"flag"
	"fmt"

	//"log"
	"os"
	//"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	//"github.com/BondMachineHQ/BondMachine/pkg/simbox"
	//"sort"
	//"strconv"
	"strings"
)

type string_slice []string

func (i *string_slice) String() string {
	return fmt.Sprint(*i)
}

func (i *string_slice) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*i = append(*i, dt)
	}
	return nil
}

type BMCluster struct {
	Cluster_file      string
	Bondmachine_ids   []int
	Bondmachine_files []string
	Bondmachine_maps  []string
}

type Transformation struct {
	Transformations []string
}

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var clusterFile = flag.String("cluster-file", "", "Cluster JSON file")

var emit_dot = flag.Bool("emit-dot", false, "Emit dot file on stdout")
var dot_detail = flag.Int("dot-detail", 1, "Detail of infos on dot file 1-5")

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
	conf.Dotdetail = uint8(*dot_detail)

	clMain := new(BMCluster)
	if *clusterFile != "" {
		if rd_json, err := os.ReadFile(*clusterFile); err == nil {
			if err := json.Unmarshal([]byte(rd_json), clMain); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}

		// Create the Redeployer data struct starting from the json redeployer file
		rd := new(bondmachine.Redeployer)

		rd.Init()

		econfig := new(etherbond.Config)
		// TODO Import register size from whitin BM internals
		econfig.Rsize = uint8(8)

		if clMain.Cluster_file != "" {
			if cluster, err := bmcluster.UnmarshalCluster(clMain.Cluster_file); err != nil {
				panic(err)
			} else {
				rd.Cluster = cluster
			}
		} else {
			panic("Wrong Cluster file")
		}

		if len(clMain.Bondmachine_ids) != len(clMain.Bondmachine_files) {
			panic("Wrong number of files")
		}
		if len(clMain.Bondmachine_ids) != len(clMain.Bondmachine_maps) {
			panic("Wrong number of files")
		}

		for i, id := range clMain.Bondmachine_ids {
			bmach := new(bondmachine.Bondmachine)
			if bondmachine_json, err := os.ReadFile(clMain.Bondmachine_files[i]); err == nil {
				var bmachj bondmachine.Bondmachine_json
				if err := json.Unmarshal([]byte(bondmachine_json), &bmachj); err == nil {
					bmach = (&bmachj).Dejsoner()
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}

			rd.Bondmachines[id] = bmach

			iomap := new(bondmachine.IOmap)
			if mapfile_json, err := os.ReadFile(clMain.Bondmachine_maps[i]); err == nil {
				if err := json.Unmarshal([]byte(mapfile_json), iomap); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}

			rd.Maps[id] = iomap

		}

		if *emit_dot {
			fmt.Print(rd.Dot(conf))
		}

		f, err := os.Create(*clusterFile)
		check(err)
		b, errj := json.Marshal(clMain)
		check(errj)
		_, err = f.WriteString(string(b))
		check(err)
	}
}
