package basm

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

// Assembler2BondMachine transform an assembled instance into a BCOF file
func (bi *BasmInstance) Assembler2BondMachine() error {
	if bi.bm == nil {
		if bi.debug {
			fmt.Println(purple("BondMachine generator") + ": " + red("An existing BondMachine has not been provided, the assembler will create a new one"))
		}
		return bi.assembler2NewBondMachine()
	} else {
		if bi.debug {
			fmt.Println(purple("BondMachine generator") + ": " + red("An existing BondMachine has been provided, the assembler will use that as assembler target"))
		}
		return bi.assembler2ExistingBondMachine()
	}
	return nil
}

func (bi *BasmInstance) assembler2NewBondMachine() error {
	if bi.debug {
		fmt.Println("\t" + green("BondMachine metadata"))
	}

	registerSize := bi.global.GetMeta("registersize")

	if registerSize == "" {
		return errors.New("register size not specified")
	}

	var rSize uint8
	if size, err := strconv.Atoi(registerSize); err == nil {
		if 0 < size && size < 256 {
			rSize = uint8(size)
		} else {
			return errors.New("wrong value for register size")
		}
	} else {
		return errors.New("register size not valid")
	}
	if bi.debug {
		fmt.Println("\t\t"+green("register size:"), rSize)
	}

	if bi.debug {
		fmt.Println("\t" + green("Processors creation"))
	}

	cps := make([]*procbuilder.Machine, len(bi.cps))
	cpIndexes := make(map[string]string)

	bi.rg.Requirement(bmreqs.ReqRequest{Node: "/", T: bmreqs.ObjectSet, Name: "bm", Value: "cps", Op: bmreqs.OpAdd})

	for i, cp := range bi.cps {
		if bi.debug {
			fmt.Print("\t\t" + green("CP: ") + yellow(cp.GetValue()))
		}
		romCode := cp.GetMeta("romcode")
		if romCode == "" {
			return errors.New("CP rom code not found")
		}
		if bi.debug {
			fmt.Println(" - " + green("rom code: ") + yellow(romCode))
		}

		bi.rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps", T: bmreqs.ObjectSet, Name: "id", Value: strconv.Itoa(i), Op: bmreqs.OpAdd})

		if cpm, err := bi.CreateConnectingProcessor(rSize, i, romCode); err == nil {
			cps[i] = cpm
			cpIndexes[cp.GetValue()] = strconv.Itoa(i)
			if bi.BMinfo != nil {
				if bi.CPNames == nil {
					bi.CPNames = make(map[int]string)
				}
				bi.CPNames[i] = cp.GetValue()
			}
			if bi.debug {
				fmt.Println(green("\t\tProcessor created dump: "))
				fmt.Println(green("\t\t----"))
				fmt.Println(cpm)
				fmt.Println(green("\t\t----"))
			}
		} else {
			return err
		}

	}

	bMach := new(bondmachine.Bondmachine)

	bMach.Rsize = uint8(rSize)

	bMach.Init()

	// Attach the Connecting processors
	for i, cp := range cps {
		bMach.Domains = append(bMach.Domains, cp)
		if _, ok := bMach.Add_processor(i); ok != nil {
			return errors.New("attach processor failed")
		}
	}

	// Insert the Shared Objects into the BM and compose and hash of name,constraint
	constrains := make(map[string]string)
	soIndexes := make(map[string]string)
	for i, so := range bi.sos {
		if bi.debug {
			fmt.Print("\t\t" + green("SO: ") + yellow(so.GetValue()))
		}
		constraint := so.GetMeta("constraint")
		if bi.debug {
			fmt.Println(" - " + green("constraint: ") + yellow(constraint))
		}
		bMach.Add_shared_objects([]string{constraint})
		constrains[so.GetValue()] = constraint
		soIndexes[so.GetValue()] = strconv.Itoa(i)
	}

	// Will keep track of the processed attach
	soattDone := make([]bool, len(bi.soAttach))

	// Process every CP
	for i, cp := range bi.cps {
		cpName := cp.GetValue()
		expectedIndex := 0
		expectedIndexS := "0"
		cpConstraints := make([]string, 0)
		for {
			indexFound := false
			// Process every SO attach searching for the couple CP index
			for j, soatt := range bi.soAttach {
				if soatt.GetMeta("cp") == cpName && soatt.GetMeta("index") == expectedIndexS {
					endpoints := make([]string, 2)
					endpoints[0] = strconv.Itoa(i)
					endpoints[1] = soIndexes[soatt.GetValue()]
					// Attach the SO to the CP
					bMach.Connect_processor_shared_object(endpoints)
					indexFound = true
					soattDone[j] = true
					cpConstraints = append(cpConstraints, constrains[soatt.GetValue()])
					break
				}
			}
			if indexFound {
				expectedIndex += 1
				expectedIndexS = strconv.Itoa(expectedIndex)
			} else {
				break
			}
		}
		// compose the Shared constraint of every processor
		bMach.Domains[i].Shared_constraints = strings.Join(cpConstraints, ",")
	}

	for _, val := range soattDone {
		if !val {
			return errors.New("processor SO index inconsistent")
		}
	}

	// Now that Shared Objects are attached, we can assemble the code (it was not possible before)

	for i, cp := range bi.cps {
		romCode := cp.GetMeta("romcode")
		myMachine := bMach.Domains[i]
		myArch := &myMachine.Arch

		prog := ""

		for _, line := range bi.sections[romCode].sectionBody.Lines {
			prog += line.Operation.GetValue()
			for _, arg := range line.Elements {
				prog += " " + arg.GetValue()
			}
			prog += "\n"
		}

		if prog, err := myArch.Assembler([]byte(prog)); err == nil {
			myMachine.Program = prog
		} else {
			return err
		}

	}

	// Will keep track of the processed attach
	ioAttDone := make([]bool, len(bi.ioAttach))

	endPoints := make([]string, len(bi.ioAttach))
	var e1 string
	var e2 string

	for i, ioAtt := range bi.ioAttach {
		if !ioAttDone[i] {
			// Process every IO attach searching for the couple CP index
			for j, ioAtt2 := range bi.ioAttach {
				if !ioAttDone[j] {
					if ioAtt.GetValue() == ioAtt2.GetValue() {
						e1cp := ioAtt.GetMeta("cp")
						e1type := ioAtt.GetMeta("type")
						e1index := ioAtt.GetMeta("index")
						if e1cp == "bm" {
							if e1type == "input" {
								e1 = "i" + e1index
							} else if e1type == "output" {
								e1 = "o" + e1index
							} else {
								return errors.New("wrong IO type")
							}
						} else {
							if e1type == "input" {
								e1 = "p" + cpIndexes[e1cp] + "i" + e1index
							} else if e1type == "output" {
								e1 = "p" + cpIndexes[e1cp] + "o" + e1index
							} else {
								return errors.New("wrong IO type")
							}
						}

						e2cp := ioAtt2.GetMeta("cp")
						e2type := ioAtt2.GetMeta("type")
						e2index := ioAtt2.GetMeta("index")
						if e2cp == "bm" {
							if e2type == "input" {
								e2 = "i" + e2index
							} else if e2type == "output" {
								e2 = "o" + e2index
							} else {
								return errors.New("wrong IO type")
							}
						} else {
							if e2type == "input" {
								e2 = "p" + cpIndexes[e2cp] + "i" + e2index
							} else if e2type == "output" {
								e2 = "p" + cpIndexes[e2cp] + "o" + e2index
							} else {
								return errors.New("wrong IO type")
							}
						}
						endPoints[i] = e1 + "," + e2
						ioAttDone[i] = true
						ioAttDone[j] = true
					}
				}
			}
		}
	}

	// ioOuts := make(map[string]string)
	// ioIns := make(map[string]string)

	// fmt.Println(bi.ioAttach)

	// Process every CP
	// for _, cp := range bi.cps {

	// 	cpName := cp.GetValue()
	// 	// Processing CP inputs
	// 	expectedIndex := 0
	// 	expectedIndexS := "0"
	// 	for {
	// 		indexFound := false
	// 		// Process every IO attach searching for the couple CP index
	// 		for j, ioatt := range bi.ioAttach {
	// 			if ioatt.GetMeta("cp") == cpName && ioatt.GetMeta("index") == expectedIndexS && ioatt.GetMeta("type") == "input" {
	// 				ioname := ioatt.GetValue()
	// 				if curr, ok := ioIns[ioname]; ok {
	// 					ioIns[ioname] = curr + ",p" + cpIndexes[cpName] + "i" + ioatt.GetMeta("index")
	// 				} else {
	// 					ioIns[ioname] = "p" + cpIndexes[cpName] + "i" + ioatt.GetMeta("index")
	// 				}

	// 				indexFound = true
	// 				ioattDone[j] = true
	// 				break
	// 			}
	// 		}
	// 		if indexFound {
	// 			expectedIndex += 1
	// 			expectedIndexS = strconv.Itoa(expectedIndex)
	// 		} else {

	// 			break
	// 		}
	// 	}
	// 	// Processing CP outputs
	// 	expectedIndex = 0
	// 	expectedIndexS = "0"
	// 	for {
	// 		indexFound := false
	// 		// Process every IO attach searching for the couple CP index
	// 		for j, ioatt := range bi.ioAttach {
	// 			if ioatt.GetMeta("cp") == cpName && ioatt.GetMeta("index") == expectedIndexS && ioatt.GetMeta("type") == "output" {
	// 				ioname := ioatt.GetValue()
	// 				fmt.Println(ioname)
	// 				if _, ok := ioOuts[ioname]; ok {
	// 					return errors.New("Multiple IO inconsistency")
	// 				} else {
	// 					ioOuts[ioname] = "p" + cpIndexes[cpName] + "o" + ioatt.GetMeta("index")
	// 				}

	// 				indexFound = true
	// 				ioattDone[j] = true
	// 				break
	// 			}
	// 		}
	// 		if indexFound {
	// 			expectedIndex += 1
	// 			expectedIndexS = strconv.Itoa(expectedIndex)
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }

	// // Processing BM inputs
	// expectedIndex := 0
	// expectedIndexS := "0"
	// for {
	// 	indexFound := false
	// 	// Process every IO attach searching for the couple CP index
	// 	for j, ioatt := range bi.ioAttach {
	// 		if ioatt.GetMeta("cp") == "bm" && ioatt.GetMeta("index") == expectedIndexS && ioatt.GetMeta("type") == "input" {
	// 			ioname := ioatt.GetValue()
	// 			if _, ok := ioOuts[ioname]; ok {
	// 				return errors.New("Multiple IO inconsistency")
	// 			} else {
	// 				ioOuts[ioname] = "i" + ioatt.GetMeta("index")
	// 			}
	// 			bMach.Add_input()
	// 			indexFound = true
	// 			ioattDone[j] = true
	// 			break
	// 		}
	// 	}
	// 	if indexFound {
	// 		expectedIndex += 1
	// 		expectedIndexS = strconv.Itoa(expectedIndex)
	// 	} else {
	// 		break
	// 	}
	// }

	// // Processing BM outputs
	// expectedIndex = 0
	// expectedIndexS = "0"
	// for {
	// 	indexFound := false
	// 	// Process every IO attach searching for the couple CP index
	// 	for j, ioatt := range bi.ioAttach {
	// 		if ioatt.GetMeta("cp") == "bm" && ioatt.GetMeta("index") == expectedIndexS && ioatt.GetMeta("type") == "output" {
	// 			ioname := ioatt.GetValue()
	// 			if curr, ok := ioIns[ioname]; ok {
	// 				ioIns[ioname] = curr + ",o" + ioatt.GetMeta("index")
	// 			} else {
	// 				ioIns[ioname] = "o" + ioatt.GetMeta("index")
	// 			}
	// 			bMach.Add_output()
	// 			indexFound = true
	// 			ioattDone[j] = true
	// 			break
	// 		}
	// 	}
	// 	if indexFound {
	// 		expectedIndex += 1
	// 		expectedIndexS = strconv.Itoa(expectedIndex)
	// 	} else {
	// 		break
	// 	}
	// }

	// for i, _ := range ioattDone {
	// 	// if !val {
	// 	name := bi.ioAttach[i].GetValue()
	// 	fmt.Println(name)
	// 	fmt.Println(ioIns[name])
	// 	fmt.Println(ioOuts[name])
	// 	// }
	// }

	for _, val := range ioAttDone {
		if !val {
			return errors.New("processor IO index inconsistent")
		}
	}

	// for ioName, ioOut := range ioOuts {
	// 	for _, ioIn := range strings.Split(ioIns[ioName], ",") {
	// 		bMach.Add_bond([]string{ioIn, ioOut})
	// 		//fmt.Println([]string{ioIn, ioOut})

	// 	}
	// }

	inToAdd := 0
	outToAdd := 0
	for _, end := range endPoints {
		if end != "" {
			ends := strings.Split(end, ",")
			if ends[0][0] == 'i' {
				inNum, _ := strconv.Atoi(ends[0][1:])
				if inNum+1 > inToAdd {
					inToAdd = inNum + 1
				}
			}
			if ends[0][0] == 'o' {
				outNum, _ := strconv.Atoi(ends[0][1:])
				if outNum+1 > outToAdd {
					outToAdd = outNum + 1
				}
			}
			if ends[1][0] == 'i' {
				inNum, _ := strconv.Atoi(ends[1][1:])
				if inNum+1 > inToAdd {
					inToAdd = inNum + 1
				}
			}
			if ends[1][0] == 'o' {
				outNum, _ := strconv.Atoi(ends[1][1:])
				if outNum+1 > outToAdd {
					outToAdd = outNum + 1
				}
			}
		}
	}

	// TODO recheck the code and include errors handling

	for i := 0; i < inToAdd; i++ {
		bMach.Add_input()
	}

	for i := 0; i < outToAdd; i++ {
		bMach.Add_output()
	}

	for _, end := range endPoints {
		if end != "" {
			//fmt.Println(end)
			bMach.Add_bond(strings.Split(end, ","))
		}
	}
	// fmt.Println(ioIns, ioOuts)

	bi.result = bMach

	return nil
}

func (bi *BasmInstance) assembler2ExistingBondMachine() error {
	// TODO
	return nil
}

func (bi *BasmInstance) CreateConnectingProcessor(rSize uint8, procid int, romCode string) (*procbuilder.Machine, error) {
	myMachine := new(procbuilder.Machine)

	myArch := &myMachine.Arch

	myArch.Rsize = uint8(rSize)

	myArch.Modes = make([]string, 1)
	myArch.Modes[0] = "ha"

	resp := bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "opcodes", Op: bmreqs.OpGet})
	if resp.Error != nil {
		return nil, resp.Error
	}

	bi.rg.Clone("/code:romtexts/sections:"+romCode, "/bm:cps/id:"+strconv.Itoa(procid))

	opCodesS := strings.Split(resp.Value, ",")

	opcodes := make([]procbuilder.Opcode, 0)

	for _, op := range procbuilder.Allopcodes {
		for _, opn := range opCodesS {
			if opn == op.Op_get_name() {
				opcodes = append(opcodes, op)
				break
			}
		}
	}

	sort.Sort(procbuilder.ByName(opcodes))

	myArch.Op = opcodes

	// Getting the registers requirements on the ROM code
	resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "registers", Op: bmreqs.OpGet})
	if resp.Error != nil {
		return nil, resp.Error
	}

	// TODO CHECK: Only the number is relevant for now
	regS := len(strings.Split(resp.Value, ","))
	myArch.R = uint8(Needed_bits(regS))

	// TODO RAM
	// myarch.L = uint8(Needed_bits(preq.Ramsize))
	myArch.L = uint8(8)

	// Getting the Input requirements on the ROM code
	resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "inputs", Op: bmreqs.OpGet})
	if resp.Error != nil {
		myArch.N = uint8(0)
		bi.Warning("No inputs found on ROM code, assuming 0")
	} else {
		// TODO CHECK: Only the number is relevant for now
		inputS := len(strings.Split(resp.Value, ","))
		myArch.N = uint8(inputS)
	}

	// Getting the Output requirements on the ROM code
	resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "outputs", Op: bmreqs.OpGet})
	if resp.Error != nil {
		myArch.M = uint8(0)
		bi.Warning("No outputs found on ROM code, assuming 0")
	} else {
		// TODO CHECK: Only the number is relevant for now
		outputS := len(strings.Split(resp.Value, ","))
		myArch.M = uint8(outputS)
	}

	myArch.O = uint8(Needed_bits(len(bi.sections[romCode].sectionBody.Lines)))

	// The shared constrains will be populated later from the basm metadata
	myArch.Shared_constraints = ""

	// This comports that program will be assembled later, after the bondmachine is created. Eventually, this will be done
	// prog := ""

	// for _, line := range bi.sections[romCode].sectionBody.Lines {
	// 	prog += line.Operation.GetValue()
	// 	for _, arg := range line.Elements {
	// 		prog += " " + arg.GetValue()
	// 	}
	// 	prog += "\n"
	// }

	// if prog, err := myArch.Assembler([]byte(prog)); err == nil {
	// 	myMachine.Program = prog
	// } else {
	// 	return nil, err
	// }

	return myMachine, nil
}

func (bi *BasmInstance) GetBondMachine() *bondmachine.Bondmachine {
	return bi.result
}
