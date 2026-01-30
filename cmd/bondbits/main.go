package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/basm"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

type string_slice []string

func (i *string_slice) String() string {
	return fmt.Sprint(*i)
}

func (i *string_slice) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*i = append(*i, dt)
	}
	return nil
}

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

// BondMachine targets
var bondmachineFile = flag.String("bondmachine", "", "Load a bondmachine JSON file")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

// Phases
var listPasses = flag.Bool("list-passes", false, "List the available passes")
var actPasses = flag.String("activate-passes", "", "List of comma separated optional passes to activate (default: none)")
var deactPasses = flag.String("deactivate-passes", "", "List of comma separated optional passes to deactivate (default: none)")

// Optimizations
var listOptimizations = flag.Bool("list-optimizations", false, "List the available optimizations")
var actOptimizations = flag.String("activate-optimizations", "", "List of comma separated optional optimizations to activate (default: none, everything: all)")

// Bondbits rules
var bondBitsRules string_slice
var saveDirectory = flag.String("save-directory", ".", "Directory where to save the generated files")

func init() {
	flag.Var(&bondBitsRules, "bbr", "Add a bondbits rule (can be repeated)")
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
			if bondmachineJSON, err := os.ReadFile(*bondmachineFile); err == nil {
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

	bi.BasmInstanceInit(bm)

	mne := basm.GetPassMnemonic()

	for _, pass := range mne {
		bi.UnsetActive(pass, "bondbits")
	}

	if *actPasses != "" {
		passA := strings.Split(*actPasses, ",")

		for _, p := range passA {
			if err := bi.SetActive(p, "bondbits"); err != nil {
				bi.Alert("Error while activating pass:", err)
				return
			}
		}
	}

	if *deactPasses != "" {
		passD := strings.Split(*deactPasses, ",")

		for _, p := range passD {
			if err := bi.UnsetActive(p, "bondbits"); err != nil {
				bi.Alert("Error while deactivating pass:", err)
				return
			}
		}
	}

	if *listPasses || *debug || *verbose {
		fmt.Println("Passes:")
		pass := uint64(1)

		for i := 1; pass != basm.LAST_PASS; i++ {
			opt := basm.IsOptionalPass("bondbits")[pass]
			if bi.ActivePass(pass) {
				if opt {
					fmt.Printf("  %02d: %s (optional)\n", i, mne[pass])
				} else {
					fmt.Printf("  %02d: %s\n", i, mne[pass])
				}
			} else {
				fmt.Printf("  %02d: %s (optional, disabled)\n", i, mne[pass])
			}
			pass = pass << 1
		}

	}

	if *actOptimizations != "" {
		optA := strings.Split(*actOptimizations, ",")
		for _, o := range optA {
			if err := bi.ActivateOptimization(o); err != nil {
				bi.Alert("Error while activating optimization:", err)
				return
			}
		}
	}

	if *listOptimizations || *debug || *verbose {
		// TODO: Finish this
	}

	startAssembling := false

	for _, asmFile := range flag.Args() {
		startAssembling = true

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

	if !startAssembling {
		return
	}

	if err := bi.RunAssembler(); err != nil {
		bi.Alert(err)
		return
	}

	// TODO Rules may need of specific passes or optimizations activated
	// Rules have to be validated against each other. (some rules may be incompatible)
	// Plus, rules may be given to the CLI or as metadata inside the assembly itself
	// Examples: annotate for a specific call type.
	// TODO Create functions to export basminstance to BASM files

	// Start processing bondbits rules
	for _, rule := range bondBitsRules {
		if *debug || *verbose {
			fmt.Println("Applying bondbits rule:", rule)
		}
		if err := bi.ApplyBondBitsRule(rule); err != nil {
			bi.Alert("Error while applying bondbits rule:", err)
			return
		}
	}

	if *saveDirectory != "" {
		// Create the save directory if it does not exist, if it exists, exit with error
		if _, err := os.Stat(*saveDirectory); err == nil {
			bi.Alert("Error: save directory already exists:", *saveDirectory)
			return
		}
		if err := os.MkdirAll(*saveDirectory, os.ModePerm); err != nil {
			bi.Alert("Error while creating save directory:", err)
			return
		}
		if err := bi.ExportBasmFiles(*saveDirectory); err != nil {
			bi.Alert("Error while exporting BASM files:", err)
			return
		}
	}
}
