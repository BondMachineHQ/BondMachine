package procbuilder

import (
	"errors"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Saj opcode is both a basic instruction and a template for other instructions.
type Saj struct{}

func (op Saj) Op_get_name() string {
	// TODO
	return "saj"
}

func (op Saj) Op_get_desc() string {
	// TODO
	return "No operation"
}

func (op Saj) Op_show_assembler(arch *Arch) string {
	// TODO
	opbits := arch.Opcodes_bits()
	result := "saj [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op Saj) Op_get_instruction_len(arch *Arch) int {
	// TODO
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op Saj) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (op Saj) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	// TODO
	result := ""
	result += "				SAJ: begin\n"
	result += "					$display(\"SAJ\");\n"
	result += "					_pc <= _pc + 1'b1 ;\n"
	result += "				end\n"

	return result
}

func (op Saj) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Saj) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Saj) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op Saj) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Saj) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Saj) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Saj) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Saj) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Saj) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Saj) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Saj) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Saj) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Saj) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Saj) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Saj) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Saj) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Saj) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Saj) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
