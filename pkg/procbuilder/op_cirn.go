package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Cirn struct{}

func (op Cirn) Op_get_name() string {
	return "cirn"
}

func (op Cirn) Op_get_desc() string {
	return "Register right shift"
}

func (op Cirn) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "cirn [" + strconv.Itoa(int(arch.R)) + "(Reg)] 	// Set a register to the logical and of its value with another register [" + strconv.Itoa(opbits+int(arch.R)) + "]\n"
	return result
}

func (op Cirn) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Cirn) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (op Cirn) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					CIRN: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + ">>> 1'b1;\n"
		result += "								$display(\"CIRN " + strings.ToUpper(Get_register_name(i)) + "\");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Cirn) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Cirn) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	reg_num := 1 << arch.R

	if len(words) != 1 {
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

	for i := opbits + int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Cirn) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	return result, nil
}

// The simulation does nothing
func (op Cirn) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regdest := get_id(instr[:reg_bits])
	regsrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 8:
		vm.Registers[regdest] = vm.Registers[regsrc].(uint8) >> 1
	case 16:
		vm.Registers[regdest] = vm.Registers[regsrc].(uint16) >> 1
	default:
		// TODO Fix
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Cirn) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Cirn) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Cirn) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Cirn) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Cirn) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cirn) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Cirn) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cirn) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cirn) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Cirn) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq, types := Sequence_to_0(words[0])

	if len(seq) > 0 && types == O_REGISTER {

		result := make([]UsageNotify, 2)
		newnot0 := UsageNotify{C_OPCODE, "cirn", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq)}
		result[1] = newnot1

		return result, nil
	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op Cirn) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Cirn) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Cirn) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Cirn) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Cirn) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
