package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

const (
	SICV2SM_IDLE = uint8(0) + iota
	SICV2SM_WAIT
	SICV2SM_COUNT
	SICV2SM_DONE
)

type Sicv2 struct{}

func (op Sicv2) Op_get_name() string {
	return "sicv2"
}

func (op Sicv2) Op_get_desc() string {
	return "Wait for an input change via valid and increments a register "
}

func (op Sicv2) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	result := "sicv2 [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(inBits) + "(Input)]	// Wait for an input change via valid and increment the register [" + strconv.Itoa(opBits+int(arch.R)+inBits) + "]\n"
	return result
}

func (op Sicv2) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	return opBits + int(arch.R) + 2*int(inBits) // The bits for the opcode + bits for a register + bits for two inputs
}

func (op Sicv2) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := "\n"
	result += "\t//State Machine for SICv2\n"
	result += "\treg [1:0] sicv2_sm;\n"
	result += "\tlocalparam SICV2SM_IDLE = 2'b00;\n"
	result += "\tlocalparam SICV2SM_WAIT = 2'b01;\n"
	result += "\tlocalparam SICV2SM_COUNT = 2'b10;\n"
	result += "\tlocalparam SICV2SM_DONE = 2'b11;\n"
	result += "\tinitial begin\n"
	result += "\t\tsicv2_sm = SICV2SM_IDLE;\n"
	result += "\tend\n"
	result += "\n"

	if arch.OnlyOne(op.Op_get_name(), unique["inputrecv"]) {

		opbits := arch.Opcodes_bits()
		rom_word := arch.Max_word()

		result += "\n"
		// Data received fro inputs
		for j := 0; j < int(arch.N); j++ {
			result += "\treg " + strings.ToLower(Get_input_name(j)) + "_recv;\n"
		}

		result += "\n"

		for j := 0; j < int(arch.N); j++ {

			objects := make([]string, 1)
			objects[0] = strings.ToLower(Get_input_name(j))

			// Process for data outputs data valid
			result += "\talways @(posedge clock_signal, posedge reset_signal)\n"
			result += "\tbegin\n"

			result += "\t\tif (reset_signal)\n"
			result += "\t\tbegin\n"
			result += "\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
			result += "\t\tend\n"
			result += "\t\telse\n"
			result += "\t\tbegin\n"

			if opbits == 1 {
				result += "\t\t\tcase(current_instruction[" + strconv.Itoa(rom_word-1) + "])\n"
			} else {
				result += "\t\t\tcase(current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "])\n"
			}

			for _, currop := range arch.Op {
				result += currop.Op_instruction_verilog_extra_block(arch, flavor, uint8(4), "input_data_received", objects)
			}

			result += "\t\t\t\tdefault: begin\n"
			result += "\t\t\t\t\tif (!" + strings.ToLower(Get_input_name(j)) + "_valid)\n"
			result += "\t\t\t\t\tbegin\n"
			result += "\t\t\t\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
			result += "\t\t\t\t\tend\n"
			result += "\t\t\t\tend\n"

			result += "\t\t\tendcase\n"

			result += "\t\tend\n"

			result += "\tend\n"
		}
	}

	return result
}

func (op Sicv2) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	tabsNum := 5
	regNum := 1 << arch.R
	inBits := arch.Inputs_bits()

	result := ""
	result += tabs(tabsNum) + "SICV2: begin\n"

	if arch.R == 1 {
		result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
	} else {
		result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {
		result += tabs(tabsNum+2) + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if inBits == 1 {
			result += tabs(tabsNum+3) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += tabs(tabsNum+3) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)) + "])\n"
		}

		for j := 0; j < int(arch.N); j++ {
			result += tabs(tabsNum+4) + strings.ToUpper(Get_input_name(j)) + " : begin\n"
			if inBits == 1 {
				result += tabs(tabsNum+5) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)-1) + "])\n"
			} else {
				result += tabs(tabsNum+5) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-2*int(inBits)) + "])\n"
			}
			for k := 0; k < int(arch.N); k++ {
				result += tabs(tabsNum+5) + strings.ToUpper(Get_input_name(k)) + " : begin\n"
				result += tabs(tabsNum+6) + "case (sicv2_sm)\n"
				result += tabs(tabsNum+6) + "SICV2SM_IDLE: begin\n"
				result += tabs(tabsNum+7) + "if (" + strings.ToLower(Get_input_name(j)) + "_valid) begin\n"
				result += tabs(tabsNum+8) + "sicv2_sm <= SICV2SM_COUNT;\n"
				result += tabs(tabsNum+7) + "end\n"
				result += tabs(tabsNum+6) + "end\n"
				result += tabs(tabsNum+6) + "SICV2SM_COUNT: begin\n"
				result += tabs(tabsNum+7) + "if (" + strings.ToLower(Get_input_name(k)) + "_valid) begin\n"
				result += tabs(tabsNum+8) + "sicv2_sm <= SICV2SM_IDLE;\n"
				result += tabs(tabsNum+8) + "_pc <= #1 _pc + 1'b1;\n"
				result += tabs(tabsNum+7) + "end else begin\n"
				result += tabs(tabsNum+8) + "_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + " + 1'b1;\n"
				result += tabs(tabsNum+7) + "end\n"
				result += tabs(tabsNum+6) + "end\n"
				result += tabs(tabsNum+6) + "endcase\n"
				result += tabs(tabsNum+5) + "end\n"
			}
			result += tabs(tabsNum+5) + "endcase\n"

			// result += NextInstruction(conf, arch, 9, "_pc + 1'b1")
			result += tabs(tabsNum+5) + "$display(\"SICV2 " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_input_name(j)) + " - PC=%d\", _pc);\n"
			result += tabs(tabsNum+4) + "end\n"

		}
		result += tabs(tabsNum+3) + "endcase\n"
		result += tabs(tabsNum+2) + "end\n"
	}
	result += tabs(tabsNum+2) + "endcase\n"

	result += tabs(tabsNum+1) + "end\n"
	return result
}

func (op Sicv2) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Sicv2) Assembler(arch *Arch, words []string) (string, error) {
	// "reference": {"support_asm": "ok"}
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	romWord := arch.Max_word()

	regNum := 2
	regNum = regNum << (arch.R - 1)

	if len(words) != 3 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"Unknown register name " + words[0]}
	}

	if partial, err := Process_input(words[1], int(arch.N)); err == nil {
		result += zeros_prefix(inBits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	if partial, err := Process_input(words[2], int(arch.N)); err == nil {
		result += zeros_prefix(inBits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opBits + int(arch.R) + 2*inBits; i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Sicv2) Disassembler(arch *Arch, instr string) (string, error) {
	// "reference": {"support_asm": "ok"}
	inBits := arch.Inputs_bits()
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	inpId1 := get_id(instr[arch.R : int(arch.R)+inBits])
	result += strings.ToLower(Get_input_name(inpId1)) + " "
	inpId2 := get_id(instr[int(arch.R)+inBits : int(arch.R)+2*inBits])
	result += strings.ToLower(Get_input_name(inpId2)) + " "

	return result, nil
}

func (op Sicv2) Simulate(vm *VM, instr string) error {
	// "reference": {"support_sim": "testing"}
	inBits := vm.Mach.Inputs_bits()
	regBits := vm.Mach.R
	reg := get_id(instr[:regBits])
	in0 := get_id(instr[regBits : int(regBits)+inBits])
	in1 := get_id(instr[int(regBits)+inBits : int(regBits)+2*inBits])

	i0Valid := vm.InputsValid[in0]
	i1Valid := vm.InputsValid[in1]

	var sicv2State uint8
	if state, ok := vm.Extra_states["sicv2_state"]; ok {
		sicv2State = state.(uint8)
	} else {
		vm.Extra_states["sicv2_state"] = SICV2SM_IDLE
		sicv2State = SICV2SM_IDLE
	}

	// switch sicv2State {
	// case SICV2SM_IDLE:
	// 	fmt.Println("SICV2 STATE: IDLE")
	// case SICV2SM_WAIT:
	// 	fmt.Println("SICV2 STATE: WAIT")
	// case SICV2SM_COUNT:
	// 	fmt.Println("SICV2 STATE: COUNT")
	// case SICV2SM_DONE:
	// 	fmt.Println("SICV2 STATE: DONE")
	// }

	switch sicv2State {
	case SICV2SM_IDLE:
		if !i0Valid {
			vm.InputsRecv[in0] = false
		} else {
			vm.Extra_states["sicv2_state"] = SICV2SM_WAIT
			vm.InputsRecv[in0] = true
		}
	case SICV2SM_WAIT:
		if !i0Valid {
			vm.InputsRecv[in0] = false
			vm.Extra_states["sicv2_state"] = SICV2SM_COUNT
			switch vm.Mach.Rsize {
			case 8:
				vm.Registers[reg] = uint8(0)
			case 16:
				vm.Registers[reg] = uint16(0)
			case 32:
				vm.Registers[reg] = uint32(0)
			case 64:
				vm.Registers[reg] = uint64(0)
			default:
				return errors.New("go simulation only works on 8,16,32 or 64 bits registers")
			}
		}
	case SICV2SM_COUNT:
		if !i1Valid {
			vm.InputsRecv[in1] = false
			switch vm.Mach.Rsize {
			case 8:
				vm.Registers[reg] = vm.Registers[reg].(uint8) + 1
			case 16:
				vm.Registers[reg] = vm.Registers[reg].(uint16) + 1
			case 32:
				vm.Registers[reg] = vm.Registers[reg].(uint32) + 1
			case 64:
				vm.Registers[reg] = vm.Registers[reg].(uint64) + 1
			default:
				return errors.New("go simulation only works on 8,16,32 or 64 bits registers")
			}
		} else {
			vm.Extra_states["sicv2_state"] = SICV2SM_DONE
			vm.InputsRecv[in1] = true
		}
	case SICV2SM_DONE:
		if !i1Valid {
			vm.InputsRecv[in1] = false
			vm.Pc = vm.Pc + 1
			vm.Extra_states["sicv2_state"] = SICV2SM_IDLE
		}
	}
	return nil
}

func (op Sicv2) Generate(arch *Arch) string {
	inBits := arch.Inputs_bits()
	regNum := 1 << arch.R
	reg := rand.Intn(regNum)
	inp := rand.Intn(int(arch.N))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(inBits, get_binary(inp))
}

func (op Sicv2) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Sicv2) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Sicv2) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Sicv2) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv2) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv2) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv2) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Sicv2) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Sicv2) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "sicv2", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1

		for i, _ := range seq1 {
			result[i+2] = UsageNotify{C_INPUT, S_NIL, i + 1}
		}

		return result, nil

	}

	return []UsageNotify{}, errors.New("wrong parameters")
}

func (Op Sicv2) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	romWord := arch.Max_word()

	result := ""

	pref := strings.Repeat("\t", int(level))

	switch blockname {
	case "input_data_received":
		result += pref + "SICV2: begin\n"
		in1 := ""
		if inBits == 1 {
			in1 = "current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + "]"
		} else {
			in1 = "current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)) + "]"
		}

		in2 := ""
		if inBits == 1 {
			in2 = "current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)-1) + "]"
		} else {
			in2 = "current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-int(inBits)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-2*int(inBits)) + "]"
		}

		result += pref + "\tif (" + in1 + " == " + strings.ToUpper(objects[0]) + " || " + in2 + " == " + strings.ToUpper(objects[0]) + ") begin\n"

		result += pref + "\t\tif (" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b1;\n"
		result += pref + "\t\tend else begin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\telse begin\n"
		result += pref + "\t\tif (!" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "end\n"
	default:
		result = ""
	}

	return result
}
func (Op Sicv2) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Sicv2) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	// "reference": {"support_hlasm": "notapplicable"}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Sicv2) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Sicv2) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	// "reference": {"support_asmeta": "notapplicable"}
	return nil, nil
}
