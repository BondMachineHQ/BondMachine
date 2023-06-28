package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BondMachineHQ/BondMachine/pkg/basm"
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

// BondMachine targets

var bondmachineFile = flag.String("bondmachine", "", "Load a bondmachine JSON file")
var outFile = flag.String("o", "", "Output file")
var target = flag.String("target", "bondmachine", "Choose the assembler target among: bondmachine, bcof (BondMachineClusteredObjectFormat) ")

// Utils
var getMeta = flag.String("getmeta", "", "Get the metadata of an internal parameter of the BondMachine")

// Optionals
var bmInfoFile = flag.String("bminfo-file", "", "Load additional information about the BondMachine")
var dumpRequirements = flag.String("dump-requirements", "", "Dump the requirements of the BondMachine in a JSON file")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

// TODO var passes = flag.String("passes", "", "List of comma separated optional passes to run (default: none)")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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
	var bm *bondmachine.Bondmachine

	if *bondmachineFile != "" {
		if _, err := os.Stat(*bondmachineFile); err == nil {
			// Open the bondmachine file is exists
			if bondmachineJSON, err := ioutil.ReadFile(*bondmachineFile); err == nil {
				var bmj bondmachine.Bondmachine_json
				if err := json.Unmarshal([]byte(bondmachineJSON), &bmj); err == nil {
					bm = (&bmj).Dejsoner()
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
			bm.Init()
		} else {
			bm = nil
		}
	}

	bi := new(basm.BasmInstance)

	if *debug {
		bi.SetDebug()
	}

	if *verbose {
		bi.SetVerbose()
	}

	if *bmInfoFile != "" {
		bi.BMinfo = new(bminfo.BMinfo)
		if bmInfoJSON, err := ioutil.ReadFile(*bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, bi.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	bi.BasmInstanceInit(bm)

	for _, asmFile := range flag.Args() {
		// Get the file extension
		ext := filepath.Ext(asmFile)

		switch ext {

		case ".basm":
			err := bi.ParseAssemblyDefault(asmFile)
			if err != nil {
				bi.Alert("Error while parsing file:", err)
				return
			}
		case ".ll":
			err := bi.ParseAssemblyLLVM(asmFile)
			if err != nil {
				bi.Alert("Error while parsing file:", err)
				return
			}
		default:
			// Default is .basm
			err := bi.ParseAssemblyDefault(asmFile)
			if err != nil {
				bi.Alert("Error while parsing file:", err)
				return
			}
		}
	}

	if err := bi.RunAssembler(); err != nil {
		bi.Alert(err)
		return
	}

	// All the utils

	if *getMeta != "" {
		if meta, err := bi.GetMeta(*getMeta); err == nil {
			fmt.Println(meta)
		} else {
			bi.Alert(err)
		}
		return
	}

	// Targets

	switch *target {
	case "bondmachine":
		if err := bi.Assembler2BondMachine(); err != nil {
			bi.Alert("Error in creating a BondMachine", err)
			return
		}

		var outF string
		if *outFile != "" {
			outF = *outFile
		} else {
			outF = "bondmachine.json"
		}

		bMach := bi.GetBondMachine()

		// Write the bondmachine file (TODO rewrite)
		f, _ := os.Create(outF)
		defer f.Close()
		b, _ := json.Marshal(bMach.Jsoner())
		f.WriteString(string(b))

		if *bmInfoFile != "" {
			// Write the config file
			if bmInfoFileJSON, err := json.MarshalIndent(bi.BMinfo, "", "  "); err == nil {
				ioutil.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
			} else {
				panic(err)
			}
		}

		if *dumpRequirements != "" {
			// Write the requirements file
			if requirementsJSON, err := json.MarshalIndent(bi.DumpRequirements(), "", "  "); err == nil {
				ioutil.WriteFile(*dumpRequirements, requirementsJSON, 0644)
			} else {
				panic(err)
			}
		}

	case "bcof":
		if err := bi.Assembler2BCOF(); err != nil {
			bi.Alert("Error in creating a BCOF file", err)
			return
		}
	default:
		bi.Alert("Unknown assembler target")
		return
	}
}
