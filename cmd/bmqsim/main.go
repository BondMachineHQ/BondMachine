package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/BondMachineHQ/BondMachine/pkg/bmbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

func init() {
	flag.Parse()

	// if *debug {
	// 	fmt.Println("basm init")
	// }

	if *linearDataRange != "" {
		if err := bmnumbers.LoadLinearDataRangesFromFile(*linearDataRange); err != nil {
			log.Fatal(err)
		}

		var lqRanges *map[int]bmnumbers.LinearDataRange
		for _, t := range bmnumbers.AllDynamicalTypes {
			if t.GetName() == "dyn_linear_quantizer" {
				lqRanges = t.(bmnumbers.DynLinearQuantizer).Ranges
			}
		}

		for i, t := range procbuilder.AllDynamicalInstructions {
			if t.GetName() == "dyn_linear_quantizer" {
				dynIst := t.(procbuilder.DynLinearQuantizer)
				dynIst.Ranges = lqRanges
				procbuilder.AllDynamicalInstructions[i] = dynIst
			}
		}
	}
}

func main() {
	bld := new(bmbuilder.BMBuilder)

	bld.BMBuilderInit()

	if *debug {
		bld.SetDebug()
	}

	if *verbose {
		bld.SetVerbose()
	}

	startAssembling := false

	for _, bmqFile := range flag.Args() {

		startAssembling = true

		// Get the file extension
		ext := filepath.Ext(bmqFile)

		switch ext {

		case ".bmq":
			err := bld.ParseBuilderDefault(bmqFile)
			if err != nil {
				bld.Alert("Error while parsing file:", err)
				return
			}
		default:
			err := bld.ParseBuilderDefault(bmqFile)
			if err != nil {
				bld.Alert("Error while parsing file:", err)
				return
			}
		}
	}

	if !startAssembling {
		return
	}

	if err := bld.RunBuilder(); err != nil {
		bld.Alert(err)
		return
	}

}
