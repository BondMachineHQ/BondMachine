package basm

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
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
		bi.result = bi.bm
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
			fmt.Println("\t\t" + green("CP: ") + yellow(cp.GetValue()))
		}
		romCode := cp.GetMeta("romcode")

		if bi.debug {
			if romCode != "" {
				fmt.Println("\t\t - " + green("rom code: ") + yellow(romCode))
			} else {
				fmt.Println("\t\t - " + green("rom code: ") + yellow("not specified"))
			}
		}

		romData := cp.GetMeta("romdata")

		if bi.debug {
			if romData != "" {
				fmt.Println("\t\t - " + green("rom data: ") + yellow(romData))
			} else {
				fmt.Println("\t\t - " + green("rom data: ") + yellow("not specified"))
			}
		}

		ramCode := cp.GetMeta("ramcode")
		if bi.debug {
			if ramCode != "" {
				fmt.Println("\t\t - " + green("ram code: ") + yellow(ramCode))
			} else {
				fmt.Println("\t\t - " + green("ram code: ") + yellow("not specified"))
			}
		}

		ramData := cp.GetMeta("ramdata")
		if bi.debug {
			if ramData != "" {
				fmt.Println("\t\t - " + green("ram data: ") + yellow(ramData))
			} else {
				fmt.Println("\t\t - " + green("ram data: ") + yellow("not specified"))
			}
		}

		execMode := cp.GetMeta("execmode")
		if execMode == "" {
			execMode = bi.global.GetMeta("defaultexecmode")
		}

		if bi.debug {
			if execMode != "" {
				fmt.Println("\t\t - " + green("execution mode: ") + yellow(execMode))
			} else {
				fmt.Println("\t\t - " + green("execution mode: ") + yellow("not specified, defaulting to 'ha'"))
			}
		}

		if execMode == "" {
			execMode = "ha"
		}

		bi.rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps", T: bmreqs.ObjectSet, Name: "id", Value: strconv.Itoa(i), Op: bmreqs.OpAdd})

		if cpm, err := bi.CreateConnectingProcessor(rSize, cp, i, romCode, romData, ramCode, ramData, execMode); err == nil {
			cps[i] = cpm
			cpIndexes[cp.GetValue()] = strconv.Itoa(i)
			if bi.BMinfo != nil {
				if bi.CPNames == nil {
					bi.CPNames = make(map[int]string)
				}
				bi.CPNames[i] = cp.GetValue()
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
		romData := cp.GetMeta("romdata")
		myMachine := bMach.Domains[i]
		myArch := &myMachine.Arch

		prog := ""

		romCodeContrib := 0

		if romCode != "" {
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

			romCodeContrib = len(bi.sections[romCode].sectionBody.Lines)
		}

		if romData != "" {
			wordSize := myArch.Max_word()
			// fmt.Println("Word size: ", wordSize)
			wordPad := ""
			for i := 0; i < int(wordSize); i++ {
				wordPad += "0"
			}
			if wordSize < 8 {
				return errors.New("word size is too small")
			}

			data := make([]string, 0)

			for _, line := range bi.sections[romData].sectionBody.Lines {
				for _, arg := range line.Elements {
					hexVal := arg.GetValue()
					if n, err := bmnumbers.ImportString(hexVal); err == nil {
						nS, _ := n.ExportBinary(false)
						nS = "00000000" + nS
						nS = nS[len(nS)-8:]
						nS = nS + wordPad
						nS = nS[:wordSize]
						data = append(data, nS)
					} else {
						return err
					}
				}
			}

			myArch.O = uint8(Needed_bits(romCodeContrib + len(data)))
			myMachine.Data.Vars = data

		}

		if bi.debug {
			fmt.Println("\t\t - " + green("romsize (post-data): ") + yellow(myArch.O))
			fmt.Println(green("\t\tProcessor created dump: "))
			fmt.Println(green("\t\t----"))
			fmt.Println(cps[i])
			fmt.Println(green("\t\t----"))
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

func (bi *BasmInstance) CreateConnectingProcessor(rSize uint8, cp *bmline.BasmElement, procid int, romCode string, romData string, ramCode string, ramData string, execMode string) (*procbuilder.Machine, error) {
	myMachine := new(procbuilder.Machine)

	myArch := &myMachine.Arch

	myArch.Rsize = uint8(rSize)

	myArch.Modes = make([]string, 1)
	myArch.Modes[0] = execMode

	var resp bmreqs.ReqResponse

	// Processing Code sections: CP has to have at least one code section
	// Getting the ROM code requirements

	if romCode == "" && ramCode == "" {
		return nil, errors.New("no code section specified, neither ROM nor RAM")
	}

	opCodesROM := make([]string, 0)

	if romCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "opcodes", Op: bmreqs.OpGet})
		if resp.Error != nil {
			return nil, resp.Error
		}
		opCodesROM = strings.Split(resp.Value, ",")
	}

	opCodesRAM := make([]string, 0)

	// Getting the RAM code requirements
	if ramCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + ramCode, Name: "opcodes", Op: bmreqs.OpGet})
		if resp.Error != nil {
			return nil, resp.Error
		}
		opCodesRAM = strings.Split(resp.Value, ",")
	}
	// The final list of opCodes is the union of the two lists
	opCodes := make([]procbuilder.Opcode, 0)

outer:
	for _, op := range procbuilder.Allopcodes {
		for _, opn := range opCodesROM {
			if opn == op.Op_get_name() {
				opCodes = append(opCodes, op)
				continue outer
			}
		}
		for _, opn := range opCodesRAM {
			if opn == op.Op_get_name() {
				opCodes = append(opCodes, op)
				continue outer
			}
		}
	}

	sort.Sort(procbuilder.ByName(opCodes))

	myArch.Op = opCodes

	// TODO: check how it is used and if it is needed, eventually remove or substitute with the merge of the two lists
	bi.rg.Clone("/code:romtexts/sections:"+romCode, "/bm:cps/id:"+strconv.Itoa(procid))

	regs := make([]string, 0)

	// Getting the registers requirements on the ROM code
	if romCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "registers", Op: bmreqs.OpGet})
		if resp.Error == nil {
			regs = strings.Split(resp.Value, ",")
		}
	}

	// Getting the registers requirements on the RAM code, appending to the previous list
	if ramCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + ramCode, Name: "registers", Op: bmreqs.OpGet})
		if resp.Error == nil {
			for _, reg := range strings.Split(resp.Value, ",") {
				if !stringInSlice(reg, regs) {
					regs = append(regs, reg)
				}
			}
		}
	}

	if len(regs) == 0 {
		return nil, errors.New("no registers found on ROM/RAM code")
	}

	// Sorting the registers list (ordering using the compareStrings function)
	sort.Slice(regs, func(i, j int) bool {
		return compareStrings(regs[i], regs[j])
	})

	// Getting the last register in the list
	lastReg := regs[len(regs)-1]

	// Getting the register number
	regNum, _ := strconv.Atoi(lastReg[1:])

	// To store up to the last register, we need regNum+1 registers
	myArch.R = uint8(Needed_bits(regNum + 1))

	// If the length of the register list is different from the number of registers, emit a warning
	if len(regs) != regNum+1 {
		bi.Warning("Register list is not complete, some registers are missing. This is not an error provided you know what you are doing.")
	}

	ins := make([]string, 0)

	// Getting the Input requirements on the ROM code
	if romCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "inputs", Op: bmreqs.OpGet})
		if resp.Error == nil {
			ins = strings.Split(resp.Value, ",")
		}
	}
	// Getting the Input requirements on the RAM code, appending to the previous list
	if ramCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + ramCode, Name: "inputs", Op: bmreqs.OpGet})
		if resp.Error == nil {
			for _, in := range strings.Split(resp.Value, ",") {
				if !stringInSlice(in, ins) {
					ins = append(ins, in)
				}
			}
		}
	}

	if len(ins) == 0 {
		bi.Warning("No inputs found on ROM/RAM code, assuming 0")
		myArch.N = uint8(0)
	} else {
		// Sorting the inputs list (ordering using the compareStrings function)
		sort.Slice(ins, func(i, j int) bool {
			return compareStrings(ins[i], ins[j])
		})

		// Getting the last input in the list
		lastIn := ins[len(ins)-1]

		// Getting the input number
		inNum, _ := strconv.Atoi(lastIn[1:])
		// To store up to the last input, we need inNum+1 inputs
		myArch.N = uint8(inNum + 1)

		if len(ins) != inNum+1 {
			bi.Warning("Input list is not complete, some inputs are missing. This is not an error, but you are wasting resources.")
		}
	}

	outs := make([]string, 0)

	// Getting the Output requirements on the ROM code
	if romCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts/sections:" + romCode, Name: "outputs", Op: bmreqs.OpGet})
		if resp.Error == nil {
			outs = strings.Split(resp.Value, ",")
		}
	}

	// Getting the Output requirements on the RAM code, appending to the previous list
	if ramCode != "" {
		resp = bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts/sections:" + ramCode, Name: "outputs", Op: bmreqs.OpGet})
		if resp.Error == nil {
			for _, out := range strings.Split(resp.Value, ",") {
				if !stringInSlice(out, outs) {
					outs = append(outs, out)
				}
			}
		}
	}

	if len(outs) == 0 {
		bi.Warning("No outputs found on ROM/RAM code, assuming 0")
		myArch.M = uint8(0)
	} else {
		// Sorting the outputs list (ordering using the compareStrings function)
		sort.Slice(outs, func(i, j int) bool {
			return compareStrings(outs[i], outs[j])
		})

		// Getting the last output in the list
		lastOut := outs[len(outs)-1]

		// Getting the output number
		outNum, _ := strconv.Atoi(lastOut[1:])
		// To store up to the last output, we need outNum+1 outputs
		myArch.M = uint8(outNum + 1)

		if len(outs) != outNum+1 {
			bi.Warning("Output list is not complete, some outputs are missing. This is not an error, but you are wasting resources.")
		}
	}

	// Here start the mess with the RAM/ROM size, word size, etc.
	// Compute the word size and the RAM/ROM size, check the eventual cp options

	// This is the sequence of check in order of priority (valid both for ROM and RAM):
	// - Check if there is a cp option for the word size and use it overriding the other options. In the case of missing resources, the assembler will fail
	// - TODO: Check requirements
	//   - TODO From the direct opcodes
	//   - TODO From the indirect opcodes
	// - Check the ROM code size and use it (eventually adding the data section size).

	romCodeContrib := 0
	if cp.GetMeta("romsize") != "" {
		if val, err := strconv.Atoi(cp.GetMeta("romsize")); err == nil {
			romCodeContrib = 2 ^ val
			myArch.O = uint8(val)
			if bi.debug {
				fmt.Println("\t\t - " + green("romsize (cp config): ") + yellow(cp.GetMeta("romsize")))
			}
		} else {
			return nil, err
		}
	} else if romCode != "" {
		romCodeContrib = len(bi.sections[romCode].sectionBody.Lines)
		myArch.O = uint8(Needed_bits(romCodeContrib))
		if bi.debug {
			fmt.Println("\t\t - " + green("romsize (pre-data): ") + yellow(strconv.Itoa(Needed_bits(len(bi.sections[romCode].sectionBody.Lines)))))
		}
	} else {
		romCodeContrib = 0
		myArch.O = uint8(0)
		if bi.debug {
			fmt.Println("\t\t - " + green("romsize (pre-data): ") + yellow("0"))
		}
	}

	ramCodeContrib := 0

	if cp.GetMeta("ramsize") != "" {
		if val, err := strconv.Atoi(cp.GetMeta("ramsize")); err == nil {
			ramCodeContrib = 2 ^ val
			myArch.L = uint8(val)
			if bi.debug {
				fmt.Println("\t\t - " + green("ramsize (cp config): ") + yellow(cp.GetMeta("ramsize")))
			}
		} else {
			return nil, err
		}

	} else if ramCode != "" {
		ramCodeContrib = len(bi.sections[ramCode].sectionBody.Lines)
		myArch.L = uint8(Needed_bits(ramCodeContrib))
		if bi.debug {
			fmt.Println("\t\t - " + green("ramsize (pre-data): ") + yellow(strconv.Itoa(Needed_bits(len(bi.sections[ramCode].sectionBody.Lines)))))
		}
	} else {
		myArch.L = uint8(0)
		ramCodeContrib = 0
		if bi.debug {
			fmt.Println("\t\t - " + green("ramsize (pre-data): ") + yellow("0"))
		}
	}

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

	// If there is a data section, we need to add it to the machine and update the myArch.O field prior to assembling
	if romData != "" {
		wordSize := myMachine.Max_word()
		// fmt.Println("Word size: ", wordSize)
		wordPad := ""
		for i := 0; i < int(wordSize); i++ {
			wordPad += "0"
		}
		if wordSize < 8 {
			return nil, errors.New("word size is too small")
		}

		data := make([]string, 0)

		for _, line := range bi.sections[romData].sectionBody.Lines {
			for _, arg := range line.Elements {
				hexVal := arg.GetValue()
				if n, err := bmnumbers.ImportString(hexVal); err == nil {
					nS, _ := n.ExportBinary(false)
					nS = nS + wordPad
					nS = nS[len(nS)-int(wordSize):]
					data = append(data, nS)
				} else {
					return nil, err
				}
			}
		}

		myArch.O = uint8(Needed_bits(romCodeContrib + len(data)))
	}

	if ramData != "" {
		rSize := myArch.Rsize
		ramPad := ""
		for i := 0; i < int(rSize); i++ {
			ramPad += "0"
		}
		if rSize < 8 {
			return nil, errors.New("register size is too small")
		}

		data := make([]string, 0)

		for _, line := range bi.sections[ramData].sectionBody.Lines {
			for _, arg := range line.Elements {
				hexVal := arg.GetValue()
				if n, err := bmnumbers.ImportString(hexVal); err == nil {
					nS, _ := n.ExportBinary(false)
					nS = nS + ramPad
					nS = nS[len(nS)-int(rSize):]
					data = append(data, nS)
				} else {
					return nil, err
				}
			}
		}

		myArch.L = uint8(Needed_bits(ramCodeContrib + len(data)))
	}

	return myMachine, nil
}

func (bi *BasmInstance) GetBondMachine() *bondmachine.Bondmachine {
	return bi.result
}
