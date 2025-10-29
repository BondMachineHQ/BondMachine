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
	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"google.golang.org/protobuf/proto"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

// BondMachine targets

var bondmachineFile = flag.String("bondmachine", "", "Load a bondmachine JSON file")
var symInFile = flag.String("si", "", "Load a symbols JSON file")
var bmOutFile = flag.String("o", "", "BondMachine Output file")
var bcofOutFile = flag.String("bo", "", "BCOF Output file")
var symOutFile = flag.String("so", "", "Symbols Output file")
var clusOutFile = flag.String("co", "", "Cluster Output file")
var basmOutPrefix = flag.String("oprefix", "", "Prefix for the output files of the assembler (default: none)")

// Utils
var getMeta = flag.String("getmeta", "", "Get the metadata of an internal parameter of the BondMachine")

// Optionals
var bmInfoFile = flag.String("bminfo-file", "", "Load additional information about the BondMachine")
var dumpRequirements = flag.String("dump-requirements", "", "Dump the requirements of the BondMachine in a JSON file")
var createMapFile = flag.String("create-mapfile", "", "Create a mapping file for the BondMachine I/O")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

// Phases
var listPasses = flag.Bool("list-passes", false, "List the available passes")
var actPasses = flag.String("activate-passes", "", "List of comma separated optional passes to activate (default: none)")
var deactPasses = flag.String("deactivate-passes", "", "List of comma separated optional passes to deactivate (default: none)")

// Optimizations
var listOptimizations = flag.Bool("list-optimizations", false, "List the available optimizations")
var actOptimizations = flag.String("activate-optimizations", "", "List of comma separated optional optimizations to activate (default: none, everything: all)")

// Config
var disableDynamicalMatching = flag.Bool("disable-dynamical-matching", false, "Disable the dynamical matching")
var chooserMinWordSize = flag.Bool("chooser-min-word-size", false, "Choose the minimum word size for the chooser")
var chooserForceSameName = flag.Bool("chooser-force-same-name", false, "Force the chooser to use the same name for the ROM and the RAM")

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

	// The BMinfo is always initialized and used internally. If the user wants to use it, it must be loaded from a file. Specifying the file is optional and it will be save only if specified.
	bi.BMinfo = new(bminfo.BMinfo)

	if *bmInfoFile != "" {
		if bmInfoJSON, err := os.ReadFile(*bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, bi.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	bi.BasmInstanceInit(bm)

	// Config options
	if *disableDynamicalMatching {
		bi.Activate(bmconfig.DisableDynamicalMatching)
	}
	if *chooserMinWordSize {
		bi.Activate(bmconfig.ChooserMinWordSize)
	}
	if *chooserForceSameName {
		bi.Activate(bmconfig.ChooserForceSameName)
	}

	mne := basm.GetPassMnemonic()

	if *actPasses != "" {
		passA := strings.Split(*actPasses, ",")

		for _, p := range passA {
			if err := bi.SetActive(p, "basm"); err != nil {
				bi.Alert("Error while activating pass:", err)
				return
			}
		}
	}

	if *deactPasses != "" {
		passD := strings.Split(*deactPasses, ",")

		for _, p := range passD {
			if err := bi.UnsetActive(p, "basm"); err != nil {
				bi.Alert("Error while deactivating pass:", err)
				return
			}
		}
	}

	if *listPasses || *debug || *verbose {
		fmt.Println("Passes:")
		pass := uint64(1)

		for i := 1; pass != basm.LAST_PASS; i++ {
			opt := basm.IsOptionalPass("basm")[pass]
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

	if !bi.IsClustered() && *bmOutFile != "" {
		if err := bi.Assembler2BondMachine(); err != nil {
			bi.Alert("Error in creating a BondMachine", err)
			return
		}

		var outF string
		if *bmOutFile != "" {
			outF = *bmOutFile
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
				os.WriteFile(*bmInfoFile, bmInfoFileJSON, 0644)
			} else {
				panic(err)
			}
		}

		if *dumpRequirements != "" {
			// Write the requirements file
			if requirementsJSON, err := json.MarshalIndent(bi.DumpRequirements(), "", "  "); err == nil {
				os.WriteFile(*dumpRequirements, requirementsJSON, 0644)
			} else {
				panic(err)
			}
		}

		if *createMapFile != "" {
			// Create the mapping file
			if err := bi.CreateMappingFile(*createMapFile); err != nil {
				panic(fmt.Sprintf("Error in creating a mapping file: %v", err))
			}
		}
	}

	if !bi.IsClustered() && *bcofOutFile != "" {
		if err := bi.Assembler2BCOF(); err != nil {
			panic("Error in creating a BCOF file")
		} else {
			bcofBytes, err := proto.Marshal(bi.GetBCOF())
			if err != nil {
				panic("failed to marshal BCOF")
			}
			if err := os.WriteFile(*bcofOutFile, bcofBytes, 0644); err != nil {
				panic("failed to write BCOF file")
			}
		}
	}

	if bi.IsClustered() && *clusOutFile != "" && *basmOutPrefix != "" {
		if err := bi.Assembler2Cluster(); err != nil {
			panic("Error in creating a Cluster")
		} else {

			// Write the cluster file
			clusterBytes, err := json.Marshal(bi.GetCluster())
			if err != nil {
				panic("failed to marshal Cluster")
			}
			if err := os.WriteFile(*clusOutFile, clusterBytes, 0644); err != nil {
				panic("failed to write Cluster file")
			}

			for bmName, bmId := range bi.GetClusteredName() {
				if *debug || *verbose {
					fmt.Printf("Writing BondMachine %s to %s%d.bmeta\n", bmName, *basmOutPrefix, bmId)
				}

				edgeFile := fmt.Sprintf("%s%d.bmeta", *basmOutPrefix, bmId)
				if err := os.WriteFile(edgeFile, []byte(bi.GetClusteredBondMachines()[bmId]), 0644); err != nil {
					panic(fmt.Sprintf("failed to write BondMachine file %s: %v", edgeFile, err))
				}

				if *debug || *verbose {
					fmt.Printf("Writing BondMachine Maps %s to %s%d_maps.json\n", bmName, *basmOutPrefix, bmId)
				}

				mapFile := fmt.Sprintf("%s%d_maps.json", *basmOutPrefix, bmId)
				assoc := bi.GetClusteredMaps()[bmId]
				if mapBytes, err := json.Marshal(assoc); err == nil {
					os.WriteFile(mapFile, mapBytes, 0644)
				} else {
					panic(fmt.Sprintf("failed to marshal BondMachine Maps file %s: %v", mapFile, err))
				}
			}
		}
	}
}
