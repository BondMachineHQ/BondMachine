package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/melbond"
	"github.com/BondMachineHQ/BondMachine/pkg/neuralbond"
	"github.com/mmirko/mel/pkg/m3number"
	"github.com/mmirko/mel/pkg/mel"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 32, "Number of bits per register (n-bit)")
var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")

var configFile = flag.String("config-file", "", "JSON description of the net configuration")
var bmInfoFile = flag.String("bminfo-file", "", "JSON description of the BondMachine abstraction")

var iomode = flag.String("io-mode", "async", "IO mode: async, sync")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
}

func main() {
	bc := new(melbond.MelBondConfig)
	bc.RegisterSize = uint8(*registerSize)

	// Load net from a JSON file the configuration
	if *configFile != "" {
		if netFileJSON, err := ioutil.ReadFile(*configFile); err == nil {
			if err := json.Unmarshal(netFileJSON, bc); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		bc.Debug = *debug
		bc.Verbose = *verbose
		bc.Params = make(map[string]string)
	}

	// Load or create the Info file
	bc.BMinfo = new(bminfo.BMinfo)

	if *bmInfoFile != "" {
		if bmInfoJSON, err := ioutil.ReadFile(*bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, bc.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	if bc.Params == nil {
		bc.Params = make(map[string]string)
	}
	if bc.List == nil {
		bc.List = make(map[string]string)
	}
	if bc.Pruned == nil {
		bc.Pruned = make([]string, 0)
	}

	if *neuronLibPath != "" {
		bc.NeuronLibPath = *neuronLibPath
	} else {
		panic("No neuron library path specified")
	}

	switch *iomode {
	case "async":
		bc.IOMode = neuralbond.ASYNC
	case "sync":
		bc.IOMode = neuralbond.SYNC
	default:
		panic("Unknown IO mode")
	}

	if *dataType != "" {
		found := false
		for _, tpy := range bmnumbers.AllTypes {
			if tpy.GetName() == *dataType {
				for opType, opName := range tpy.ShowInstructions() {
					bc.Params[opType] = opName
				}
				bc.DataType = *dataType
				bc.TypePrefix = tpy.ShowPrefix()
				bc.Params["typeprefix"] = tpy.ShowPrefix()
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
								bc.Params[opType] = opName
							}
							bc.DataType = *dataType
							bc.TypePrefix = tpy.ShowPrefix()
							bc.Params["typeprefix"] = tpy.ShowPrefix()
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
		if bc.DataType == "" {
			panic("No data type specified")
		}
	}

	a := new(m3number.M3numberMe3li)
	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	a.Mel3Object.DefaultCreator = bc.BasmCreator
	c.Debug = false
	a.MelInit(c, ep)

	prog := new(melbond.MelBondProgram)
	prog.MelBondConfig = bc
	prog.M3numberMe3li = a

	if len(flag.Args()) != 1 {
		panic("No mel file specified")
	}

	for _, melFile := range flag.Args() {
		if source, err := ioutil.ReadFile(melFile); err != nil {
			panic(err)
		} else {
			prog.Source = string(source)
		}
	}

	if *saveBasm != "" {
		if basmFile, err := prog.WriteBasm(); err == nil {
			ioutil.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

	if *bmInfoFile != "" {
		// Write the info file
		if bmInfoFileJSON, err := json.MarshalIndent(bc.BMinfo, "", "  "); err == nil {
			ioutil.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
		} else {
			panic(err)
		}
	}

	// Remove the info file from the config prior to saving it
	bc.BMinfo = nil
	if *configFile != "" {
		// Write the eventually updated config file
		if configFileJSON, err := json.MarshalIndent(bc, "", "  "); err == nil {
			ioutil.WriteFile(*configFile, configFileJSON, 0644)
		} else {
			panic(err)
		}
	}

}
