package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The And opcode is both a basic instruction and a template for other instructions.
type Xor struct{}

func (op Xor) Op_get_name() string {
	return "xor"
}

func (op Xor) Op_get_desc() string {
	return "Register bitwise xor"
}

func (op Xor) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "xor [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the logical and of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Xor) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Xor) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (op Xor) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					XOR: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {
			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(j)) + " ^ _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								$display(\"XOR " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"

		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"
	return result
}

func (op Xor) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Xor) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Xor) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The simulation does nothing
func (op Xor) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regdest := get_id(instr[:reg_bits])
	regsrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 8:
		vm.Registers[regdest] = vm.Registers[regdest].(uint8) ^ vm.Registers[regsrc].(uint8)
	case 16:
		vm.Registers[regdest] = vm.Registers[regdest].(uint16) ^ vm.Registers[regsrc].(uint16)
	default:
		// TODO Fix
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Xor) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Xor) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Xor) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Xor) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Xor) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Xor) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Xor) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Xor) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Xor) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Xor) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_REGISTER {

		result := make([]UsageNotify, 3)
		newnot0 := UsageNotify{C_OPCODE, "xor", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1
		newnot2 := UsageNotify{C_REGSIZE, S_NIL, len(seq1)}
		result[2] = newnot2

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op Xor) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Xor) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Xor) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Xor) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
