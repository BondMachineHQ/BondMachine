package procbuilder

import (
	"errors"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Hlt opcode is both a basic instruction and a template for other instructions.
type Hlt struct{}

func (op Hlt) Op_get_name() string {
	return "hlt"
}

func (op Hlt) Op_get_desc() string {
	return "Halt the processor"
}

func (op Hlt) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "hlt [" + strconv.Itoa(opbits) + "]	// Halt the processor [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op Hlt) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op Hlt) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (Op Hlt) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Hlt) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Hlt) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Hlt) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	result := ""
	result += "				HLT: begin\n"
	result += "					$display(\"HLT\");\n"
	result += "				end\n"

	return result
}

func (op Hlt) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Hlt) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Hlt) Disassembler(arch *Arch, instr string) (string, error) {
	return "", nil
}

// The simulation does nothing
func (op Hlt) Simulate(vm *VM, instr string) error {
	// TODO
	return nil
}

// The random genaration does nothing
func (op Hlt) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Hlt) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Hlt) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Hlt) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Hlt) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Hlt) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Hlt) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Hlt) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Hlt) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Hlt) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
