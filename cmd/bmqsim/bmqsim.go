package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmatrix"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bmqsim"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

type qbitSwaps []string

func (q *qbitSwaps) String() string {
	return fmt.Sprintf("%v", *q)
}
func (i *qbitSwaps) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var swaps qbitSwaps

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Debug")

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
var softwareSimulationInput = flag.String("software-simulation-input", "", "Software simulation mode input file. If not provided, the input will be zero-state")
var softwareSimulationOutput = flag.String("software-simulation-output", "", "Software simulation mode output file, if not provided the output will be printed to stdout")

// 5
var buildMatrixSeqHLS = flag.Bool("build-matrix-seq-hls", false, "Build a matrix sequence HLS code with hardcoded quantum circuit")

// Common options

var buildApp = flag.Bool("build-app", false, "Build an hardware connected app")
var appFlavor = flag.String("app-flavor", "", "App flavor for the selected operating mode")
var appFlavorList = flag.Bool("app-flavor-list", false, "List of available app flavors")
var appFile = flag.String("app-file", "a.out", "App file to be used as output")

var bmFile = flag.String("save-bondmachine", "bondmachine.json", "Bondmachine file to be used as output")
var basmFile = flag.String("save-basm", "a.out.basm", "Basm file to be used as output")
var compiledFile = flag.String("compiled-file", "a.out.json", "Compiled file to be used as output")
var bundleDir = flag.String("bundle-dir", "", "Bundle directory to be used as output")

var hardwareFlavor = flag.String("hw-flavor", "", "Hardware flavor for the selected operating mode")
var hardwareFlavorList = flag.Bool("hw-flavor-list", false, "List of available hardware flavors")

// Other options
var showMatrices = flag.Bool("show-matrices", false, "Show the matrices")
var showCircuitMatrix = flag.Bool("show-circuit-matrix", false, "Show the circuit matrix")

var emitBMAPIMaps = flag.Bool("emit-bmapi-maps", false, "Emit the BMAPIMaps")
var bmAPIMapsFile = flag.String("bmapi-maps-file", "bmapi.json", "BMAPIMaps file to be used as output")

func init() {
	flag.Var(&swaps, "swap", "Swap qbits in the circuit")
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
	if *buildMatrixSeqHLS {
		numOp++
	}
	if *softwareSimulation {
		numOp++
	}
	if numOp == 0 {
		log.Fatal("No build mode selected")
	}
	if numOp > 1 {
		log.Fatal("Only one build mode can be selected among: build-full-hw-hardcoded, build-matrix-seq(_compiled), build-matrix-seq-hardcoded, software-simulation, build-matrix-seq-hls")
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

			if len(swaps) > 0 {
				for _, swap := range swaps {
					swap1s := strings.Split(swap, ",")[0]
					swap2s := strings.Split(swap, ",")[1]
					swap1, _ := strconv.Atoi(swap1s)
					swap2, _ := strconv.Atoi(swap2s)
					mm = sim.SwapQbits(mm, swap1, swap2)
				}
			}

			fmt.Println(green("Whole circuit matrix:"))
			fmt.Println(mm.StringColor(green))
		}

		if *emitBMAPIMaps {
			if startBuilding && *hardwareFlavor != "" {
				if _, ok := bmqsim.HardwareFlavors[*hardwareFlavor]; ok {
					if fileData, err := sim.EmitBMAPIMaps(*hardwareFlavor); err != nil {
						bld.Alert(err)
						return
					} else {
						os.WriteFile(*bmAPIMapsFile, []byte(fileData), 0644)
					}
				} else {
					bld.Alert("Hardware flavor not found")
				}
			} else {
				bld.Alert("No quantum circuit to build the BM API maps")
			}
		}

	}

	if *buildMatrixSeqHardcoded {
		// Build a matrix sequence BM with hardcoded quantum circuit

		modeTags := []string{"real", "complex"}

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

	if *buildMatrixSeqHLS {
		// Build a matrix sequence HLS code
		modeTags := []string{"real", "complex"}

		if *hardwareFlavorList {
			// List of available hardware flavors for the selected operating mode
			for t := range bmqsim.HLSFlavors {
				flavorTags := bmqsim.HLSFlavorsTags[t]
				for _, tag := range modeTags {
					if bmqsim.StringInSlice(tag, flavorTags) {
						fmt.Println(t)
					}
				}
			}
		} else if *hardwareFlavor != "" {
			// The output of this mode is a bundle directory
			if *bundleDir == "" {
				bld.Alert("Bundle directory must be provided")
				return
			}
			if startBuilding {
				if _, ok := bmqsim.HLSFlavors[*hardwareFlavor]; ok {
					if sim.VerifyConditions(*hardwareFlavor) == nil {
						if err := sim.ApplyTemplateBundle(*hardwareFlavor, *bundleDir); err != nil {
							bld.Alert(err)
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

	if *softwareSimulation {
		if startBuilding {
			// Software simulation mode
			if *softwareSimulationInput != "" {
				// Load the input data from a json file
				inputs := new([]bmqsim.StateArray)
				if inputJSON, err := os.ReadFile(*softwareSimulationInput); err == nil {
					if err := json.Unmarshal([]byte(inputJSON), inputs); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
				sim.Inputs = *inputs
			} else {
				// Zero state
				inputs := new(bmqsim.StateArray)
				inputs.Vector = make([]bmmatrix.Complex32, sim.StateSize())
				for i := range inputs.Vector {
					if i == 0 {
						inputs.Vector[i] = bmmatrix.Complex32{Real: 1, Imag: 0}
					} else {
						inputs.Vector[i] = bmmatrix.Complex32{Real: 0, Imag: 0}
					}
				}
				sim.Inputs = make([]bmqsim.StateArray, 1)
				sim.Inputs[0] = *inputs
			}

			if err := sim.RunSoftwareSimulation(); err != nil {
				bld.Alert(err)
				return
			}

			// Save the output data to a json file or print it to stdout
			if outputJSON, err := json.Marshal(sim.Outputs); err == nil {
				if *softwareSimulationOutput != "" {
					os.WriteFile(*softwareSimulationOutput, outputJSON, 0644)
				} else {
					fmt.Println(string(outputJSON))
				}
			} else {
				panic(err)
			}

		} else {
			bld.Alert("No quantum circuit to run the software simulation")
		}
	}

	if *buildApp {
		// Build an hardware connected app

		modeTags := []string{"real", "complex"}

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
