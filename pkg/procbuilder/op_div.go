package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Div opcode is both a basic instruction and a template for other instructions.
type Div struct{}

func (op Div) Op_get_name() string {
	return "div"
}

func (op Div) Op_get_desc() string {
	return "Register div"
}

func (op Div) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "div [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the division of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Div) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Div) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (op Div) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					DIV: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {
			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + " / _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								$display(\"DIV " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"

		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Div) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Div) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	reg_num := 1 << arch.R

	if len(words) != 2 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < reg_num; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"Unknown register name " + words[0]}
	}

	partial := ""
	for i := 0; i < reg_num; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", Prerror{"Unknown register name " + words[1]}
	}

	result += partial

	for i := opbits + 2*int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Div) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The Div opcode simulates the execution of a div instruction
func (op Div) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regdest := get_id(instr[:reg_bits])
	regsrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 8:
		vm.Registers[regdest] = vm.Registers[regdest].(uint8) / vm.Registers[regsrc].(uint8)
	case 16:
		vm.Registers[regdest] = vm.Registers[regdest].(uint16) / vm.Registers[regsrc].(uint16)
	case 32:
		vm.Registers[regdest] = vm.Registers[regdest].(uint32) / vm.Registers[regsrc].(uint32)
	case 64:
		vm.Registers[regdest] = vm.Registers[regdest].(uint64) / vm.Registers[regsrc].(uint64)
	default:
		return errors.New("invalid register size")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Div) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Div) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Div) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Div) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Div) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Div) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Div) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Div) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Div) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Div) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_REGISTER {

		result := make([]UsageNotify, 3)
		newnot0 := UsageNotify{C_OPCODE, "div", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1
		newnot2 := UsageNotify{C_REGSIZE, S_NIL, len(seq1)}
		result[2] = newnot2

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op Div) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Div) HLAssemblerMatch(arch *Arch) []string {

	result := make([]string, 0)
	result = append(result, "div::*--type=reg::*--type=reg")
	return result
}
func (Op Div) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "div":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}

func (Op Div) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Div) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
