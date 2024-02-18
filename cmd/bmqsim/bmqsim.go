package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/BondMachineHQ/BondMachine/pkg/bmbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bmqsim"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

var buildMatrixSeqHardcoded = flag.String("build-matrix-seq-hardcoded", "", "Build a matrix sequence BM with hardcoded quantum circuit")
var buildMatrixSeq = flag.String("build-matrix-seq", "", "Build a matrix sequence BM with a loadable quantum circuit file")
var buildMatrixSeqCompiled = flag.String("build-matrix-seq-compiled", "", "Build a binary for a matrix sequence BM")
var buildFullHardwareHardcoded = flag.String("build-full-hw-hardcoded", "", "Build a full hardware BM with hardcoded quantum circuit")

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
	sim := new(bmqsim.BmQSimulator)

	bld.BMBuilderInit()
	sim.BmQSimulatorInit()

	if *debug {
		bld.SetDebug()
		sim.SetDebug()
	}

	if *verbose {
		bld.SetVerbose()
		sim.SetVerbose()
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

	if *buildFullHardwareHardcoded != "" {
		// Build a full hardware BM with hardcoded quantum circuit this is a special case uncompatible with the rest of the modes

		// Run the builder with the full set of passes
		if err := bld.RunBuilder(); err != nil {
			bld.Alert(err)
			return
		}

		fmt.Println("Under construction")

	} else {

		// Run the builder with a minimal set of passes

		bld.UnsetActive("generatorsexec")

		if err := bld.RunBuilder(); err != nil {
			bld.Alert(err)
			return
		}

		var body *bmline.BasmBody

		// Export the BasmBody to generate the matrices
		if v, err := bld.ExportBasmBody(); err != nil {
			if *buildFullHardwareHardcoded != "" {
				// Build a full hardware BM with hardcoded quantum circuit
				fmt.Println("Under construction")
			}

			bld.Alert(err)
			return
		} else {
			body = v
		}

		fmt.Println("Quantum circuit:")
		fmt.Println(body)

		var mtx []*bmmatrix.BmMatrixSquareComplex

		// Get the circuit matrices from the BasmBody
		if matrices, err := sim.QasmToBmMatrices(body); err != nil {
			bld.Alert(err)
			return
		} else {
			mtx = make([]*bmmatrix.BmMatrixSquareComplex, len(matrices))
			copy(mtx, matrices)
		}

		fmt.Println(mtx)

		if *buildMatrixSeqHardcoded != "" {
			// Build a matrix sequence BM with hardcoded quantum circuit
			fmt.Println("Under construction")
		}

		if *buildMatrixSeq != "" {
			// Build a matrix sequence BM with a loadable quantum circuit file
			fmt.Println("Under construction")
		}

		if *buildMatrixSeqCompiled != "" {
			// Build a binary for a matrix sequence BM
			fmt.Println("Under construction")
		}
	}
}
