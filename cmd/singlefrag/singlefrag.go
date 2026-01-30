package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/fragtester"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Debug")

var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var describe = flag.Bool("describe", false, "Describe the fragment without running it")

var saveBasm = flag.String("save-basm", "", "Create a basm file")

var createBmApi = flag.String("create-bmapi", "", "Create a mapping file for the BondMachine I/O over BMAPI")

var createSicv2Endpoints = flag.String("sicv2-endpoints", "", "Create SICv2 endpoints and activate the benchcore module")

var buildApp = flag.Bool("build-app", false, "Build an hardware connected app")
var appFlavor = flag.String("app-flavor", "", "App flavor for the selected operating mode")
var appFlavorList = flag.Bool("app-flavor-list", false, "List of available app flavors")
var appFile = flag.String("app-file", "a.out", "App file to be used as output")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")
var fragmentFile = flag.String("fragment-file", "", "Name of the fragment file")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
}

func main() {

	ft := fragtester.NewFragTester()

	if *verbose {
		ft.Verbose = true
	}
	if *debug {
		ft.Debug = true
	}

	// Requests that does not need fragment analysis

	if *appFlavorList {
		flavors := ft.ListAppFlavors()
		for _, flavor := range flavors {
			fmt.Println(flavor)
		}
		os.Exit(0)
	}

	// Set the data type
	if *dataType != "" {
		found := false
		for _, tpy := range bmnumbers.AllTypes {
			if tpy.GetName() == *dataType {
				for opType, opName := range tpy.ShowInstructions() {
					ft.Params[opType] = opName
					ft.OpString += ", " + opType + ":" + opName
				}
				ft.DataType = *dataType
				ft.TypePrefix = tpy.ShowPrefix()
				ft.OpString += ", prefix:" + tpy.ShowPrefix()
				ft.Params["typeprefix"] = tpy.ShowPrefix()
				ft.RegisterSize = tpy.GetSize()
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
								ft.Params[opType] = opName
								ft.OpString += ", " + opType + ":" + opName
							}
							ft.DataType = *dataType
							ft.TypePrefix = tpy.ShowPrefix()
							ft.OpString += ", prefix:" + tpy.ShowPrefix()
							ft.Params["typeprefix"] = tpy.ShowPrefix()
							ft.RegisterSize = tpy.GetSize()
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
		if ft.DataType == "" {
			panic("No data type specified")
		}
	}

	// Set the neuron library path
	if *neuronLibPath != "" {
		ft.NeuronLibPath = *neuronLibPath
	} else {
		panic("No neuron library path specified")
	}

	// Load the fragment file and analyze it
	if *fragmentFile != "" {
		// Load the fragment file
		if file, err := os.ReadFile(*neuronLibPath + "/" + *fragmentFile); err != nil {
			panic(err)
		} else {
			ft.AnalyzeFragment(string(file))
		}
	} else {
		panic("No fragment file specified")
	}

	// Check if the fragment is valid
	if !ft.Valid {
		os.Exit(1)
	}

	// Requests that need fragment analysis and are mutually exclusive

	// Describe the fragment
	if *describe {
		ft.DescribeFragment()

	} else {
		// Run the fragment and produce outputs

		// Activate benchcoreV2 if requested
		if *createSicv2Endpoints != "" {
			if sicv2Ends, err := ft.WriteSicv2Endpoints(); err == nil {
				os.WriteFile(*createSicv2Endpoints, []byte(sicv2Ends), 0644)
				ft.BenchcoreV2 = true
			} else {
				panic(err)
			}
		}

		// Create the BMAPI mapping file if requested
		if *createBmApi != "" {
			if err := ft.CreateMappingFile(*createBmApi); err != nil {
				panic(err)
			}
		}

		// Build the app if requested
		if *buildApp && *appFile != "" {
			if appSource, err := ft.WriteApp(*appFlavor); err == nil {
				os.WriteFile(*appFile, []byte(appSource), 0644)
			} else {
				panic(err)
			}
		}

		// Produce the requested output files
		if *saveBasm != "" {
			if basmFile, err := ft.WriteBasm(); err == nil {
				os.WriteFile(*saveBasm, []byte(basmFile), 0644)
			} else {
				panic(err)
			}
		}
	}
}
