package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/neuralbond"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 32, "Number of bits per register (n-bit)")
var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")

var netFile = flag.String("net-file", "", "JSON description of the net")
var configFile = flag.String("config-file", "", "JSON description of the net configuration")
var bmInfoFile = flag.String("bminfo-file", "", "JSON description of the BondMachine abstraction")

var operatingMode = flag.String("operating-mode", "romcode", "Operating mode: romcode, fragment")
var iomode = flag.String("io-mode", "async", "IO mode: async, sync")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
}

func main() {
	net := new(neuralbond.TrainedNet)

	// Load net from a JSON file
	if *netFile != "" {
		if netFileJSON, err := ioutil.ReadFile(*netFile); err == nil {
			if err := json.Unmarshal(netFileJSON, net); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		panic("No net file specified")
	}

	net.RegisterSize = *registerSize

	switch *operatingMode {
	case "romcode":
		net.OperatingMode = neuralbond.ROMCODE
	case "fragment":
		net.OperatingMode = neuralbond.FRAGMENT
	default:
		panic("Unknown operating mode")
	}

	config := new(neuralbond.Config)

	// Load net from a JSON file the configuration
	if *configFile != "" {
		if netFileJSON, err := ioutil.ReadFile(*configFile); err == nil {
			if err := json.Unmarshal(netFileJSON, config); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		config.Debug = *debug
		config.Verbose = *verbose
		config.Params = make(map[string]string)
	}

	// Load or create the Info file
	config.BMinfo = new(bminfo.BMinfo)

	if *bmInfoFile != "" {
		if bmInfoJSON, err := ioutil.ReadFile(*bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, config.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	if config.Params == nil {
		config.Params = make(map[string]string)
	}
	if config.List == nil {
		config.List = make(map[string]string)
	}
	if config.Pruned == nil {
		config.Pruned = make([]string, 0)
	}

	if *neuronLibPath != "" {
		config.NeuronLibPath = *neuronLibPath
	} else {
		panic("No neuron library path specified")
	}

	if err := net.Init(config); err != nil {
		panic(err)
	}

	switch *iomode {
	case "async":
		net.IOMode = neuralbond.ASYNC
	case "sync":
		net.IOMode = neuralbond.SYNC
	default:
		panic("Unknown IO mode")
	}

	net.Normalize()

	if *dataType != "" {
		found := false
		for _, tpy := range bmnumbers.AllTypes {
			if tpy.GetName() == *dataType {
				for opType, opName := range tpy.ShowInstructions() {
					config.Params[opType] = opName
				}
				config.DataType = *dataType
				config.TypePrefix = tpy.ShowPrefix()
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

	if *saveBasm != "" {
		if basmFile, err := net.WriteBasm(); err == nil {
			ioutil.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

	if *bmInfoFile != "" {
		// Write the info file
		if bmInfoFileJSON, err := json.MarshalIndent(config.BMinfo, "", "  "); err == nil {
			ioutil.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
		} else {
			panic(err)
		}
	}

	// Remove the info file from the config prior to saving it
	config.BMinfo = nil
	if *configFile != "" {
		// Write the eventually updated config file
		if configFileJSON, err := json.MarshalIndent(config, "", "  "); err == nil {
			ioutil.WriteFile(*configFile, configFileJSON, 0644)
		} else {
			panic(err)
		}
	}

}
