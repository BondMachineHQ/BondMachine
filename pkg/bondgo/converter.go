package bondgo

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
	"github.com/BondMachineHQ/BondMachine/pkg/udpbond"
)

func (bg *BondgoCheck) Create_Connecting_Processor(rsize int, procid int) (*procbuilder.Machine, bool) {
	mymachine := new(procbuilder.Machine)

	myarch := &mymachine.Arch

	myarch.Rsize = uint8(rsize)

	myarch.Modes = make([]string, 1)
	myarch.Modes[0] = "ha"

	preq := bg.Procr[procid]

	opcodes := make([]procbuilder.Opcode, 0)

	for _, op := range procbuilder.Allopcodes {
		for _, opn := range preq.Opcodes {
			if opn == op.Op_get_name() {
				opcodes = append(opcodes, op)
				break
			}
		}
	}

	sort.Sort(procbuilder.ByName(opcodes))

	myarch.Op = opcodes

	myarch.R = uint8(Needed_bits(preq.Registersize))
	myarch.L = uint8(Needed_bits(preq.Ramsize))
	myarch.N = uint8(preq.Inputs)
	myarch.M = uint8(preq.Outputs)
	myarch.O = uint8(Needed_bits(preq.Romsize))
	myarch.Shared_constraints = strings.Join(preq.SharedObjects, ",")

	prog := bg.Write_assembly(procid)

	//fmt.Println(prog)

	if prog, err := myarch.Assembler([]byte(prog)); err == nil {
		mymachine.Program = prog
	} else {
		fmt.Println(err)
	}

	return mymachine, true
}

func (bg *BondgoCheck) Create_Bondmachine(rsize int, filter string) (*bondmachine.Bondmachine, *bondmachine.Residual, error) {
	bmach := new(bondmachine.Bondmachine)

	bmach.Rsize = uint8(rsize)

	bmach.Init()

	if bg.Cascading_io {
		if _, ok := bmach.Add_input(); ok != nil {
			return nil, nil, errors.New("Input add failed")
		}
		if _, ok := bmach.Add_output(); ok != nil {
			fmt.Println("Output add failed")
			return nil, nil, errors.New("Output add failed")
		}
	}

	// TODO Remove the fucking len with something smarter
	proc_id := 0
	rel_procids := make(map[int]int)

	for i := 0; i < len(bg.Program); i++ {
		if bg.Procr[i].Device == filter {
			if mymachine, ok := bg.Create_Connecting_Processor(rsize, i); ok {
				bmach.Domains = append(bmach.Domains, mymachine)
				if _, ok := bmach.Add_processor(proc_id); ok != nil {
					return nil, nil, errors.New("Attach processor failed")
				}
				rel_procids[i] = proc_id
				if bg.Cascading_io {
					if i == 0 {
						bmach.Add_bond([]string{"i0", "p0i0"})
					} else if i == len(bg.Program)-1 {
						bmach.Add_bond([]string{"p" + strconv.Itoa(i) + "i0", "p" + strconv.Itoa(i-1) + "o0"})
						bmach.Add_bond([]string{"o0", "p" + strconv.Itoa(i) + "o0"})
					} else {
						bmach.Add_bond([]string{"p" + strconv.Itoa(i) + "i0", "p" + strconv.Itoa(i-1) + "o0"})
					}
				}
				proc_id += 1
			} else {
				return nil, nil, errors.New("Creating processor failed")
			}
		}
	}

	// Creation of the IO topology

	unconnected_inputs := make(map[int]bool)

	ext_input := make(map[int]string)
	ext_input_keys := make([]int, 0)
	ext_output := make(map[int]string)
	ext_output_keys := make([]int, 0)

	for i, ioref := range bg.IOr {
		if bg.Procr[i].Device == filter {
			for _, in_id := range ioref.Inputs_ids {
				if _, ok := unconnected_inputs[in_id]; !ok {
					unconnected_inputs[in_id] = true
				}
			}
		}
	}

	for proc1_id, ioref1 := range bg.IOr {
		if bg.Procr[proc1_id].Device == filter {
			for out_ord, out_id := range ioref1.Outputs_ids {
				connected := false
				for proc2_id, ioref2 := range bg.IOr {
					if bg.Procr[proc2_id].Device == filter {
						for in_ord, in_id := range ioref2.Inputs_ids {
							if out_id == in_id {
								unconnected_inputs[in_id] = false
								connected = true
								bmach.Add_bond([]string{"p" + strconv.Itoa(rel_procids[proc2_id]) + "i" + strconv.Itoa(in_ord), "p" + strconv.Itoa(rel_procids[proc1_id]) + "o" + strconv.Itoa(out_ord)})
								//fmt.Println([]string{"p" + strconv.Itoa(proc2_id) + "i" + strconv.Itoa(in_ord), "p" + strconv.Itoa(proc1_id) + "o" + strconv.Itoa(out_ord)})
							}
						}
					}
				}
				if !connected {
					ext_output[out_id] = "p" + strconv.Itoa(rel_procids[proc1_id]) + "o" + strconv.Itoa(out_ord)
					ext_output_keys = append(ext_output_keys, out_id)
				}
			}
		}
	}

	for in_id, id_unconn := range unconnected_inputs {
		if id_unconn {
			for proc_id, ioref := range bg.IOr {
				for in_ord, in_idc := range ioref.Inputs_ids {
					if in_id == in_idc {
						ext_input[in_id] = "p" + strconv.Itoa(rel_procids[proc_id]) + "i" + strconv.Itoa(in_ord)
						ext_input_keys = append(ext_input_keys, in_id)
						//bmach.Add_bond([]string{"p" + strconv.Itoa(proc_id) + "i" + strconv.Itoa(in_ord), "i" + strconv.Itoa(bmach.Inputs-1)})
					}
				}
			}
		}
	}

	residual := new(bondmachine.Residual)
	residual.Map.Assoc = make(map[string]string)

	sort.Ints(ext_input_keys)
	sort.Ints(ext_output_keys)
	//fmt.Println(ext_input_keys, ext_output_keys)
	for _, inp := range ext_input_keys {
		inps := ext_input[inp]
		if _, ok := bmach.Add_input(); ok != nil {
			return nil, nil, errors.New("Input add failed")
		} else {
			bmach.Add_bond([]string{inps, "i" + strconv.Itoa(bmach.Inputs-1)})
			residual.Map.Assoc["i"+strconv.Itoa(bmach.Inputs-1)] = strconv.Itoa(inp)
		}
	}

	for _, out := range ext_output_keys {
		outs := ext_output[out]
		if _, ok := bmach.Add_output(); ok != nil {
			return nil, nil, errors.New("Output add failed")
		} else {
			bmach.Add_bond([]string{outs, "o" + strconv.Itoa(bmach.Outputs-1)})
			residual.Map.Assoc["o"+strconv.Itoa(bmach.Outputs-1)] = strconv.Itoa(out)
		}
	}

	// TODO Only channels here, include the others
	creqs := bg.Chanr

	for _, _ = range creqs {
		bmach.Add_shared_objects([]string{"channel:"})
	}
	for chanid, creq := range creqs {
		for _, proc_id := range creq.Connected {
			endpoints := make([]string, 2)
			endpoints[0] = strconv.Itoa(proc_id)
			endpoints[1] = strconv.Itoa(chanid)
			bmach.Connect_processor_shared_object(endpoints)
		}
	}

	return bmach, residual, nil
}

func (bg *BondgoRequirements) Abstract_assembler(rsize int, asmcode []string, used chan UsageNotify) error {
	for _, line := range asmcode {
		words := make([]string, 0)
		for _, word := range strings.Split(line, " ") {
			if word != "" {
				words = append(words, strings.TrimSpace(word))
			}
		}
		if len(words) != 0 {
			if words[0][0] != '#' {
				for _, op := range procbuilder.Allopcodes {
					if op.Op_get_name() == words[0] {
						if results, err := op.AbstractAssembler(nil, words[1:]); err == nil {
							for _, result := range results {
								if result.ComponentType == C_OUTPUT {
									// TODO Very temporary code
									result.Componenti += 100
								}
								used <- UsageNotify{TR_PROC, 0, result.ComponentType, result.Components, result.Componenti}
							}
						} else {
							return errors.New(err.Error() + ", error processing " + op.Op_get_name())
						}
					}
				}
			} else {
				if words[0] == "#archinclude" {
					// TODO Finish
					if len(words) > 1 {
						for _, word := range words[1:len(words)] {
							seq0, types0 := Sequence_to_0(word)

							switch types0 {
							case C_INPUT:

								for i, _ := range seq0 {
									used <- UsageNotify{TR_PROC, 0, C_INPUT, S_NIL, i + 1}
								}
							}
						}
					}
				}
			}
		} else {
			return nil
		}
	}
	return nil
}

func (bg *BondgoCheck) Create_Etherbond_Cluster(rsize int, extc *bmcluster.Cluster) (*bmcluster.Cluster, []uint32, []*bondmachine.Bondmachine, []*bondmachine.IOmap, []*bondmachine.Residual, error) {
	devlist := make([]string, 0)

	resultbond := make([]*bondmachine.Bondmachine, 0)
	resultpeerids := make([]uint32, 0)
	resultio := make([]*bondmachine.IOmap, 0)
	resultresi := make([]*bondmachine.Residual, 0)

	partialio := make([]*bondmachine.Residual, 0)

	for i := 0; i < len(bg.Program); i++ {
		check := false
		currdev := bg.Procr[i].Device
		for _, idev := range devlist {
			if idev == currdev {
				check = true
				break
			}
		}
		if !check {
			devlist = append(devlist, bg.Procr[i].Device)
		}
	}

	for _, device := range devlist {
		if bmach, residual, err := bg.Create_Bondmachine(rsize, device); err == nil {
			resultbond = append(resultbond, bmach)
			partialio = append(partialio, residual)
		} else {
			return nil, nil, nil, nil, nil, errors.New("BondMachine creation failed")
		}
	}

	for _, _ = range resultbond {
		// Find a free peerid
	nextpid:
		for nexpid := uint32(1); nexpid < MAXPEERID; nexpid++ {
			if extc != nil {
				for _, peer := range extc.Peers {
					if peer.PeerId == nexpid {
						continue nextpid
					}
				}
			}

			for _, i := range resultpeerids {
				if i == nexpid {
					continue nextpid
				}
			}

			resultpeerids = append(resultpeerids, nexpid)

			break
		}
	}

	resultcluster := new(bmcluster.Cluster)
	resultcluster.ClusterId = uint32(0)
	resultcluster.Peers = make([]bmcluster.Peer, 0)

	if extc != nil {

		// Precess every peer of the extra cluster, its peers connot be residual
		for _, peer := range extc.Peers {
			newinputs := make([]uint32, 0)
			newoutputs := make([]uint32, 0)
			newchannels := make([]uint32, 0)

			// TODO Include the channels whenever their implementation will be ready

			for _, inp := range peer.Inputs {
				// TODO For now the external cluater will be considered checked, in the future a real check will be desiderable
				newinputs = append(newinputs, inp)
			}

			for _, outp := range peer.Outputs {
				newoutputs = append(newoutputs, outp)
			}
			newpeer := bmcluster.Peer{peer.PeerId, "", newchannels, newinputs, newoutputs}
			resultcluster.Peers = append(resultcluster.Peers, newpeer)
		}
	}

	// Process the created bondmachines
	for myid, res := range partialio {
		newio := new(bondmachine.IOmap)
		newio.Assoc = make(map[string]string)
		newresidual := new(bondmachine.Residual)
		newresidual.Map.Assoc = make(map[string]string)

		newinputs := make([]uint32, 0)
		newoutputs := make([]uint32, 0)
		newchannels := make([]uint32, 0)

		for bio, ioid := range res.Map.Assoc {
			if bio[0] == 'i' {
				connected := false
				multi := false

				// Search the input among ext cluster outputs
				if extc != nil {
				extloop:
					for _, peer := range extc.Peers {
						for _, outp := range peer.Outputs {
							if strconv.Itoa(int(outp)) == ioid {
								if connected {
									multi = true
									break extloop
								} else {
									connected = true
								}
							}
						}
					}
				}
				// Search the input among other bondmachines
			bmloop:
				for otherid, otherres := range partialio {
					// Not searching on the same BM
					if otherid != myid {
						for otherbio, otherioid := range otherres.Map.Assoc {
							if ioid == otherioid {
								if otherbio[0] == 'i' {
									// That's ok to have other inputs connected to the same output
								} else if otherbio[0] == 'o' {
									if connected {
										multi = true
										break bmloop
									} else {
										connected = true
									}

								}
							}
						}
					}
				}

				if multi {
					return nil, nil, nil, nil, nil, errors.New("Input connected to multiple outputs")
				}

				if connected {
					// The Input is within the cluster
					newio.Assoc[bio] = ioid
					ioid_n, _ := strconv.Atoi(ioid)
					newinputs = append(newinputs, uint32(ioid_n))
				} else {
					// Its residual
					newresidual.Map.Assoc[bio] = ioid
				}

			} else if bio[0] == 'o' {
				connected := false
				multioutput := false

				// Search the output among ext cluster inputs
				if extc != nil {
					for _, peer := range extc.Peers {
						for _, inp := range peer.Inputs {
							if strconv.Itoa(int(inp)) == ioid {
								connected = true
								break
							}
						}

						for _, outp := range peer.Outputs {
							if strconv.Itoa(int(outp)) == ioid {
								multioutput = true
								break
							}
						}
					}
				}

				for otherid, otherres := range partialio {
					// Not searching on the same BM
					if otherid != myid {
						for otherbio, otherioid := range otherres.Map.Assoc {
							if ioid == otherioid {
								if otherbio[0] == 'o' {
									multioutput = true
								} else if otherbio[0] == 'i' {
									connected = true
								}
							}
						}
					}
				}

				if multioutput {
					return nil, nil, nil, nil, nil, errors.New("Multiple outputs have the same id")
				}

				if connected {
					// The Output is within the cluster
					newio.Assoc[bio] = ioid
					ioid_n, _ := strconv.Atoi(ioid)
					newoutputs = append(newoutputs, uint32(ioid_n))
				} else {
					// Its residual
					newresidual.Map.Assoc[bio] = ioid
				}

			}
		}

		resultio = append(resultio, newio)
		resultresi = append(resultresi, newresidual)

		newpeer := bmcluster.Peer{resultpeerids[myid], "", newchannels, newinputs, newoutputs}
		resultcluster.Peers = append(resultcluster.Peers, newpeer)

	}

	return resultcluster, resultpeerids, resultbond, resultio, resultresi, nil
}

func (bg *BondgoCheck) Create_Udpbond_Cluster(rsize int, extc *udpbond.Cluster) (*udpbond.Cluster, []uint32, []*bondmachine.Bondmachine, []*bondmachine.IOmap, []*bondmachine.Residual, error) {
	devlist := make([]string, 0)

	resultbond := make([]*bondmachine.Bondmachine, 0)
	resultpeerids := make([]uint32, 0)
	resultio := make([]*bondmachine.IOmap, 0)
	resultresi := make([]*bondmachine.Residual, 0)

	partialio := make([]*bondmachine.Residual, 0)

	for i := 0; i < len(bg.Program); i++ {
		check := false
		currdev := bg.Procr[i].Device
		for _, idev := range devlist {
			if idev == currdev {
				check = true
				break
			}
		}
		if !check {
			devlist = append(devlist, bg.Procr[i].Device)
		}
	}

	for _, device := range devlist {
		if bmach, residual, err := bg.Create_Bondmachine(rsize, device); err == nil {
			resultbond = append(resultbond, bmach)
			partialio = append(partialio, residual)
		} else {
			return nil, nil, nil, nil, nil, errors.New("BondMachine creation failed")
		}
	}

	for _, _ = range resultbond {
		// Find a free peerid
	nextpid:
		for nexpid := uint32(1); nexpid < MAXPEERID; nexpid++ {
			if extc != nil {
				for _, peer := range extc.Peers {
					if peer.PeerId == nexpid {
						continue nextpid
					}
				}
			}

			for _, i := range resultpeerids {
				if i == nexpid {
					continue nextpid
				}
			}

			resultpeerids = append(resultpeerids, nexpid)

			break
		}
	}

	resultcluster := new(udpbond.Cluster)
	resultcluster.ClusterId = uint32(0)
	resultcluster.Peers = make([]udpbond.Peer, 0)

	if extc != nil {

		// Precess every peer of the extra cluster, its peers connot be residual
		for _, peer := range extc.Peers {
			newinputs := make([]uint32, 0)
			newoutputs := make([]uint32, 0)
			newchannels := make([]uint32, 0)

			// TODO Include the channels whenever their implementation will be ready

			for _, inp := range peer.Inputs {
				// TODO For now the external cluater will be considered checked, in the future a real check will be desiderable
				newinputs = append(newinputs, inp)
			}

			for _, outp := range peer.Outputs {
				newoutputs = append(newoutputs, outp)
			}
			newpeer := udpbond.Peer{peer.PeerId, newchannels, newinputs, newoutputs}
			resultcluster.Peers = append(resultcluster.Peers, newpeer)
		}
	}

	// Process the created bondmachines
	for myid, res := range partialio {
		newio := new(bondmachine.IOmap)
		newio.Assoc = make(map[string]string)
		newresidual := new(bondmachine.Residual)
		newresidual.Map.Assoc = make(map[string]string)

		newinputs := make([]uint32, 0)
		newoutputs := make([]uint32, 0)
		newchannels := make([]uint32, 0)

		for bio, ioid := range res.Map.Assoc {
			if bio[0] == 'i' {
				connected := false
				multi := false

				// Search the input among ext cluster outputs
				if extc != nil {
				extloop:
					for _, peer := range extc.Peers {
						for _, outp := range peer.Outputs {
							if strconv.Itoa(int(outp)) == ioid {
								if connected {
									multi = true
									break extloop
								} else {
									connected = true
								}
							}
						}
					}
				}
				// Search the input among other bondmachines
			bmloop:
				for otherid, otherres := range partialio {
					// Not searching on the same BM
					if otherid != myid {
						for otherbio, otherioid := range otherres.Map.Assoc {
							if ioid == otherioid {
								if otherbio[0] == 'i' {
									// That's ok to have other inputs connected to the same output
								} else if otherbio[0] == 'o' {
									if connected {
										multi = true
										break bmloop
									} else {
										connected = true
									}

								}
							}
						}
					}
				}

				if multi {
					return nil, nil, nil, nil, nil, errors.New("Input connected to multiple outputs")
				}

				if connected {
					// The Input is within the cluster
					newio.Assoc[bio] = ioid
					ioid_n, _ := strconv.Atoi(ioid)
					newinputs = append(newinputs, uint32(ioid_n))
				} else {
					// Its residual
					newresidual.Map.Assoc[bio] = ioid
				}

			} else if bio[0] == 'o' {
				connected := false
				multioutput := false

				// Search the output among ext cluster inputs
				if extc != nil {
					for _, peer := range extc.Peers {
						for _, inp := range peer.Inputs {
							if strconv.Itoa(int(inp)) == ioid {
								connected = true
								break
							}
						}

						for _, outp := range peer.Outputs {
							if strconv.Itoa(int(outp)) == ioid {
								multioutput = true
								break
							}
						}
					}
				}

				for otherid, otherres := range partialio {
					// Not searching on the same BM
					if otherid != myid {
						for otherbio, otherioid := range otherres.Map.Assoc {
							if ioid == otherioid {
								if otherbio[0] == 'o' {
									multioutput = true
								} else if otherbio[0] == 'i' {
									connected = true
								}
							}
						}
					}
				}

				if multioutput {
					return nil, nil, nil, nil, nil, errors.New("Multiple outputs have the same id")
				}

				if connected {
					// The Output is within the cluster
					newio.Assoc[bio] = ioid
					ioid_n, _ := strconv.Atoi(ioid)
					newoutputs = append(newoutputs, uint32(ioid_n))
				} else {
					// Its residual
					newresidual.Map.Assoc[bio] = ioid
				}

			}
		}

		resultio = append(resultio, newio)
		resultresi = append(resultresi, newresidual)

		newpeer := udpbond.Peer{resultpeerids[myid], newchannels, newinputs, newoutputs}
		resultcluster.Peers = append(resultcluster.Peers, newpeer)

	}

	return resultcluster, resultpeerids, resultbond, resultio, resultresi, nil
}

func Assembly_2_Processor(rsize int, asmcode []string) (*procbuilder.Machine, error) {

	config := new(BondgoConfig)
	config.Debug = false
	config.Verbose = false
	config.Mpm = false
	config.Cascading_io = false

	switch rsize {
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
		config.Rsize = uint8(8)
		config.Basic_type = "uint8"
		config.Basic_chantype = "chan uint8"
	}

	usagedone := make(chan bool)

	results := new(BondgoResults) // Results go in here
	results.Init_Results(config)

	messages := new(BondgoMessages) // Compiler logs and errors
	messages.Init_Messages(config)

	reqmnts := new(BondgoRequirements) // The pointer to the requirements struct
	reqmnts.Init_Requirements(config)

	usagenotify := make(chan UsageNotify) // Used to notify the used resource

	go reqmnts.Usage_Monitor(usagenotify, usagedone) // Spawn the usage monitor

	run := new(BondgoRuninfo) // Running data
	run.Init_Runinfo(config)

	varreq := make(chan VarReq) // Variable request
	varans := make(chan VarAns) // Variable response

	functs := new(BondgoFunctions) // Functions
	functs.Init_Functions(config, messages)

	vars := make(map[string]VarCell)
	returns := make([]VarCell, 0)

	bgmain := &BondgoCheck{results, config, reqmnts, run, messages, functs, usagenotify, varreq, varans, nil, nil, vars, returns, "", "", "device_0", 0}

	bgmain.Used <- UsageNotify{TR_PROC, 0, C_DEVICE, bgmain.CurrentDevice, I_NIL}

	bgmain.Abstract_assembler(rsize, asmcode, usagenotify)

	for _, currline := range asmcode {
		bgmain.WriteLine(0, currline)
	}

	for procid, rout := range bgmain.Program {
		// TODO Recheck
		linesn := len(rout.Lines)
		bgmain.Used <- UsageNotify{TR_PROC, procid, C_ROMSIZE, S_NIL, linesn}
	}

	bgmain.Used <- UsageNotify{TR_EXIT, 0, 0, S_NIL, I_NIL}
	<-usagedone

	if mymachine, ok := bgmain.Create_Connecting_Processor(rsize, 0); ok {
		return mymachine, nil
	} else {
		return nil, errors.New("Creating processor failed")
	}

	return nil, nil
}

func MultiAsm2BondMachine(rsize int, aafile *Abs_assembly) (*bondmachine.Bondmachine, error) {
	bmach := new(bondmachine.Bondmachine)
	bmach.Rsize = uint8(rsize)
	bmach.Init()

	for i, proc_prog := range aafile.ProcProgs {
		if mach, err := Assembly_2_Processor(rsize, strings.Split(proc_prog, "\n")); err == nil {
			bmach.Domains = append(bmach.Domains, mach)
			if _, ok := bmach.Add_processor(i); ok != nil {
				return nil, errors.New("Attach processor failed")
			}

		} else {
			return nil, errors.New("Creating processor failed")
		}
	}

	for _, bond := range aafile.Bonds {
		endpoints := strings.Split(bond, ",")
		if endpoints[0][0] == 'i' {
			for {
				max := "i" + strconv.Itoa(bmach.Inputs)
				if max <= endpoints[0] {
					bmach.Add_input()
				} else {
					break
				}
			}
		}
		if endpoints[0][0] == 'o' {
			for {
				max := "o" + strconv.Itoa(bmach.Outputs)
				if max <= endpoints[0] {
					bmach.Add_output()
				} else {
					break
				}
			}
		}

		if endpoints[1][0] == 'i' {
			for {
				max := "i" + strconv.Itoa(bmach.Inputs)
				if max <= endpoints[1] {
					bmach.Add_input()
				} else {
					break
				}
			}
		}
		if endpoints[1][0] == 'o' {
			for {
				max := "o" + strconv.Itoa(bmach.Outputs)
				if max <= endpoints[1] {
					bmach.Add_output()
				} else {
					break
				}
			}
		}

		bmach.Add_bond(endpoints)

	}

	return bmach, nil
}
