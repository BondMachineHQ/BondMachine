package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/neuralbond"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 32, "Number of bits per register (n-bit)")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")

var netFile = flag.String("net-file", "", "JSON description of the net")
var configFile = flag.String("config-file", "", "JSON description of the net configuration")
var bmInfoFile = flag.String("bminfo-file", "", "JSON description of the BondMachine abstraction")

var operatingMode = flag.String("operating-mode", "romcode", "Operating mode: romcode, fragment")

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

	net.Normalize()

	if *saveBasm != "" {
		if basmFile, err := net.WriteBasm(); err == nil {
			ioutil.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

	if *bmInfoFile != "" {
		// Write the config file
		if bmInfoFileJSON, err := json.MarshalIndent(config.BMinfo, "", "  "); err == nil {
			ioutil.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
		} else {
			panic(err)
		}
	}
}
