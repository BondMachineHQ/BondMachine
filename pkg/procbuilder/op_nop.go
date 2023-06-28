package procbuilder

import (
	"errors"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Nop opcode is both a basic instruction and a template for other instructions.
type Nop struct{}

func (op Nop) Op_get_name() string {
	return "nop"
}

func (op Nop) Op_get_desc() string {
	return "No operation"
}

func (op Nop) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "nop [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op Nop) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits
}

func (op Nop) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (Op Nop) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Nop) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Nop) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Nop) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	result := ""
	result += "					NOP: begin\n"
	result += "						$display(\"NOP\");\n"
	result += "						_pc <= _pc + 1'b1 ;\n"
	result += "					end\n"

	return result
}

func (op Nop) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Nop) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Nop) Disassembler(arch *Arch, instr string) (string, error) {
	return "", nil
}

// The simulation does nothing
func (op Nop) Simulate(vm *VM, instr string) error {
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Nop) Generate(arch *Arch) string {
	return ""
}

func (op Nop) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Nop) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Nop) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Nop) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Nop) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "nop", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op Nop) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Nop) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 2)
	result[0] = "nop"
	result[1] = "noop"
	return result
}
func (Op Nop) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "nop":
		return line, nil
	case "noop":
		line.Operation.SetValue("nop")
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Nop) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Nop) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}

