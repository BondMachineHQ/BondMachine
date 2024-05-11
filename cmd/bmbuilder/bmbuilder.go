package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BondMachineHQ/BondMachine/pkg/bmbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

// Build modes

var bmFile = flag.String("save-bondmachine", "a.out.json", "Save the BM file")

var emitBMAPIMaps = flag.Bool("emit-bmapi-maps", false, "Emit the BMAPIMaps")
var bmAPIMapsFile = flag.String("bmapi-maps-file", "bmapi.json", "BMAPIMaps file to be used as output")

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

	startBuilding := false

	for _, bmbFile := range flag.Args() {

		startBuilding = true

		// Get the file extension
		ext := filepath.Ext(bmbFile)

		switch ext {

		case ".bmb":
			err := bld.ParseBuilderDefault(bmbFile)
			if err != nil {
				bld.Alert("Error while parsing file:", err)
				return
			}
		default:
			err := bld.ParseBuilderDefault(bmbFile)
			if err != nil {
				bld.Alert("Error while parsing file:", err)
				return
			}
		}
	}

	if startBuilding {

		if err := bld.RunBuilder(); err != nil {
			bld.Alert(err)
			return
		}

		if *debug {
			fmt.Println(purple("BmBuilder completed"))
		}
	}

	if *bmFile != "" {
		if err := bld.BuildBondMachine(); err != nil {
			bld.Alert("Error in creating a BondMachine", err)
			return
		}

		var outF string
		if *bmFile != "" {
			outF = *bmFile
		} else {
			outF = "bondmachine.json"
		}

		bMach := bld.GetBondMachine()

		// Write the bondmachine file (TODO rewrite)
		f, _ := os.Create(outF)
		defer f.Close()
		b, _ := json.Marshal(bMach.Jsoner())
		f.WriteString(string(b))

		// if *bmInfoFile != "" {
		// 	// Write the config file
		// 	if bmInfoFileJSON, err := json.MarshalIndent(bi.BMinfo, "", "  "); err == nil {
		// 		os.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
		// 	} else {
		// 		panic(err)
		// 	}
		// }

		// if *dumpRequirements != "" {
		// 	// Write the requirements file
		// 	if requirementsJSON, err := json.MarshalIndent(bi.DumpRequirements(), "", "  "); err == nil {
		// 		os.WriteFile(*dumpRequirements, requirementsJSON, 0644)
		// 	} else {
		// 		panic(err)
		// 	}
		// }
	}

}
