package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bondgo"
	"github.com/BondMachineHQ/BondMachine/pkg/etherbond"
	"github.com/BondMachineHQ/BondMachine/pkg/udpbond"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var registerSize = flag.Int("register-size", 8, "Number of bits per register (n-bit)")

var mpm = flag.Bool("mpm", false, "Use the multi processor mode")
var cascading_io = flag.Bool("cascading-io", false, "Connect the processors in cascading io")

var o = flag.Int("O", 0, "Optimization level")

var showRequirements = flag.Bool("show-requirements", false, "Show bondmachine requirements")

var input_file = flag.String("input-file", "", "Go input file")

var multiAbstractAssemblyInput = flag.Bool("multi-abstract-assembly-input", false, "Input from a multi abstract assembly JSON file")
var abstractAssemblyInput = flag.Bool("abstract-assembly-input", false, "Input from abstract assembly")
var assemblyInput = flag.Bool("assembly-input", false, "Input from assembly")
var goInput = flag.Bool("go-input", true, "Input from Go file")

// Compiler modes:
//
//	standard: produce assembly files
//	checking: produce assembly files and check if the given machine/bondmachine is suitable for the code execution
//	enforcing: the produced assembly has to run on the given machine/bondmachine
//	optimizing: the optimized machine/bondmachine is created.
var compiler_mode = flag.String("compiler-mode", "standard", "Compiler mode: standard, checking, enforcing, optimizing")

// Loaded things
// For checking,enforcing
var loadBondmachine = flag.String("load-bondmachine", "", "Filename of the bondmachine to load")
var loadMachine = flag.String("load-machine", "", "Filename of the machine")

// Saved things
// For optimizing mandatory, optional for the others
var saveMachine = flag.String("save-machine", "", "Create a machine JSON file")
var saveBondmachine = flag.String("save-bondmachine", "", "Create a bondmachine JSON file")

// For standard, checking, enforcing
var saveAssembly = flag.String("save-assembly", "", "Machine or bondmachine (numbered per domain) assembly output file")

var useEtherbond = flag.Bool("use-etherbond", false, "Build including etherbond support")
var etherbond_external = flag.String("etherbond-external", "", "Etherbond external peers description file")
var save_etherbond_cluster = flag.String("save-etherbond-cluster", "ebcluster", "Create several BM files and the cluster file with the given prefix")

var use_udpbond = flag.Bool("use-udpbond", false, "Build including udpbond support")
var udpbond_external = flag.String("udpbond-external", "", "Udobond external peers description file")
var save_udpbond_cluster = flag.String("save-udpbond-cluster", "ebcluster", "Create several BM files and the cluster file with the given prefix")

var save_redeployer_file = flag.String("save-redeployer-file", "", "Create a redeployer file out of the cluster")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	flag.Parse()
	if *saveAssembly == "" {
		if *mpm {
			if *saveBondmachine == "" {
				*saveAssembly = "a.out.asm"
			}
		} else {
			if *saveMachine == "" {
				*saveAssembly = "a.out.asm"
			}
		}
	}

	if *abstractAssemblyInput || *multiAbstractAssemblyInput {
		*goInput = false
	}
}

func main() {

	fset := token.NewFileSet()

	// Setting the global compiler option that will be exported everywhere
	config := new(bondgo.BondgoConfig)
	config.Debug = *debug
	config.Verbose = *verbose
	if *mpm {
		config.Mpm = true
	} else {
		config.Mpm = false
	}

	if *cascading_io {
		config.Cascading_io = true
	} else {
		config.Cascading_io = false
	}

	switch *registerSize {
	case 8:
		config.Rsize = uint8(8)
		config.Basic_type = "uint8"
		config.Basic_chantype = "chan uint8"
	case 16:
		config.Rsize = uint8(16)
		config.Basic_type = "uint16"
		config.Basic_chantype = "chan uint16"
	case 32:
		config.Rsize = uint8(32)
		config.Basic_type = "uint32"
		config.Basic_chantype = "chan uint32"
	case 64:
		config.Rsize = uint8(64)
		config.Basic_type = "uint64"
		config.Basic_chantype = "chan uint64"
	default:
		fmt.Println("The specified register_size is not usable, defaulting to 8")
		config.Rsize = uint8(8)
		config.Basic_type = "uint8"
		config.Basic_chantype = "chan uint8"
	}

	if *input_file != "" {
		if *goInput {
			f, err := parser.ParseFile(fset, *input_file, nil, 0)
			if err != nil {
				fmt.Println(err)
				return
			}

			usagedone := make(chan bool)
			assignerdone := make(chan bool)

			results := new(bondgo.BondgoResults) // Results go in here
			results.Init_Results(config)

			messages := new(bondgo.BondgoMessages) // Compiler logs and errors
			messages.Init_Messages(config)

			reqmnts := new(bondgo.BondgoRequirements) // The pointer to the requirements struct
			reqmnts.Init_Requirements(config)

			usagenotify := make(chan bondgo.UsageNotify) // Used to notify the used resource

			go reqmnts.Usage_Monitor(usagenotify, usagedone) // Spawn the usage monitor

			run := new(bondgo.BondgoRuninfo) // Running data
			run.Init_Runinfo(config)

			varreq := make(chan bondgo.VarReq) // Variable request
			varans := make(chan bondgo.VarAns) // Variable response

			go run.Var_assigner(varreq, varans, usagenotify, assignerdone) // Spawn the variable allocator

			functs := new(bondgo.BondgoFunctions) // Functions
			functs.Init_Functions(config, messages)

			vars := make(map[string]bondgo.VarCell)
			returns := make([]bondgo.VarCell, 0)

			bgmain := &bondgo.BondgoCheck{results, config, reqmnts, run, messages, functs, usagenotify, varreq, varans, nil, nil, vars, returns, "", "", "device_0", 0}

			bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, 0, bondgo.C_DEVICE, bgmain.CurrentDevice, bondgo.I_NIL}

			if config.Debug {
				ast.Print(fset, f)
			}

			// Load all the functions
			ast.Walk(functs, f)

			if !bgmain.Is_faulty() {

				// Start the main function
				executable := false
				for ifuncname, ifunc := range functs.Functions {
					if ifuncname == "main" {
						start_ast := ifunc.Body
						ast.Walk(bgmain, start_ast)
						executable = true
						break
					}
				}

				if !executable {
					bgmain.Set_faulty("main function not found.")
				}

				for procid, _ := range bgmain.Program {
					// Add io connections
					if bgmain.Cascading_io {
						bgmain.WriteLine(procid, "i2r r0 i0")
						bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_OPCODE, "i2r", bondgo.I_NIL}
						bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_INPUT, bondgo.S_NIL, 1}
						bgmain.WriteLine(procid, "r2o r0 o0")
						bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_OPCODE, "r2o", bondgo.I_NIL}
						bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_OUTPUT, bondgo.S_NIL, 1}
					}
				}

				for procid, rout := range bgmain.Program {
					// TODO Recheck
					linesn := len(rout.Lines)
					bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_ROMSIZE, bondgo.S_NIL, linesn}
				}

				bgmain.Used <- bondgo.UsageNotify{bondgo.TR_EXIT, 0, 0, bondgo.S_NIL, bondgo.I_NIL}
				<-usagedone

				gent, _ := bondgo.Type_from_string(bgmain.Basic_type)
				bgmain.Reqs <- bondgo.VarReq{bondgo.REQ_EXIT, 0, bondgo.VarCell{gent, 0, 0, 0, 0, 0, 0, 0}}
				<-assignerdone
			}

			fmt.Print(bgmain.Dump_log())

			if !bgmain.Is_faulty() {

				if *showRequirements {
					fmt.Print(bgmain.Dump_Requirements())
				}

				switch *compiler_mode {
				case "standard":

					var machine_to_save string
					savedsomething := false

					if *mpm {
						machine_to_save = *saveBondmachine
					} else {
						machine_to_save = *saveMachine
					}

					if *saveAssembly != "" {
						savedsomething = true
						if *mpm || *useEtherbond || *use_udpbond {
							for i, _ := range bgmain.Program {
								if _, err := os.Stat(*saveAssembly + "_" + strconv.Itoa(i)); os.IsNotExist(err) {
									f, err := os.Create(*saveAssembly + "_" + strconv.Itoa(i))
									check(err)
									defer f.Close()
									f.WriteString(bgmain.Write_assembly(i))
								}
							}
						} else {
							if _, err := os.Stat(*saveAssembly); os.IsNotExist(err) {
								f, err := os.Create(*saveAssembly)
								check(err)
								defer f.Close()
								f.WriteString(bgmain.Write_assembly(0))
							}
						}
					}

					if *useEtherbond {
						if *save_etherbond_cluster != "" {

							var external_cluster *etherbond.Cluster

							if *etherbond_external != "" {

								ebconfig := new(etherbond.Config)
								ebconfig.Rsize = uint8(*registerSize)

								if *debug {
									ebconfig.Debug = true
								}

								if external_cluster_t, err := etherbond.UnmarshallCluster(ebconfig, *etherbond_external); err != nil {
									panic(err)
								} else {
									external_cluster = external_cluster_t
								}
							}

							mycluster, peerids, mymachines, myethio, myres, err := bgmain.Create_Etherbond_Cluster(*registerSize, external_cluster)
							check(err)
							//fmt.Println(mycluster, peerids, mymachines, myethio, myres)

							if _, err := os.Stat(*save_etherbond_cluster + ".json"); os.IsNotExist(err) {
								f, err := os.Create(*save_etherbond_cluster + ".json")
								check(err)
								b, errj := json.Marshal(mycluster)
								check(errj)
								_, err = f.WriteString(string(b))
								check(err)
								f.Close()
							}

							for i, peerid := range peerids {
								machine_to_save := *save_etherbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_bm.json"
								if _, err := os.Stat(machine_to_save); os.IsNotExist(err) {
									f, err := os.Create(machine_to_save)
									check(err)
									b, errj := json.Marshal(mymachines[i].Jsoner())
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
								io_to_save := *save_etherbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_io.json"
								if _, err := os.Stat(io_to_save); os.IsNotExist(err) {
									f, err := os.Create(io_to_save)
									check(err)
									b, errj := json.Marshal(myethio[i])
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
								residual_to_save := *save_etherbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_residual.json"
								if _, err := os.Stat(residual_to_save); os.IsNotExist(err) {
									f, err := os.Create(residual_to_save)
									check(err)
									b, errj := json.Marshal(myres[i])
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
							}

						} else {
							fmt.Println("Missing output file name")
						}

					} else if *use_udpbond {
						if *save_udpbond_cluster != "" {

							var external_cluster *udpbond.Cluster

							if *udpbond_external != "" {

								ebconfig := new(udpbond.Config)
								ebconfig.Rsize = uint8(*registerSize)

								if *debug {
									ebconfig.Debug = true
								}

								if external_cluster_t, err := udpbond.UnmarshallCluster(ebconfig, *udpbond_external); err != nil {
									panic(err)
								} else {
									external_cluster = external_cluster_t
								}
							}

							mycluster, peerids, mymachines, myethio, myres, err := bgmain.Create_Udpbond_Cluster(*registerSize, external_cluster)
							check(err)
							//fmt.Println(mycluster, peerids, mymachines, myethio, myres)

							if _, err := os.Stat(*save_udpbond_cluster + ".json"); os.IsNotExist(err) {
								f, err := os.Create(*save_udpbond_cluster + ".json")
								check(err)
								b, errj := json.Marshal(mycluster)
								check(errj)
								_, err = f.WriteString(string(b))
								check(err)
								f.Close()
							}

							for i, peerid := range peerids {
								machine_to_save := *save_udpbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_bm.json"
								if _, err := os.Stat(machine_to_save); os.IsNotExist(err) {
									f, err := os.Create(machine_to_save)
									check(err)
									b, errj := json.Marshal(mymachines[i].Jsoner())
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
								io_to_save := *save_udpbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_io.json"
								if _, err := os.Stat(io_to_save); os.IsNotExist(err) {
									f, err := os.Create(io_to_save)
									check(err)
									b, errj := json.Marshal(myethio[i])
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
								residual_to_save := *save_udpbond_cluster + "_peer_" + strconv.Itoa(int(peerid)) + "_residual.json"
								if _, err := os.Stat(residual_to_save); os.IsNotExist(err) {
									f, err := os.Create(residual_to_save)
									check(err)
									b, errj := json.Marshal(myres[i])
									check(errj)
									_, err = f.WriteString(string(b))
									check(err)
									defer f.Close()
								}
							}

						} else {
							fmt.Println("Missing output file name")
						}

					} else {
						if machine_to_save != "" {
							savedsomething = true
							if *mpm {
								if mymachine, _, err := bgmain.Create_Bondmachine(*registerSize, "device_0"); err == nil {
									if _, err := os.Stat(machine_to_save); os.IsNotExist(err) {
										f, err := os.Create(machine_to_save)
										check(err)
										defer f.Close()
										b, errj := json.Marshal(mymachine.Jsoner())
										check(errj)
										_, err = f.WriteString(string(b))
										check(err)
									}
								} else {
									fmt.Println("Creating bondmachine failed")
								}

							} else {
								if mymachine, ok := bgmain.Create_Connecting_Processor(*registerSize, 0); ok {
									if _, err := os.Stat(machine_to_save); os.IsNotExist(err) {
										f, err := os.Create(machine_to_save)
										check(err)
										defer f.Close()
										b, errj := json.Marshal(mymachine.Jsoner())
										check(errj)
										_, err = f.WriteString(string(b))
										check(err)
									}
								} else {
									fmt.Println("Creating processor failed")
								}
							}
						}

						if !savedsomething {
							fmt.Println("Missing output file name")
						}
					}
				case "checking":
					fmt.Println("TODO Uninplemented")
				case "enforcing":
					fmt.Println("TODO Uninplemented")
				case "optimizing":
					fmt.Println("TODO Uninplemented")
				default:
					fmt.Println("Unknown operating mode")
				}
			}
		} else if *assemblyInput {
			// Placeholder for the new assembly engine
		} else if *abstractAssemblyInput {
			usagedone := make(chan bool)

			results := new(bondgo.BondgoResults) // Results go in here
			results.Init_Results(config)

			messages := new(bondgo.BondgoMessages) // Compiler logs and errors
			messages.Init_Messages(config)

			reqmnts := new(bondgo.BondgoRequirements) // The pointer to the requirements struct
			reqmnts.Init_Requirements(config)

			usagenotify := make(chan bondgo.UsageNotify) // Used to notify the used resource

			go reqmnts.Usage_Monitor(usagenotify, usagedone) // Spawn the usage monitor

			run := new(bondgo.BondgoRuninfo) // Running data
			run.Init_Runinfo(config)

			varreq := make(chan bondgo.VarReq) // Variable request
			varans := make(chan bondgo.VarAns) // Variable response

			functs := new(bondgo.BondgoFunctions) // Functions
			functs.Init_Functions(config, messages)

			vars := make(map[string]bondgo.VarCell)
			returns := make([]bondgo.VarCell, 0)

			bgmain := &bondgo.BondgoCheck{results, config, reqmnts, run, messages, functs, usagenotify, varreq, varans, nil, nil, vars, returns, "", "", "device_0", 0}

			// Establish the Device parameter for the next goroutine
			bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, 0, bondgo.C_DEVICE, bgmain.CurrentDevice, bondgo.I_NIL}

			f, err := os.Open(*input_file)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			source_asm := make([]string, 0)

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				currline := scanner.Text()
				source_asm = append(source_asm, currline)
				bgmain.WriteLine(0, currline)
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			bgmain.Abstract_assembler(*registerSize, source_asm, usagenotify)

			for procid, rout := range bgmain.Program {
				// TODO Recheck
				linesn := len(rout.Lines)
				bgmain.Used <- bondgo.UsageNotify{bondgo.TR_PROC, procid, bondgo.C_ROMSIZE, bondgo.S_NIL, linesn}
			}

			bgmain.Used <- bondgo.UsageNotify{bondgo.TR_EXIT, 0, 0, bondgo.S_NIL, bondgo.I_NIL}
			<-usagedone

			if *saveBondmachine != "" {
				if mymachine, _, err := bgmain.Create_Bondmachine(*registerSize, "device_0"); err == nil {
					if _, err := os.Stat(*saveBondmachine); os.IsNotExist(err) {
						f, err := os.Create(*saveBondmachine)
						check(err)
						defer f.Close()
						b, errj := json.Marshal(mymachine.Jsoner())
						check(errj)
						_, err = f.WriteString(string(b))
						check(err)
					}
				} else {
					fmt.Println("Creating processor failed")
				}

			}

			if *saveMachine != "" {
				if mymachine, ok := bgmain.Create_Connecting_Processor(*registerSize, 0); ok {
					if _, err := os.Stat(*saveMachine); os.IsNotExist(err) {
						f, err := os.Create(*saveMachine)
						check(err)
						defer f.Close()
						b, errj := json.Marshal(mymachine.Jsoner())
						check(errj)
						_, err = f.WriteString(string(b))
						check(err)
					}
				} else {
					fmt.Println("Creating processor failed")
				}
			}
		} else if *multiAbstractAssemblyInput {

			aafile := new(bondgo.Abs_assembly)

			if _, err := os.Stat(*input_file); err == nil {
				if jsonfile, err := ioutil.ReadFile(*input_file); err == nil {
					if err := json.Unmarshal([]byte(jsonfile), aafile); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}

			if *saveBondmachine != "" {
				if mymachine, err := bondgo.MultiAsm2BondMachine(*registerSize, aafile); err == nil {
					if _, err := os.Stat(*saveBondmachine); os.IsNotExist(err) {
						f, err := os.Create(*saveBondmachine)
						check(err)
						defer f.Close()
						b, errj := json.Marshal(mymachine.Jsoner())
						check(errj)
						_, err = f.WriteString(string(b))
						check(err)
					}
				} else {
					fmt.Println("Creating processor failed")
				}
			}
		} else {
			fmt.Println("Missing input file type")
		}
	} else {
		fmt.Println("Missing input file")
	}
}
