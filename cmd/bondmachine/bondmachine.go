package main

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BondMachineHQ/BondMachine/pkg/bcof"
	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/bondirect"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/etherbond"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
	"github.com/BondMachineHQ/BondMachine/pkg/udpbond"
	"google.golang.org/protobuf/proto"
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

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")
var commentedVerilog = flag.Bool("comment-verilog", false, "Comment generated verilog")

var register_size = flag.Int("register-size", 8, "Number of bits per register (n-bit)")

var bondmachine_file = flag.String("bondmachine-file", "", "Filename of the bondmachine")

var bcofInFile = flag.String("bcof-file", "", "Use a BCOF file as input for RAM initialization")

// Verilog processing
var create_verilog = flag.Bool("create-verilog", false, "Create default verilog files")
var verilog_flavor = flag.String("verilog-flavor", "iverilog", "Choose the type of verilog device. currently supported: iverilog,de10nano.")
var verilog_mapfile = flag.String("verilog-mapfile", "", "File mapping the device IO to bondmachine IO")
var verilog_simulation = flag.Bool("verilog-simulation", false, "Create simulation oriented verilog as default.")

var show_program_alias = flag.Bool("show-program-alias", false, "Show program alias for the processor")

// Domains processing
var list_domains = flag.Bool("list-domains", false, "Domain list")
var add_domains string_slice
var del_domains string_slice

// Processors
var list_processors = flag.Bool("list-processors", false, "Processor list")
var add_processor = flag.Int("add-processor", -1, "Add a processor of the given domain")
var enumProcessors = flag.Bool("enum-processors", false, "Enumerate all the processors")

// TODO del-processor

// Inputs
var list_inputs = flag.Bool("list-inputs", false, "Inputs list")
var add_inputs = flag.Int("add-inputs", 0, "Inputs to add") // When adding we need how many new inputs, when removing we need which (thats why the list)
var del_inputs string_slice

// Outputs
var list_outputs = flag.Bool("list-outputs", false, "Outputs list")
var add_outputs = flag.Int("add-outputs", 0, "Outputs to add")
var del_outputs string_slice

var list_bonds = flag.Bool("list-bonds", false, "Bonds list")
var enumBonds = flag.Bool("enum-bonds", false, "Enumerate all the bonds")
var add_bond string_slice
var del_bonds string_slice

var specs = flag.Bool("specs", false, "Show BondMachine specs")

// TODO Shared objects
var list_shared_objects = flag.Bool("list-shared-objects", false, "Shared object list")
var add_shared_objects string_slice
var del_shared_objects string_slice
var list_processor_shared_object_links = flag.Bool("list-processor-shared-object-links", false, "Processor shared object link list")
var connect_processor_shared_object string_slice
var disconnect_processor_shared_object string_slice

var list_internal_inputs = flag.Bool("list-internal-inputs", false, "Internal inputs list")
var list_internal_outputs = flag.Bool("list-internal-outputs", false, "Internal outputs list")

// Dot output
var emit_dot = flag.Bool("emit-dot", false, "Emit dot file on stdout")
var dot_detail = flag.Int("dot-detail", 1, "Detail of infos on dot file 1-5")

// Assembly output
var show_program_disassembled = flag.Bool("show-program-disassembled", false, "Show disassebled program")
var multi_abstract_assembly_file = flag.String("multi-abstract-assembly-file", "", "Save the bondmachine as multi abstract assembly file")

var simbox_file = flag.String("simbox-file", "", "Filename of the simulation data file")

var sim = flag.Bool("sim", false, "Simulate bond machine")
var simInteractions = flag.Int("sim-interactions", 10, "Simulation interaction")
var simStopOnValidOf = flag.Int("sim-stop-on-valid-of", -1, "Stop simulation when a valid output is produced on the given output")
var simReport = flag.String("sim-report", "", "Simulation report file")

var emu = flag.Bool("emu", false, "Emulate bond machine")
var emu_interactions = flag.Int("emu-interactions", 10, "Emulation interaction (0 means forever)")

var clusterSpec = flag.String("cluster-spec", "", "Cluster Spec File ")
var peerID = flag.Int("peer-id", -1, "Peer ID of the BondMachine within the cluster")

var use_etherbond = flag.Bool("use-etherbond", false, "Build including etherbond support")
var etherbond_flavor = flag.String("etherbond-flavor", "enc60j28", "Choose the type of ethernet device. currently supported: enc60j28.")
var etherbond_mapfile = flag.String("etherbond-mapfile", "", "File mapping the bondmachine IO the etherbond.")
var etherbond_macfile = flag.String("etherbond-macfile", "", "File mapping the bondmachine peers to MAC addresses.")

var use_udpbond = flag.Bool("use-udpbond", false, "Build including udpbond support")
var udpbond_flavor = flag.String("udpbond-flavor", "esp8266", "Choose the type of network device. currently supported: esp8266.")
var udpbondMapfile = flag.String("udpbond-mapfile", "", "File mapping the bondmachine IO the udpbond.")
var udpbond_ipfile = flag.String("udpbond-ipfile", "", "File mapping the bondmachine peers to IP addresses.")
var udpbond_netconfig = flag.String("udpbond-netconfig", "", "JSON file containing the network configuration for udpbond")

var useBondirect = flag.Bool("use-bondirect", false, "Build including bondirect support")
var bondirectFlavor = flag.String("bondirect-flavor", "basic", "Choose the type of bondirect device. currently supported: basic.")
var bondirectMapfile = flag.String("bondirect-mapfile", "", "File mapping the bondmachine IO the bondirect.")
var bondirectMesh = flag.String("bondirect-mesh", "", "Bondirect mesh File ")

var usebmapi = flag.Bool("use-bmapi", false, "Build a BMAPI interface")
var bmapiLanguage = flag.String("bmapi-language", "go", "Choose BMAPI language (go,c,python)")
var bmapiFramework = flag.String("bmapi-framework", "", "Choose BMAPI framework (pynq)")
var bmapiFlavor = flag.String("bmapi-flavor", "", "Choose the BMAPI interconnect")
var bmapiFlavorVersion = flag.String("bmapi-flavor-version", "", "Choose the BMAPI interconnect version")
var bmapiMapfile = flag.String("bmapi-mapfile", "", "File mapping the bondmachine IO the BMAPI.")
var bmapiLibOutDir = flag.String("bmapi-liboutdir", "", "Output directory for the BMAPI library.")
var bmapiModOutDir = flag.String("bmapi-modoutdir", "", "Output directory for the BMAPI kernel module.")
var bmapiAuxOutDir = flag.String("bmapi-auxoutdir", "", "Output directory for the BMAPI auxiliary material.")
var bmapiPackageName = flag.String("bmapi-packagename", "", "GO package name.")
var bmapiModuleName = flag.String("bmapi-modulename", "", "GO module name.")
var bmapiGenerateExample = flag.String("bmapi-generate-example", "", "Generate an example program using the BMAPI.")
var bmapiDataType = flag.String("bmapi-data-type", "float32", "Data type for the BMAPI.")

var board_slow = flag.Bool("board-slow", false, "Board slow support")
var board_slow_factor = flag.Int("board-slow-factor", 1, "Board slow factor")

var counter = flag.Bool("counter", false, "Counter support")
var counterMap = flag.String("counter-map", "", "Counter mappings")
var counterSlowFactor = flag.String("counter-slow-factor", "23", "Counter slow factor")

var uart = flag.Bool("uart", false, "UART support")
var uartMapFile = flag.String("uart-mapfile", "", "UART mappings")

var basys3_7segment = flag.Bool("basys3-7segment", false, "Basys3 7 segments display support")
var basys3_7segment_map = flag.String("basys3-7segment-map", "", "Basys3 7 segments display mappings")

var basys3Leds = flag.Bool("basys3-leds", false, "Basys3 leds support")
var basys3LedsMap = flag.String("basys3-leds-map", "", "Basys3 leds mappings")
var basys3LedsName = flag.String("basys3-leds-name", "led", "Basys3 leds name")

var iceBreakerLeds = flag.Bool("icebreaker-leds", false, "Icebreaker leds support")
var iceBreakerLedsMap = flag.String("icebreaker-leds-map", "", "Icebreaker leds mappings")

var iceFunLeds = flag.Bool("icefun-leds", false, "IceFun leds support")
var iceFunLedsMap = flag.String("icefun-leds-map", "", "IceFun leds mappings")

var Ice40Lp1kLeds = flag.Bool("ice40lp1k-leds", false, "Ice40lp1k leds support")
var Ice40Lp1kLedsMap = flag.String("ice40lp1k-leds-map", "", "Ice40lp1k leds mappings")

var vgatext = flag.Bool("vgatext", false, "Multi CP VGA Textual support")
var vgatextFlavor = flag.String("vgatext-flavor", "800x600", "VGA Textual flavor. currently supported: 800x600")
var vgatextFonts = flag.String("vgatext-fonts", "", "VGA Textual fonts file")
var vgatextHeader = flag.String("vgatext-header", "", "VGA Textual header rom file")

var ps2keyboardIo = flag.Bool("ps2-keyboard-io", false, "PS2 Keyboard support via IO")
var ps2keyboardIoMap = flag.String("ps2-keyboard-io-map", "", "PS2 Keyboard via IO mappings")

var ps2keyboard = flag.Bool("ps2-keyboard", false, "PS2 Keyboard support via SO")

var attach_benchmark_core string_slice
var attachBenchmarkCoreV2 string_slice

var bmInfoFile = flag.String("bminfo-file", "", "File containing the bondmachine extra info")
var bmRequirementsFile = flag.String("bmrequirements-file", "", "File containing the bondmachine requirements")
var hwOptimizations = flag.String("hw-optimizations", "", "comma separated hardware optimizations")

var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the sintax index,filename)")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	rand.Seed(int64(time.Now().Unix()))

	flag.Var(&add_domains, "add-domains", "Comma-separated list of JSON machine files to add")
	flag.Var(&del_domains, "del-domains", "Comma-separated list of domain ID to delete")
	flag.Var(&del_inputs, "del-inputs", "Comma-separated list of input ID to delete")
	flag.Var(&del_outputs, "del-outputs", "Comma-separated list of output ID to delete")
	flag.Var(&del_bonds, "del-bonds", "Comma-separated list of bond ID to delete")
	flag.Var(&add_bond, "add-bond", "Bond with comma-separated endpoints")
	flag.Var(&add_shared_objects, "add-shared-objects", "Add a shared object")
	flag.Var(&del_shared_objects, "del-shared-objects", "Delete a shared object")
	flag.Var(&connect_processor_shared_object, "connect-processor-shared-object", "Connect a processor to a shared object")
	flag.Var(&disconnect_processor_shared_object, "disconnect-processor-shared-object", "Disconnect a processor from a shared object")
	flag.Var(&attach_benchmark_core, "attach-benchmark-core", "Attach a benchmark core")
	flag.Var(&attachBenchmarkCoreV2, "attach-benchmark-core-v2", "Attach a benchmark core v2")

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

func lastAddr(n *net.IPNet) (net.IP, error) { // works when the n is a prefix, otherwise...
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses.")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}
func main() {
	conf := new(bondmachine.Config)
	conf.Debug = *debug
	conf.HwOptimizations = 0
	conf.Dotdetail = uint8(*dot_detail)
	conf.CommentedVerilog = *commentedVerilog
	if *bcofInFile != "" {
		conf.BCOFEntry = new(bcof.BCOFEntry)
		bcofBytes, err := os.ReadFile(*bcofInFile)
		if err != nil {
			panic("failed to read BCOF file")
		}
		if err := proto.Unmarshal(bcofBytes, conf.BCOFEntry); err != nil {
			panic("failed to unmarshal BCOF")
		}
	}

	if *hwOptimizations != "" {
		for _, opt := range strings.Split(*hwOptimizations, ",") {
			if id := procbuilder.HwOptimizationId(opt); id != 0 {
				conf.HwOptimizations = procbuilder.SetHwOptimization(conf.HwOptimizations, id)
			} else {
				panic("Unknown hardware optimization: " + opt)
			}
		}
	}

	var bmach *bondmachine.Bondmachine

	if *bmInfoFile != "" {
		if bmInfoJSON, err := os.ReadFile(*bmInfoFile); err == nil {
			conf.BMinfo = new(bminfo.BMinfo)
			if err := json.Unmarshal(bmInfoJSON, conf.BMinfo); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	if *bmRequirementsFile != "" {
		if bmRequirementsJSON, err := os.ReadFile(*bmRequirementsFile); err == nil {
			reqs := new(bmreqs.ExportedReqs)
			if err := json.Unmarshal(bmRequirementsJSON, reqs); err != nil {
				panic(err)
			}
			newRg, _ := bmreqs.Import(reqs)
			conf.ReqRoot = newRg
			// fmt.Println(newRg.Requirement(bmreqs.ReqRequest{Node: "/", Op: bmreqs.OpDump}))
		} else {
			panic(err)
		}
	}

	if *bondmachine_file != "" {
		if _, err := os.Stat(*bondmachine_file); err == nil {
			// Open the bondmachine file is exists
			if bondmachine_json, err := os.ReadFile(*bondmachine_file); err == nil {
				var bmachj bondmachine.Bondmachine_json
				if err := json.Unmarshal([]byte(bondmachine_json), &bmachj); err == nil {
					bmach = (&bmachj).Dejsoner()
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			// Or create a new one
			bmach = new(bondmachine.Bondmachine)
			bmach.Rsize = uint8(*register_size)
		}

		bmach.Init()

		if &attach_benchmark_core != nil && len(attach_benchmark_core) == 2 {
			err := bmach.Attach_benchmark_core(attach_benchmark_core)
			check(err)
		}

		if &attachBenchmarkCoreV2 != nil && len(attachBenchmarkCoreV2) == 2 {
			err := bmach.AttachBenchmarkCoreV2(attachBenchmarkCoreV2)
			check(err)
		}

		// Eventually create verilog files
		if *create_verilog {
			iomap := new(bondmachine.IOmap)
			if *verilog_mapfile != "" {
				if mapfile_json, err := os.ReadFile(*verilog_mapfile); err == nil {
					if err := json.Unmarshal([]byte(mapfile_json), iomap); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}

			}
			//fmt.Println(iomap)

			// Precess the possible extra modules
			extramodules := make([]bondmachine.ExtraModule, 0)

			// Slower
			if *board_slow {
				em := new(bondmachine.Slow_extra)
				em.Slow_factor = strconv.Itoa(*board_slow_factor)

				if err := em.Check(bmach); err != nil {
					panic(err)
				}
				extramodules = append(extramodules, em)
			}

			// Etherbond
			//TODO
			if *use_etherbond {
				ethb := new(bondmachine.Etherbond_extra)

				config := new(etherbond.Config)
				config.Rsize = uint8(*register_size)

				ethb.Config = config
				ethb.Flavor = *etherbond_flavor

				if *clusterSpec != "" {
					if cluster, err := bmcluster.UnmarshalCluster(*clusterSpec); err != nil {
						panic(err)
					} else {
						ethb.Cluster = cluster
					}
				} else {
					panic("A Cluster spec file is needed")
				}

				ethiomap := new(bondmachine.IOmap)
				if *etherbond_mapfile != "" {
					if mapfile_json, err := os.ReadFile(*etherbond_mapfile); err == nil {
						if err := json.Unmarshal([]byte(mapfile_json), ethiomap); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}

				} else {
					panic(errors.New("Etherbond Mapfile needed"))
				}

				macmap := new(etherbond.Macs)
				if *etherbond_macfile != "" {
					if macfile_json, err := os.ReadFile(*etherbond_macfile); err == nil {
						if err := json.Unmarshal([]byte(macfile_json), macmap); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}

				ethb.Macs = macmap
				ethb.Maps = ethiomap
				ethb.PeerID = uint32(*peerID)
				ethb.Mac = "0288" + fmt.Sprintf("%08d", *peerID)

				if err := ethb.Check(bmach); err != nil {
					panic(err)
				}
				extramodules = append(extramodules, ethb)
			}
			if *use_udpbond {
				udpb := new(bondmachine.Udpbond_extra)

				// TODO Import the wiki configuration from file

				config := new(udpbond.Config)
				config.Rsize = uint8(*register_size)

				udpb.Config = config
				udpb.Flavor = *udpbond_flavor

				if *clusterSpec != "" {
					if cluster, err := udpbond.UnmarshallCluster(config, *clusterSpec); err != nil {
						panic(err)
					} else {
						udpb.Cluster = cluster
					}
				} else {
					panic("A Cluster spec file is needed")
				}

				ethiomap := new(bondmachine.IOmap)
				if *udpbondMapfile != "" {
					if mapfile_json, err := os.ReadFile(*udpbondMapfile); err == nil {
						if err := json.Unmarshal([]byte(mapfile_json), ethiomap); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}

				} else {
					panic(errors.New("Udpbond Mapfile needed"))
				}

				macmap := new(udpbond.Ips)
				if *udpbond_ipfile != "" {
					if macfile_json, err := os.ReadFile(*udpbond_ipfile); err == nil {
						if err := json.Unmarshal([]byte(macfile_json), macmap); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}

				netparams := new(bondmachine.NetParameters)
				if *udpbond_netconfig != "" {
					if netconfig_json, err := os.ReadFile(*udpbond_netconfig); err == nil {
						if err := json.Unmarshal([]byte(netconfig_json), netparams); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}
				fmt.Println(netparams)
				udpb.NetParams = netparams
				udpb.Ips = macmap
				udpb.Maps = ethiomap
				udpb.PeerID = uint32(*peerID)
				if ipst, ok := macmap.Assoc["peer_"+strconv.Itoa(*peerID)]; ok {
					ip := strings.Split(ipst, "/")[0]
					port := strings.Split(ipst, ":")[1]
					udpb.Ip = ip
					udpb.Port = port
					_, nip, _ := net.ParseCIDR(strings.Split(ipst, ":")[0])
					brd, _ := lastAddr(nip)
					udpb.Broadcast = brd.String()
					//udpb.Netmask = nip.String()
				} else {
					panic(errors.New("Wrong IP"))
				}
				if err := udpb.Check(bmach); err != nil {
					panic(err)
				}
				extramodules = append(extramodules, udpb)
			}

			if *useBondirect {
				bdir := new(bondmachine.Bondirect_extra)

				config := new(bondirect.Config)
				config.Rsize = uint8(*register_size)
				bdir.BondirectElement = new(bondirect.BondirectElement)

				bdir.Config = config
				bdir.Flavor = *bondirectFlavor
				bdir.PeerID = uint32(*peerID)

				if *clusterSpec != "" {
					if cluster, err := bmcluster.UnmarshalCluster(*clusterSpec); err != nil {
						panic(err)
					} else {
						bdir.Cluster = cluster
					}
				} else {
					panic("A Cluster spec file is needed")
				}

				if *bondirectMesh != "" {
					if mesh, err := bondirect.UnmarshalMesh(config, *bondirectMesh); err != nil {
						panic(err)
					} else {
						bdir.Mesh = mesh
					}
				} else {
					panic("A Mesh spec file is needed")
				}

				bmdirmap := new(bondmachine.IOmap)
				if *bondirectMapfile != "" {
					if mapfileJSON, err := os.ReadFile(*bondirectMapfile); err == nil {
						if err := json.Unmarshal([]byte(mapfileJSON), bmdirmap); err != nil {
							panic(err)
						} else {
							bdir.Maps = bmdirmap
						}
					} else {
						panic(err)
					}

				} else {
					panic(errors.New("bondirect mapfile needed"))
				}

				// Peer name taken from the mesh
				for _, peer := range bdir.Cluster.Peers {
					if peer.PeerId == bdir.PeerID {
						bdir.PeerName, _ = bdir.AnyNameToMeshName(peer.PeerName)
						break
					}
				}

				clusterNodeName, err := bdir.BondirectElement.AnyNameToClusterName(bdir.PeerName)
				if err != nil {
					panic(err)
				}
				bdir.BondirectElement.InitTData()
				bdir.BondirectElement.PopulateIOData(clusterNodeName)
				bdir.BondirectElement.PopulateWireData(clusterNodeName)

				extramodules = append(extramodules, bdir)
			}

			if *usebmapi {
				bmapi := new(bondmachine.BMAPIExtra)

				bmapi.Rsize = uint8(*register_size)

				// TOTO Error checking
				bmapi.Language = *bmapiLanguage
				bmapi.Flavor = *bmapiFlavor
				bmapi.FlavorVersion = *bmapiFlavorVersion
				bmapi.Framework = *bmapiFramework
				bmapi.LibOutDir = *bmapiLibOutDir
				bmapi.ModOutDir = *bmapiModOutDir
				bmapi.AuxOutDir = *bmapiAuxOutDir
				bmapi.ModuleName = *bmapiModuleName
				bmapi.PackageName = *bmapiPackageName
				bmapi.GenerateExample = *bmapiGenerateExample
				bmapi.DataType = *bmapiDataType

				bmAPIMap := new(bondmachine.IOmap)
				if *bmapiMapfile != "" {
					if mapfileJSON, err := os.ReadFile(*bmapiMapfile); err == nil {
						if err := json.Unmarshal([]byte(mapfileJSON), bmAPIMap); err != nil {
							panic(err)
						} else {
							bmapi.Maps = bmAPIMap
						}
					} else {
						panic(err)
					}

				} else {
					panic(errors.New("BMAPI Mapfile needed"))
				}

				extramodules = append(extramodules, bmapi)

			}

			if *uart {
				uart := new(bondmachine.UartExtra)
				uartMap := new(bondmachine.IOmap)
				if *uartMapFile != "" {
					if mapfileJSON, err := os.ReadFile(*uartMapFile); err == nil {
						if err := json.Unmarshal([]byte(mapfileJSON), uartMap); err != nil {
							panic(err)
						} else {
							uart.Maps = uartMap
						}
					} else {
						panic(err)
					}
				}
				extramodules = append(extramodules, uart)
			}

			if *basys3_7segment {
				b37s := new(bondmachine.B37s)
				b37s.Mapped_output = *basys3_7segment_map
				extramodules = append(extramodules, b37s)
			}

			if *basys3Leds {
				b3l := new(bondmachine.Basys3Leds)
				b3l.MappedOutput = *basys3LedsMap
				b3l.LedName = *basys3LedsName
				if *register_size > 16 {
					panic(errors.New("Basys3 leds can be mapped to a maximum of 16 bits"))
				}
				b3l.Width = strconv.Itoa(*register_size)
				extramodules = append(extramodules, b3l)

			}

			if *counter {
				cnt := new(bondmachine.CounterExtra)
				cnt.MappedInput = *counterMap
				cnt.SlowFactor = *counterSlowFactor
				cnt.Width = strconv.Itoa(*register_size - 1)
				if err := cnt.Check(bmach); err != nil {
					panic(err)
				}
				extramodules = append(extramodules, cnt)
			}

			if *iceBreakerLeds {
				IBL := new(bondmachine.IcebreakerLeds)
				IBL.MappedOutput = *iceBreakerLedsMap
				extramodules = append(extramodules, IBL)
			}

			if *iceFunLeds {
				IFL := new(bondmachine.IceFunLeds)
				IFL.MappedOutput = *iceFunLedsMap
				extramodules = append(extramodules, IFL)
			}

			if *Ice40Lp1kLeds {
				IFL := new(bondmachine.Ice40Lp1kLeds)
				IFL.MappedOutput = *Ice40Lp1kLedsMap
				extramodules = append(extramodules, IFL)
			}

			// Inclusion of PS2 keyboard extra module
			if *ps2keyboardIo {
				ps2 := new(bondmachine.Ps2KeyboardIoExtra)
				ps2.MappedInput = *ps2keyboardIoMap
				extramodules = append(extramodules, ps2)
			}

			// Inclusion of the multi CPU textual VGA module
			if *vgatext {
				if *vgatextFlavor == "800x600" {
					vgat := new(bondmachine.Vga800x600Extra)
					vgat.Header = *vgatextHeader
					vgat.Fonts = *vgatextFonts
					extramodules = append(extramodules, vgat)
				} else {
					panic(errors.New("Unsopported VGA flavor"))
				}
			}

			var flavor string

			if *verilog_simulation {
				flavor = *verilog_flavor + "_simulation"
			} else {
				flavor = *verilog_flavor
			}

			var sbox *simbox.Simbox

			if *verilog_simulation {
				if *simbox_file != "" {
					sbox = new(simbox.Simbox)
					if _, err := os.Stat(*simbox_file); err == nil {
						// Open the simbox file is exists
						if simbox_json, err := os.ReadFile(*simbox_file); err == nil {
							if err := json.Unmarshal([]byte(simbox_json), sbox); err != nil {
								panic(err)
							}
						} else {
							panic(err)
						}
					}

				}
			}

			bmach.Write_verilog(conf, flavor, iomap, extramodules, sbox)

			if *usebmapi {
				if err := bmach.WriteBMAPI(conf, flavor, iomap, extramodules, sbox); err != nil {
					panic(err)
				}
			}
		}

		// All the operation are exclusive
		if *list_domains {
			fmt.Println(bmach.List_domains())
		} else if &add_domains != nil && len(add_domains) != 0 {
			for _, load_machine := range add_domains {
				if _, err := os.Stat(load_machine); err == nil {
					if jsonfile, err := os.ReadFile(load_machine); err == nil {
						var machj procbuilder.Machine_json
						if err := json.Unmarshal([]byte(jsonfile), &machj); err == nil {
							mymachine := (&machj).Dejsoner()
							bmach.Domains = append(bmach.Domains, mymachine)
						} else {
							panic(err)
						}
					} else {
						panic(err)
					}
				} else {
					fmt.Println(load_machine + " file not found, ignoring it.")
				}
			}
		} else if (&del_domains != nil) && len(del_domains) != 0 {
			for _, remove_domain := range del_domains {
				if remove_domain_id, err := strconv.Atoi(remove_domain); err == nil {
					if remove_domain_id < len(bmach.Domains)-1 {
						bmach.Domains = append(bmach.Domains[:remove_domain_id], bmach.Domains[remove_domain_id+1:]...)
					} else if remove_domain_id == len(bmach.Domains)-1 {
						bmach.Domains = bmach.Domains[:remove_domain_id]
					} else {
						fmt.Println(remove_domain + " not a valid domain id, ignoring it.")
					}
				} else {
					fmt.Println(remove_domain + " not a valid domain id, ignoring it.")
				}
			}
			// TODO Include the check of unbounded processors
		} else if *specs {
			fmt.Printf(bmach.Specs())
		} else if *list_inputs {
			for i, inp := range bmach.List_inputs() {
				fmt.Printf("%d %s\n", i, inp)
			}
		} else if *add_inputs != 0 {
			for i := 0; i < *add_inputs; i++ {
				message, err := bmach.Add_input()
				check(err)
				if *debug {
					log.Println(message)
				} else if *verbose {
					fmt.Println(message)
				}
			}
		} else if (&del_inputs != nil) && len(del_inputs) != 0 {
			// Reorder the inputs to delete, last first
			todelete := make([]int, 0)
			for _, inp := range del_inputs {
				if value, ok := strconv.Atoi(inp); ok == nil {
					pcheck := false
					for _, i := range todelete {
						if i == value {
							pcheck = true
							break
						}
					}
					if !pcheck && value < bmach.Inputs {
						todelete = append(todelete, value)
					}
				}
			}
			sort.Ints(todelete)
			for i, _ := range todelete {
				// Remove the inputs, higher first
				bmach.Del_input(todelete[len(todelete)-i-1])
			}
		} else if *list_outputs {
			for i, outp := range bmach.List_outputs() {
				fmt.Printf("%d %s\n", i, outp)
			}
		} else if *add_outputs != 0 {
			for i := 0; i < *add_outputs; i++ {
				message, err := bmach.Add_output()
				check(err)
				if *debug {
					log.Println(message)
				} else if *verbose {
					fmt.Println(message)
				}
			}
		} else if (&del_outputs != nil) && len(del_outputs) != 0 {
			// Reorder the outputs to delete, last first
			todelete := make([]int, 0)
			for _, inp := range del_outputs {
				if value, ok := strconv.Atoi(inp); ok == nil {
					pcheck := false
					for _, i := range todelete {
						if i == value {
							pcheck = true
							break
						}
					}
					if !pcheck && value < bmach.Outputs {
						todelete = append(todelete, value)
					}
				}
			}
			sort.Ints(todelete)
			for i, _ := range todelete {
				// Remove the outputs, higher first
				bmach.Del_output(todelete[len(todelete)-i-1])
			}
		} else if *enumProcessors {
			fmt.Println(bmach.EnumProcessors())
		} else if *enumBonds {
			fmt.Println(bmach.EnumBonds())
		} else if *list_processors {
			fmt.Print(bmach.List_processors())
		} else if *add_processor != -1 {
			message, err := bmach.Add_processor(*add_processor)
			check(err)
			fmt.Println(message)
		} else if *list_bonds {
			for i, bond := range bmach.List_bonds() {
				fmt.Printf("%d %s\n", i, bond)
			}
		} else if *list_internal_inputs {
			for _, inp := range bmach.List_internal_inputs() {
				fmt.Println(inp)
			}
		} else if *list_internal_outputs {
			for _, outp := range bmach.List_internal_outputs() {
				fmt.Println(outp)
			}
		} else if *emit_dot && !*sim {
			fmt.Print(bmach.Dot(conf, "", nil, nil))
		} else if *show_program_disassembled {
			// TODO Finish
		} else if *multi_abstract_assembly_file != "" {
			// TODO Temporary, clean up code!
			mu, _ := bmach.GetMultiAssembly()
			// Write the multi_abstract_assembly_file file
			mufile, err := os.Create(*multi_abstract_assembly_file)
			check(err)
			defer mufile.Close()
			b, errj := json.Marshal(mu)
			check(errj)
			_, err = mufile.WriteString(string(b))
			check(err)
		} else if *show_program_alias {
			pa, _ := bmach.GetProgramsAlias()
			for i, al := range pa {
				alfile, err := os.Create("p" + strconv.Itoa(i) + ".alias")
				check(err)
				defer alfile.Close()
				_, err = alfile.WriteString(string(al))
				check(err)
			}
		} else if &add_bond != nil && len(add_bond) == 2 {
			bmach.Add_bond(add_bond)
		} else if (&del_bonds != nil) && len(del_bonds) != 0 {
			for _, remove_bond := range del_bonds {
				if remove_bond_id, err := strconv.Atoi(remove_bond); err == nil {
					if remove_bond_id < len(bmach.Links) {
						bmach.Del_bond(remove_bond_id)
					} else {
						fmt.Println(remove_bond + " not a valid bond id, ignoring it.")
					}
				} else {
					fmt.Println(remove_bond + " not a valid bond id, ignoring it.")
				}
			}
		} else if *list_shared_objects {
			fmt.Print(bmach.List_shared_objects())
		} else if &add_shared_objects != nil && len(add_shared_objects) > 0 {
			bmach.Add_shared_objects(add_shared_objects)
		} else if *list_processor_shared_object_links {
			fmt.Print(bmach.List_processor_shared_object_links())
		} else if &connect_processor_shared_object != nil && len(connect_processor_shared_object) == 2 {
			bmach.Connect_processor_shared_object(connect_processor_shared_object)
		} else if *sim {
			var sbox *simbox.Simbox
			if *simbox_file != "" {
				sbox = new(simbox.Simbox)
				if _, err := os.Stat(*simbox_file); err == nil {
					// Open the simbox file is exists
					if simbox_json, err := os.ReadFile(*simbox_file); err == nil {
						if err := json.Unmarshal([]byte(simbox_json), sbox); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}

			}

			// Build the simulation VM
			vm := new(bondmachine.VM)
			vm.Bmach = bmach
			err := vm.Init()
			check(err)

			var pstatevm *bondmachine.VM

			// Build the simulation configuration
			sconfig := new(bondmachine.SimConfig)
			scerr := sconfig.Init(sbox, vm, conf)
			check(scerr)

			// Build the simulation driver
			sdrive := new(bondmachine.SimDrive)
			sderr := sdrive.Init(conf, sbox, vm)
			check(sderr)

			// Build the simultion report
			srep := new(bondmachine.SimReport)
			srerr := srep.Init(sbox, vm)
			check(srerr)

			lerr := vm.Launch_processors(sbox)
			check(lerr)

			var intlen_s string

			if *emit_dot {
				pstatevm = new(bondmachine.VM)
				pstatevm.Bmach = bmach
				err := pstatevm.Init()
				check(err)

				sim_int_s := strconv.Itoa(*simInteractions)
				intlen := len(sim_int_s)
				intlen_s = strconv.Itoa(intlen)
			}

			var reportData *csv.Writer
			if *simReport != "" {
				rf, err := os.Create(*simReport)
				check(err)
				defer rf.Close()
				reportData = csv.NewWriter(rf)

				if sconfig.GetTicks {
					if err := reportData.Write(append([]string{"tick"}, srep.ReportablesNames...)); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
				} else {
					if err := reportData.Write(srep.ReportablesNames); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
				}

				reportData.Flush()

				if err := reportData.Error(); err != nil {
					log.Fatal(err)
				}

			}

			var oldRecordC *[]string

			if *simStopOnValidOf != -1 {
				if *simStopOnValidOf >= len(vm.OutputsValid) {
					log.Fatal("sim-stop-on-valid-of index out of range")
				}
			}

			// Main simulation loop, tick by tick
			for i := uint64(0); i < uint64(*simInteractions); i++ {

				if *simStopOnValidOf != -1 {
					if vm.OutputsValid[*simStopOnValidOf] {
						if *debug {
							log.Printf("Stopping simulation at tick %d due to sim-stop-on-valid-of\n", i)
						}
						break
					}
				}

				// Manage the valid/recv states of the inputs
				for inIdx, inRecv := range vm.InputsRecv {
					if inRecv {
						vm.InputsValid[inIdx] = false
					}
				}

				// This will get actions eventually to do on this tick
				if act, exist_actions := sdrive.AbsSet[i]; exist_actions {
					for k, val := range act {
						*sdrive.Injectables[k] = val
						if inIdx, ok := sdrive.NeedValid[k]; ok {
							vm.InputsValid[inIdx] = true
						}
					}
				}

				// TODO Periodic set

				if *emit_dot {
					gvfile := bmach.Dot(conf, "", vm, pstatevm)
					filename := fmt.Sprintf("graphviz%0"+intlen_s+"d", int(i))
					f, err := os.Create(filename + ".dot")
					check(err)
					_, err = f.WriteString(gvfile)
					check(err)
					f.Close()

					pstatevm.CopyState(vm)
				}

				result, err := vm.Step(sconfig)
				check(err)

				// Manage the valid/recv states of the outputs
				for outIdx, outValid := range vm.OutputsValid {
					if outValid {
						vm.OutputsRecv[outIdx] = true
					} else {
						vm.OutputsRecv[outIdx] = false
					}
				}

				fmt.Print(result)

				showList := make([]int, 0, len(srep.Showables))

				// This will get value to show on this tick
				if slist, exist_shows := srep.AbsShow[i]; exist_shows {
					for k, _ := range slist {
						showList = append(showList, k)
					}
				}

				// This will get value to show on periodic ticks
				for j, slist := range srep.PerShow {
					if i%j == 0 {
						for k, _ := range slist {
							alredtIn := false
							for _, v := range showList {
								if v == k {
									alredtIn = true
									break
								}
							}
							if !alredtIn {
								showList = append(showList, k)
							}
						}
					}
				}

				sort.Ints(showList)

				// Show the tick values
				for _, k := range showList {
					nType := srep.ShowablesTypes[k]
					if _, err := bmnumbers.EventuallyCreateType(nType, nil); err != nil {
						log.Fatal(err)
					}
					if v := bmnumbers.GetType(nType); v == nil {
						log.Fatal("Error: Unknown type")
					} else {
						bits := v.GetSize()
						if number, err := bmnumbers.ImportUint(*srep.Showables[k], bits); err != nil {
							log.Fatal(err)
						} else {
							if err := bmnumbers.CastType(number, v); err != nil {
								log.Fatal(err)
							} else {
								if numberS, err := number.ExportString(nil); err != nil {
									log.Fatal(err)
								} else {
									fmt.Print(numberS + " ")
								}
							}
						}
					}
				}
				if len(showList) > 0 {
					fmt.Println()
				}

				if *simReport != "" {

					repList := make([]int, 0, len(srep.Reportables))
					recordC := make([]string, len(srep.Reportables))

					if sconfig.GetTicks {
						recordC = append(recordC, "")
						recordC[0] = strconv.FormatUint(i, 10)
					}

					if sconfig.GetAll || sconfig.GetAllInternal {
						for j := range srep.Reportables {
							repList = append(repList, j)
						}
					} else {
						if rep, exist_reports := srep.AbsGet[i]; exist_reports {
							for k := range rep {
								repList = append(repList, k)
							}
						}

						// This will get value to show on periodic ticks
						for j, rep := range srep.PerGet {
							if i%j == 0 {
								for k, _ := range rep {
									alredtIn := false
									for _, v := range repList {
										if v == k {
											alredtIn = true
											break
										}
									}
									if !alredtIn {
										repList = append(repList, k)
									}
								}
							}
						}
					}

					// sort.Ints(repList)

					someToReport := false
					for _, k := range repList {
						nType := srep.ReportablesTypes[k]
						if _, err := bmnumbers.EventuallyCreateType(nType, nil); err != nil {
							log.Fatal(err)
						}
						if v := bmnumbers.GetType(nType); v == nil {
							log.Fatal("Error: Unknown type")
						} else {
							bits := v.GetSize()
							if number, err := bmnumbers.ImportUint(*srep.Reportables[k], bits); err != nil {
								log.Fatal(err)
							} else {
								if err := bmnumbers.CastType(number, v); err != nil {
									log.Fatal(err)
								} else {
									if numberS, err := number.ExportString(nil); err != nil {
										log.Fatal(err)
									} else {
										someToReport = true
										if sconfig.GetTicks {
											recordC[k+1] = numberS
										} else {
											recordC[k] = numberS
										}
									}
								}
							}
						}
					}

					if sconfig.GetTicks || someToReport {

						reportWrite := make([]string, len(recordC))

						for recI, recV := range recordC {
							if sconfig.GetTicks {
								if conf.FormatSimReports && (oldRecordC == nil || ((recI != 0) && (*oldRecordC)[recI] != recV)) {
									reportWrite[recI] = "\033[31m" + fmt.Sprintf("%-25s", recV) + "\033[0m"
								} else if conf.FormatSimReports {
									reportWrite[recI] = fmt.Sprintf("%-25s", recV)
								} else {
									reportWrite[recI] = recV
								}
							} else {
								if conf.FormatSimReports && (oldRecordC == nil || (*oldRecordC)[recI] != recV) {
									reportWrite[recI] = "\033[31m" + recV + "\033[0m"
								} else {
									reportWrite[recI] = recV
								}
							}
						}

						if err := reportData.Write(reportWrite); err != nil {
							log.Fatalln("error writing record to csv:", err)
						}

						reportData.Flush()

						if err := reportData.Error(); err != nil {
							log.Fatal(err)
						}
						oldRecordC = &recordC
					}
				}
			}
		} else if *emu {

			emuDrivers := make([]bondmachine.EmuDriver, 0)

			// Inclusion of the multi CPU textual VGA module
			if *vgatext {
				if *vgatextFlavor == "800x600" {
					vga := new(bondmachine.Vga800x600Emu)
					emuDrivers = append(emuDrivers, vga)
				} else {
					panic(errors.New("unsopported VGA flavor"))
				}
			}

			vm := new(bondmachine.VM)
			vm.Bmach = bmach
			vm.EmuDrivers = emuDrivers
			err := vm.Init()
			check(err)

			// the emulation configuration is not really needed
			sconfig := new(bondmachine.SimConfig)
			scerr := sconfig.Init(nil, vm, conf)
			check(scerr)

			lerr := vm.Launch_processors(nil)
			check(lerr)

			for i := uint64(0); ; {
				if *emu_interactions != 0 {
					if i >= uint64(*emu_interactions) {
						break
					}
				}

				// Manage the valid/recv states of the inputs
				for inIdx, inRecv := range vm.InputsRecv {
					if inRecv {
						vm.InputsValid[inIdx] = false
					}
				}

				_, err := vm.Step(sconfig)
				check(err)

				// Manage the valid/recv states of the outputs
				for outIdx, outValid := range vm.OutputsValid {
					if outValid {
						vm.OutputsRecv[outIdx] = true
					} else {
						vm.OutputsRecv[outIdx] = false
					}
				}

				if *emu_interactions != 0 {
					if *verbose {
						fmt.Println("Interaction:", i)
					}
					i++
				}

			}
		}
		// Write the bondmachine file
		f, err := os.Create(*bondmachine_file)
		check(err)
		defer f.Close()
		b, errj := json.Marshal(bmach.Jsoner())
		check(errj)
		_, err = f.WriteString(string(b))
		check(err)
	} else {
		panic("bondmachine-file is a mandatory option")
	}
}
