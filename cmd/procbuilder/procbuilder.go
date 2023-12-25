package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"

	//"mel"

	"math/rand"
	"os"

	//"strconv"
	"sort"
	"strings"
	"time"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var execution_model = flag.String("execution-model", "ha", "Execution model: vn (Von Neumann), ha (Harvard), hy (Hybrid)")

var register_size = flag.Int("register-size", 8, "Number of bits per register (n-bit)")

var enabledOpcodes = flag.String("opcodes", "nop", "Enabled opcodes")
var listOpcodes = flag.Bool("list-opcodes", false, "List all opcodes")

var rbit = flag.Int("registers", 3, "Number of n-bit registers 2^")
var lbit = flag.Int("ram", 8, "Number of n-bit RAM memory cells 2^")
var nbit = flag.Int("inputs", 1, "Number of n-bit inputs")
var mbit = flag.Int("outputs", 1, "Number of n-bit outputs")
var obit = flag.Int("rom", 8, "Number of ROM memory cells 2^")

var inputAssembly = flag.String("input-assembly", "", "Take assembly file as input")
var inputBinary = flag.String("input-binary", "", "Take binary file as input")
var inputRandom = flag.Bool("input-random", false, "Generate a random input")

var opcode_optimizer = flag.Bool("opcode-optimizer", false, "Activate opecode optimizator for assembly input")

var create_verilog = flag.Bool("create-verilog", false, "Create default verilog files")
var create_verilog_processor = flag.String("create-verilog-processor", "", "Filename of verilog processor")
var create_verilog_ram = flag.String("create-verilog-ram", "", "Filename of verilog ram")
var create_verilog_rom = flag.String("create-verilog-rom", "", "Filename of verilog rom")
var create_verilog_arch = flag.String("create-verilog-arch", "", "Filename of verilog arch")
var create_verilog_testbench = flag.String("create-verilog-testbench", "", "Filename of verilog testbench")
var create_verilog_main = flag.String("create-verilog-main", "", "Filename of verilog main file for FPGA")
var verilog_flavor = flag.String("verilog-flavor", "iverilog", "Choose the type of verilog device. currently supported: iverilog,kintex7.")

var show_instructions_alias = flag.Bool("show-instructions-alias", false, "Show instructions alias for the processor")
var show_program_alias = flag.Bool("show-program-alias", false, "Show program alias for the processor")

var load_machine = flag.String("load-machine", "", "Load machine in JSON format")
var save_machine = flag.String("save-machine", "", "Save machine in JSON format")

var sharedConstraints = flag.String("shared-constraints", "", "List of shared objects connected to the processor")

var show_program_binary = flag.Bool("show-program-binary", false, "Show program binary")
var showProgramDisassembled = flag.Bool("show-program-disassembled", false, "Show disassebled program")

var showOpcodes = flag.Bool("show-opcodes", false, "Show loaded opcodes")
var showOpcodesDetails = flag.Bool("show-opcodes-details", false, "Show loaded opcodes and details")

var hex = flag.Bool("hex", false, "Use HEX")
var numlines = flag.Bool("numlines", false, "Use line numbers")

var simboxFile = flag.String("simbox-file", "", "Filename of the simulation data file")

var sim = flag.Bool("sim", false, "Simulate machine")
var simInteractions = flag.Int("sim-interactions", 10, "Simulation interaction")

var run = flag.Bool("run", false, "Run machine")
var runInteractions = flag.Int("run-interactions", 1000, "Run interaction")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	rand.Seed(int64(time.Now().Unix()))
	flag.Parse()
}

func main() {
	//ep := new(mel.Evolution_parameters)
	//ep.Pars = make(map[string]string)

	ri := new(procbuilder.RuntimeInfo)
	ri.Init()

	conf := new(procbuilder.Config)
	conf.Debug = *debug
	conf.Runinfo = ri

	var mymachine *procbuilder.Machine
	var myarch *procbuilder.Arch

	if *load_machine != "" {
		if _, err := os.Stat(*load_machine); err == nil {
			if jsonfile, err := os.ReadFile(*load_machine); err == nil {
				var machj procbuilder.Machine_json
				if err := json.Unmarshal([]byte(jsonfile), &machj); err == nil {
					mymachine = (&machj).Dejsoner()
					myarch = &mymachine.Arch
				} else {
					panic(err)
				}

			} else {
				panic(err)
			}

		} else {
			panic(err)
		}
	} else {
		mymachine = new(procbuilder.Machine)

		myarch = &mymachine.Arch

		myarch.Rsize = uint8(*register_size)

		if *execution_model == "ha" || *execution_model == "vn" || *execution_model == "hy" {
			myarch.Modes = make([]string, 1)
			myarch.Modes[0] = *execution_model
		} else {
			panic("Unknown execution model")
		}

		myarch.R = uint8(*rbit)
		myarch.L = uint8(*lbit)
		myarch.N = uint8(*nbit)
		myarch.M = uint8(*mbit)
		myarch.O = uint8(*obit)
		myarch.Shared_constraints = *sharedConstraints

		//ep.Pars["procbuilder:opcodes"] = *enabled_opcodes
		//ep.Pars["procbuilder:r"] = strconv.Itoa(*rbit)
		//ep.Pars["procbuilder:l"] = strconv.Itoa(*lbit)
		//ep.Pars["procbuilder:m"] = strconv.Itoa(*mbit)
		//ep.Pars["procbuilder:n"] = strconv.Itoa(*nbit)
		//ep.Pars["procbuilder:o"] = strconv.Itoa(*obit)

		if *listOpcodes {
			for _, op := range procbuilder.Allopcodes {
				fmt.Println(op.Op_get_name())
			}
		}

		// Processing enabled opcodes
		if *inputAssembly != "" && *opcode_optimizer {
			if _, err := os.Stat(*inputAssembly); err == nil {
				if prog, err := os.ReadFile(*inputAssembly); err == nil {

					// TODO keep the opcodes ordered by name
					opcodes := make([]procbuilder.Opcode, 0)

					curLine := make([]byte, 256)

					iLine := 0

					for _, ch := range prog {
						if ch == 10 {
							curLine[iLine] = ' '
							if len(strings.Split(string(curLine), " ")) > 0 {
								tCheck := false
								opn := strings.Split(string(curLine), " ")[0]
								for _, op := range opcodes {
									if opn == op.Op_get_name() {
										tCheck = true
										break
									}
								}

								if !tCheck {
									for _, op := range procbuilder.Allopcodes {
										if opn == op.Op_get_name() {
											opcodes = append(opcodes, op)
											break
										}
									}
								}
							}
							iLine = 0

						} else {
							curLine[iLine] = ch
							iLine = iLine + 1
						}
					}

					sort.Sort(procbuilder.ByName(opcodes))

					myarch.Op = opcodes
					mymachine.Arch = *myarch

				} else {
					panic(err)
				}

			} else {
				panic(err)
			}
		} else {
			var eOps []string
			if *enabledOpcodes != "" {
				eOps = strings.Split(*enabledOpcodes, ",")
			} else {
				panic("Missing opcodes")
			}

			opcodes := make([]procbuilder.Opcode, len(eOps))
			checked := make(map[string]struct{})

			for i, opName := range eOps {
				// This is the support for dynamic opcodes. It eventually creates the opcode if it does not exist and follows the naming convention
				if _, err := procbuilder.EventuallyCreateInstruction(opName); err != nil {
					panic(err)
				}

				found := false
				for _, op := range procbuilder.Allopcodes {
					if op.Op_get_name() == opName {
						if _, ok := checked[opName]; ok {
							panic("Duplicate opcode " + opName)
						}
						opcodes[i] = op
						found = true
						checked[opName] = struct{}{}
						break
					}
				}
				if !found {
					panic("Unknown opcode " + opName)
				}
			}

			sort.Sort(procbuilder.ByName(opcodes))

			myarch.Op = opcodes
			mymachine.Arch = *myarch
		}

		// Precessing assembly
		if *inputAssembly != "" {
			if _, err := os.Stat(*inputAssembly); err == nil {
				if prog, err := os.ReadFile(*inputAssembly); err == nil {
					if prog, err := myarch.Assembler(prog); err == nil {
						mymachine.Program = prog
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}

			} else {
				panic(err)
			}
		} else if *inputBinary != "" {
			//TODO input from binary file
		} else if *inputRandom {
			//mymachine = procbuilder.Machine_Program_Generate(ep).(*procbuilder.Machine)
		} else {
			fmt.Println("Warning no program loaded")
		}
	}

	if checks, ok := mymachine.ConstraintCheck(); ok {

		if *verbose {
			fmt.Print(checks)
		}

		// Eventually show alias instrictions data
		if *show_instructions_alias {
			if alias_text, err := mymachine.Instructions_alias(); err == nil {
				fmt.Print(alias_text)
			} else {
				panic(err)
			}
		}

		// Eventually show alias instrictions data
		if *show_program_alias {
			if alias_text, err := mymachine.Program_alias(); err == nil {
				fmt.Print(alias_text)
			} else {
				panic(err)
			}
		}

		// Eventually show opcodes
		if *showOpcodes {
			for _, op := range myarch.Op {
				fmt.Println(op.Op_get_name())
			}
		}

		// Eventually show opcodes details
		if *showOpcodesDetails {
			for _, op := range myarch.Op {
				fmt.Println(op.Op_show_assembler(&mymachine.Arch))
			}
		}

		// Eventually show program
		if *showProgramDisassembled {
			if disass_text, err := mymachine.Disassembler(); err == nil {
				fmt.Print(disass_text)
			} else {
				panic(err)
			}
		}

		if *show_program_binary {
			for i, inst := range mymachine.Program.Slocs {
				if *numlines {
					fmt.Printf("%5d %s\n", i, inst)
				} else {
					fmt.Printf("%s\n", inst)
				}
			}
		}

		// Eventually create JSON machines file
		if *save_machine != "" {
			if _, err := os.Stat(*save_machine); os.IsNotExist(err) {
				f, err := os.Create(*save_machine)
				check(err)
				defer f.Close()
				b, errj := json.Marshal(mymachine.Jsoner())
				check(errj)
				_, err = f.WriteString(string(b))
				check(err)
			}
		}

		// Eventually create verilog files
		if *create_verilog {
			if _, err := os.Stat("arch.v"); os.IsNotExist(err) {
				f, err := os.Create("arch.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Write_verilog("a0", map[string]string{"processor": "p0", "rom": "p0rom", "ram": "p0ram"}, *verilog_flavor))
				check(err)
			}

			if _, err := os.Stat("processor.v"); os.IsNotExist(err) {
				f, err := os.Create("processor.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Conproc.Write_verilog(conf, myarch, "p0", *verilog_flavor))
				check(err)
			}

			if _, err := os.Stat("ram.v"); os.IsNotExist(err) {
				f, err := os.Create("ram.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Ram.Write_verilog(nil, mymachine, "p0ram", *verilog_flavor))
				check(err)
			}

			if _, err := os.Stat("rom.v"); os.IsNotExist(err) {
				f, err := os.Create("rom.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Rom.Write_verilog(mymachine, "p0rom", *verilog_flavor))
				check(err)
			}

			if _, err := os.Stat("testbench.v"); os.IsNotExist(err) {
				f, err := os.Create("testbench.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Write_verilog_testbench("a0", "processor", "memory", *verilog_flavor))
				check(err)
			}

			if _, err := os.Stat("main.v"); os.IsNotExist(err) {
				f, err := os.Create("main.v")
				check(err)
				defer f.Close()
				_, err = f.WriteString(myarch.Write_verilog_main("p0", "p0rom", "processor", "memory", *verilog_flavor))
				check(err)
			}
		} else {
			if *create_verilog_processor != "" {
				if _, err := os.Stat(*create_verilog_processor); os.IsNotExist(err) {
					f, err := os.Create(*create_verilog_processor)
					check(err)
					defer f.Close()
					_, err = f.WriteString(myarch.Conproc.Write_verilog(conf, myarch, "p0", *verilog_flavor))
					check(err)
				}
			}

			if *create_verilog_ram != "" {
				if _, err := os.Stat(*create_verilog_ram); os.IsNotExist(err) {
					f, err := os.Create(*create_verilog_ram)
					check(err)
					defer f.Close()
					_, err = f.WriteString(myarch.Ram.Write_verilog(nil, mymachine, "p0ram", *verilog_flavor))
					check(err)
				}
			}

			if *create_verilog_rom != "" {
				if _, err := os.Stat(*create_verilog_rom); os.IsNotExist(err) {
					f, err := os.Create(*create_verilog_rom)
					check(err)
					defer f.Close()
					_, err = f.WriteString(myarch.Rom.Write_verilog(mymachine, "p0rom", *verilog_flavor))
					check(err)
				}
			}

			if *create_verilog_testbench != "" {
				if _, err := os.Stat(*create_verilog_testbench); os.IsNotExist(err) {
					f, err := os.Create(*create_verilog_testbench)
					check(err)
					defer f.Close()
					_, err = f.WriteString(myarch.Write_verilog_testbench("a0", "processor", "memory", *verilog_flavor))
					check(err)
				}
			}

			if *create_verilog_main != "" {
				if _, err := os.Stat(*create_verilog_main); os.IsNotExist(err) {
					f, err := os.Create(*create_verilog_main)
					check(err)
					defer f.Close()
					_, err = f.WriteString(myarch.Write_verilog_main("p0", "p0rom", "processor", "memory", *verilog_flavor))
					check(err)
				}
			}
		}
		if *sim {
			var sbox *simbox.Simbox
			if *simboxFile != "" {
				sbox = new(simbox.Simbox)
				if _, err := os.Stat(*simboxFile); err == nil {
					// Open the simbox file is exists
					if simbox_json, err := os.ReadFile(*simboxFile); err == nil {
						if err := json.Unmarshal([]byte(simbox_json), sbox); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}

			}

			// Build the VM
			vm := new(procbuilder.VM)
			vm.Mach = mymachine
			err := vm.Init()
			check(err)

			sconfig := new(procbuilder.Sim_config)
			scerr := sconfig.Init(sbox, vm)
			check(scerr)

			// Build the simulation driver
			sdrive := new(procbuilder.Sim_drive)
			sderr := sdrive.Init(sbox, vm)
			check(sderr)

			// Build the simultion report
			srep := new(procbuilder.Sim_report)
			srerr := srep.Init(sbox, vm)
			check(srerr)

			for i := uint64(0); i < uint64(*simInteractions); i++ {
				if sconfig.Show_pc {
					fmt.Println("Program Counter:", vm.Pc)
				}
				fmt.Println("Instruction: ", vm.Mach.Slocs[vm.Pc])
				fmt.Println("Registers before: ", vm.DumpRegisters())
				fmt.Println("IO before: ", vm.DumpIO())

				// This will get actions eventually to do on this tick
				if act, exist_actions := sdrive.AbsSet[i]; exist_actions {
					for i, val := range act {
						*sdrive.Injectables[i] = val
					}
				}

				_, err := vm.Step(sconfig)
				check(err)

				// This will get value to report on this tick
				if rep, exist_reports := srep.AbsGet[i]; exist_reports {
					for i, _ := range rep {
						rep[i] = *srep.Reportables[i]
					}
					fmt.Println(rep)
				}
				fmt.Println("Registers after: ", vm.DumpRegisters())
				fmt.Println("IO after: ", vm.DumpIO(), "\n")
			}
		} else if *run {
			// TODO The sdrive and report goes also here
			vm := new(procbuilder.VM)
			vm.Mach = mymachine
			err := vm.Init()
			check(err)
			for i := 0; i < *runInteractions; i++ {
				_, err := vm.Step(nil)
				check(err)
			}
		}
	} else {
		fmt.Println("Constraint check failed: " + checks)
	}
}
