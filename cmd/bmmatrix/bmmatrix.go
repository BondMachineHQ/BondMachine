package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Debug")

var registerSize = flag.Int("register-size", 32, "Number of bits per register (n-bit)")
var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")

var bmInfoFile = flag.String("bminfo-file", "", "JSON description of the BondMachine abstraction")

var iomode = flag.String("io-mode", "async", "IO mode: async, sync")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
}

func main() {

	// Create the config struct
	config := new(bmmatrix.Config)

	config.Debug = *debug
	config.Verbose = *verbose

	// Load or create the Info file
	config.BMinfo = new(bminfo.BMinfo)

	if *bmInfoFile != "" {
		if bmInfoJSON, err := os.ReadFile(*bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, config.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	if *neuronLibPath != "" {
		config.NeuronLibPath = *neuronLibPath
	} else {
		// panic("No neuron library path specified")
	}

	if *dataType != "" {
		found := false
		for _, tpy := range bmnumbers.AllTypes {
			if tpy.GetName() == *dataType {
				for opType, opName := range tpy.ShowInstructions() {
					config.Params[opType] = opName
				}
				config.DataType = *dataType
				config.TypePrefix = tpy.ShowPrefix()
				config.Params["typeprefix"] = tpy.ShowPrefix()
				found = true
				break
			}
		}
		if !found {
			if created, err := bmnumbers.EventuallyCreateType(*dataType, nil); err == nil {
				if created {
					for _, tpy := range bmnumbers.AllTypes {
						if tpy.GetName() == *dataType {
							for opType, opName := range tpy.ShowInstructions() {
								config.Params[opType] = opName
							}
							config.DataType = *dataType
							config.TypePrefix = tpy.ShowPrefix()
							config.Params["typeprefix"] = tpy.ShowPrefix()
							break
						}
					}
				} else {
					panic("Unknown data type")
				}

			} else {
				panic(err)
			}
		}
	} else {
		if config.DataType == "" {
			panic("No data type specified")
		}
	}

	mo := new(bmmatrix.MatrixOpertions)
	mo.RegisterSize = *registerSize

	switch *iomode {
	case "async":
		mo.IOMode = bmmatrix.ASYNC
	case "sync":
		mo.IOMode = bmmatrix.SYNC
	default:
		panic("Unknown IO mode")
	}

	if *saveBasm != "" {
		if basmFile, err := mo.WriteBasm(); err == nil {
			os.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

	if *bmInfoFile != "" {
		// Write the info file
		if bmInfoFileJSON, err := json.MarshalIndent(config.BMinfo, "", "  "); err == nil {
			os.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
		} else {
			panic(err)
		}
	}
}
