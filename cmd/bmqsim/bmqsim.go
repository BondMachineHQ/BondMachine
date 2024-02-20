package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

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

// Build modes

// 1
var buildMatrixSeqHardcoded = flag.String("build-matrix-seq-hardcoded", "", "Build a matrix sequence BM with hardcoded quantum circuit")

// 2
var buildMatrixSeq = flag.String("build-matrix-seq", "", "Build a matrix sequence BM with a loadable quantum circuit file")
var buildMatrixSeqCompiled = flag.String("build-matrix-seq-compiled", "", "Build a binary for a matrix sequence BM")

// 3
var buildFullHardwareHardcoded = flag.String("build-full-hw-hardcoded", "", "Build a full hardware BM with hardcoded quantum circuit")

var hardwareFlavor = flag.String("hw-flavor", "", "Hardware flavor for the selected operating mode")
var hardwareFlavorList = flag.Bool("hw-flavor-list", false, "List of available hardware flavors")

// Other options
var showMatrices = flag.Bool("show-matrices", false, "Show the matrices")
var showCircuitMatrix = flag.Bool("show-circuit-matrix", false, "Show the circuit matrix")

func init() {
	flag.Parse()

	// if *debug {
	// 	fmt.Println("basm init")
	// }

	numOp := 0
	if *buildFullHardwareHardcoded != "" {
		numOp++
	}
	if *buildMatrixSeqHardcoded != "" {
		numOp++
	}
	if *buildMatrixSeq != "" {
		numOp++
	}
	if numOp > 1 {
		log.Fatal("Only one build mode can be selected among: build-full-hw-hardcoded, build-matrix-seq, build-matrix-seq-hardcoded")
	}

	if *buildMatrixSeqCompiled != "" && numOp == 1 && *buildMatrixSeq == "" {
		log.Fatal("A loadable quantum circuit file must be used alone or in combination with build-matrix-seq option")
	}

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

	startBuilding := false

	for _, bmqFile := range flag.Args() {

		startBuilding = true

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

	if !startBuilding {
		return
	}

	if *buildFullHardwareHardcoded != "" {
		// Build a full hardware BM with hardcoded quantum circuit this is a special case incompatible with the rest of the modes
		// Matrices won't be generated, the hardware will be built directly

		// Run the builder with the full set of passes
		if err := bld.RunBuilder(); err != nil {
			bld.Alert(err)
			return
		}

		// TODO: Finish this
		fmt.Println("Under construction")
	} else {
		// All the other modes run the builder with a minimal set of passes only to parse the quantum circuit
		// and generate the matrices

		bld.UnsetActive("generatorsexec")

		if err := bld.RunBuilder(); err != nil {
			bld.Alert(err)
			return
		}

		if *debug {
			fmt.Println(purple("BmBuilder completed"))
		}

		var body *bmline.BasmBody

		if *debug {
			fmt.Println(purple("Exporting circuit"))
		}

		// Export the BasmBody to generate the matrices
		if v, err := bld.ExportBasmBody(); err != nil {
			bld.Alert(err)
			return
		} else {
			body = v
		}

		if *debug {
			fmt.Println(purple("Processing circuit to matrices"))
		}
		// Get the circuit matrices from the BasmBody
		if matrices, err := sim.QasmToBmMatrices(body); err != nil {
			bld.Alert(err)
			return
		} else {
			sim.Mtx = make([]*bmmatrix.BmMatrixSquareComplex, len(matrices))
			copy(sim.Mtx, matrices)
		}
	}

	if *showMatrices {
		if sim.Mtx == nil {
			bld.Alert("No matrices to show")
			return
		} else {
			for i, m := range sim.Mtx {
				fmt.Println(green("Matrix:"), yellow(strconv.Itoa(i)))
				fmt.Println(m.StringColor(green))
			}
		}
	}

	if *showCircuitMatrix {
		mm := sim.Mtx[len(sim.Mtx)-1]
		for i := len(sim.Mtx) - 2; i >= 0; i-- {
			mm = bmmatrix.MatrixProductComplex(mm, sim.Mtx[i])
		}
		fmt.Println(green("Whole circuit matrix:"))
		fmt.Println(mm.StringColor(green))
	}

	if *buildMatrixSeqHardcoded != "" {
		// Build a matrix sequence BM with hardcoded quantum circuit
		fmt.Println("Under construction")

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
		} else if *hardwareFlavor != "" {
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}

	if *buildMatrixSeq != "" {
		// Build a matrix sequence BM with a loadable quantum circuit file
		fmt.Println("Under construction")

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
		} else if *hardwareFlavor != "" {
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}

	if *buildMatrixSeqCompiled != "" {
		// Build a binary for a matrix sequence BM
		fmt.Println("Under construction")

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
		} else if *hardwareFlavor != "" {
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}
}
