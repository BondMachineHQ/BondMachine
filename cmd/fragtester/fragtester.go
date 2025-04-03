package main

import (
	"flag"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/fragtester"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Debug")

var registerSize = flag.Int("register-size", 32, "Number of bits per register (n-bit)")
var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")
var fragmentFile = flag.String("fragment-file", "", "Name of the fragment file")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
}

func main() {

	config := fragtester.NewConfig()

	if *neuronLibPath != "" {
		config.NeuronLibPath = *neuronLibPath
	} else {
		panic("No neuron library path specified")
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

	if *fragmentFile != "" {
		// Load the fragment file
		if file, err := os.ReadFile(*neuronLibPath + "/" + *fragmentFile); err != nil {
			panic(err)
		} else {
			config.AnalyzeFragment(string(file))
		}
	} else {
		panic("No fragment file specified")
	}

	if *saveBasm != "" {
		if basmFile, err := config.WriteBasm(); err == nil {
			os.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

}
