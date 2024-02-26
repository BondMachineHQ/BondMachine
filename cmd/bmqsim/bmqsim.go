package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
var buildMatrixSeqHardcoded = flag.Bool("build-matrix-seq-hardcoded", false, "Build a matrix sequence BM with hardcoded quantum circuit")

// 2
var buildMatrixSeq = flag.Bool("build-matrix-seq", false, "Build a matrix sequence BM with a loadable quantum circuit file")
var buildMatrixSeqCompiled = flag.Bool("build-matrix-seq-compiled", false, "Build a binary for a matrix sequence BM")

// 3
var buildFullHardwareHardcoded = flag.Bool("build-full-hw-hardcoded", false, "Build a full hardware BM with hardcoded quantum circuit")

// 4
var softwareSimulation = flag.Bool("software-simulation", false, "Software simulation mode")
var softwareSimulationInput = flag.String("software-simulation-input", "", "Software simulation mode input file")
var softwareSimulationOutput = flag.String("software-simulation-output", "", "Software simulation mode output file")

// Common options

var buildApp = flag.Bool("build-app", false, "Build an hardware connected app")
var appFlavor = flag.String("app-flavor", "", "App flavor for the selected operating mode")
var appFlavorList = flag.Bool("app-flavor-list", false, "List of available app flavors")
var appFile = flag.String("app-file", "a.out", "App file to be used as output")

var basmFile = flag.String("save-basm", "a.out.basm", "Basm file to be used as output")
var compiledFile = flag.String("compiled-file", "a.out.json", "Compiled file to be used as output")

var hardwareFlavor = flag.String("hw-flavor", "", "Hardware flavor for the selected operating mode")
var hardwareFlavorList = flag.Bool("hw-flavor-list", false, "List of available hardware flavors")

// Other options
var showMatrices = flag.Bool("show-matrices", false, "Show the matrices")
var showCircuitMatrix = flag.Bool("show-circuit-matrix", false, "Show the circuit matrix")

var emitBMAPIMaps = flag.Bool("emit-bmapi-maps", false, "Emit the BMAPIMaps")
var bmAPIMapsFile = flag.String("bmapi-maps-file", "bmapi.json", "BMAPIMaps file to be used as output")

func init() {
	flag.Parse()

	// if *debug {
	// 	fmt.Println("basm init")
	// }

	numOp := 0
	if *buildFullHardwareHardcoded {
		numOp++
	}
	if *buildMatrixSeqHardcoded {
		numOp++
	}
	if *buildMatrixSeq || *buildMatrixSeqCompiled {
		numOp++
	}
	if *softwareSimulation {
		numOp++
	}
	if numOp == 0 {
		log.Fatal("No build mode selected")
	}
	if numOp > 1 {
		log.Fatal("Only one build mode can be selected among: build-full-hw-hardcoded, build-matrix-seq(_compiled), build-matrix-seq-hardcoded, software-simulation")
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

	if startBuilding {

		if *buildFullHardwareHardcoded {
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

		if *emitBMAPIMaps {
			if fileData, err := sim.EmitBMAPIMaps(); err != nil {
				bld.Alert(err)
				return
			} else {
				os.WriteFile(*bmAPIMapsFile, []byte(fileData), 0644)
			}
		}

	}

	if *buildMatrixSeqHardcoded {
		// Build a matrix sequence BM with hardcoded quantum circuit

		modeTags := []string{"real"}

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
			for t := range bmqsim.HardwareFlavors {
				flavorTags := bmqsim.HardwareFlavorsTags[t]
				for _, tag := range modeTags {
					if bmqsim.StringInSlice(tag, flavorTags) {
						fmt.Println(t)
					}
				}
			}
		} else if *hardwareFlavor != "" {
			if startBuilding {
				if _, ok := bmqsim.HardwareFlavors[*hardwareFlavor]; ok {
					if sim.VerifyConditions(*hardwareFlavor) == nil {
						if basmFileData, err := sim.ApplyTemplate(*hardwareFlavor); err != nil {
							bld.Alert(err)
						} else {
							os.WriteFile(*basmFile, []byte(basmFileData), 0644)
						}
					} else {
						bld.Alert("Hardware flavor not compatible with the quantum circuit")
					}
				} else {
					bld.Alert("Hardware flavor not found")
				}
			} else {
				bld.Alert("No quantum circuit to build the matrix sequence")
			}
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}

	if *buildMatrixSeq || *buildMatrixSeqCompiled {
		// Build a matrix sequence BM with a loadable quantum circuit file
		fmt.Println("Under construction")

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
		} else if *hardwareFlavor != "" {
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}

	if *softwareSimulation {
		// Build a binary for a matrix sequence BM
		fmt.Println("Under construction")

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
		} else if *hardwareFlavor != "" {
		} else {
			bld.Alert("Hardware flavor must be selected")
		}
	}

	if *buildApp {
		// Build an hardware connected app

		modeTags := []string{"real"}

		if *appFlavorList {
			// List of available app flavors for the selected operating mode
			for t := range bmqsim.AppFlavors {
				flavorTags := bmqsim.AppFlavorsTags[t]
				for _, tag := range modeTags {
					if bmqsim.StringInSlice(tag, flavorTags) {
						fmt.Println(t)
					}
				}
			}
		} else if *appFlavor != "" {
			if startBuilding {
				if _, ok := bmqsim.AppFlavors[*appFlavor]; ok {
					if appFileData, err := sim.ApplyTemplate(*appFlavor); err != nil {
						bld.Alert(err)
					} else {
						os.WriteFile(*appFile, []byte(appFileData), 0644)
					}
				} else {
					bld.Alert("App flavor not found")
				}
			} else {
				bld.Alert("No quantum circuit to build the matrix sequence")
			}
		} else {
			bld.Alert("App flavor must be selected")
		}
	}
}
