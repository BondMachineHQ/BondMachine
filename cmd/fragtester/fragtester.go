package main

import (
	"flag"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/fragtester"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Debug")

var dataType = flag.String("data-type", "float32", "bmnumbers data types")

var saveBasm = flag.String("save-basm", "", "Create a basm file")
var saveExpression = flag.String("save-expression", "", "Create a expression file")

var neuronLibPath = flag.String("neuron-lib-path", "", "Path to the neuron library to use")
var fragmentFile = flag.String("fragment-file", "", "Name of the fragment file")

func init() {
	flag.Parse()
	if *saveBasm == "" {
		*saveBasm = "out.basm"
	}
	if *saveExpression == "" {
		*saveExpression = "expression.py"
	}
}

func main() {

	ft := fragtester.NewFragTester()

	if *neuronLibPath != "" {
		ft.NeuronLibPath = *neuronLibPath
	} else {
		panic("No neuron library path specified")
	}

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

	if !ft.Valid {
		os.Exit(1)
	}

	if *saveBasm != "" {
		if basmFile, err := ft.WriteBasm(); err == nil {
			os.WriteFile(*saveBasm, []byte(basmFile), 0644)
		} else {
			panic(err)
		}
	}

	if *saveExpression != "" {
		if expressionFile, err := ft.WriteSympy(); err == nil {
			os.WriteFile(*saveExpression, []byte(expressionFile), 0644)
		} else {
			panic(err)
		}
	}

}
