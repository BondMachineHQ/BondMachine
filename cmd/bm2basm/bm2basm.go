package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/bm2basm"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

// Build modes

var bmFile = flag.String("bondmachine-file", "", "Load the BM file")
var basmFile = flag.String("basm-file", "bondmachine.basm", "save the BASM file")

func init() {
	flag.Parse()

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
	if *bmFile == "" {
		log.Fatal("No source BM file specified")
	}

	if *basmFile == "" {
		log.Fatal("No target BASM file specified")
	}

	if *verbose {
		log.Printf("Loading BM file %s\n", *bmFile)
	}
	var bm *bondmachine.Bondmachine

	if _, err := os.Stat(*bmFile); err == nil {
		// Open the bondmachine file is exists
		if bondmachineJSON, err := os.ReadFile(*bmFile); err == nil {
			var bmj bondmachine.Bondmachine_json
			if err := json.Unmarshal([]byte(bondmachineJSON), &bmj); err == nil {
				bm = (&bmj).Dejsoner()
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
		bm.Init()
	} else {
		log.Fatal(err)
	}

	b := new(bm2basm.Bm2Basm)

	if basmContent, err := b.Convert(bm); err != nil {
		log.Fatal(err)
	} else {

		// Write the BASM file
		if *verbose {
			log.Printf("Saving BASM file %s\n", *basmFile)
		}

		if err := os.WriteFile(*basmFile, []byte(basmContent), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
