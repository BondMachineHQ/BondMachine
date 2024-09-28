package procbuilder

import (
	"errors"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Nop opcode is a no operation instruction
type Nop struct{}

func (op Nop) Op_get_name() string {
	return "nop"
}

func (op Nop) Op_get_desc() string {
	// "reference": {"desc":"The NOP instruction does nothing. It is used to fill the instruction memory with no operation instructions. And waste a clock cycle."}
	return "No operation"
}

func (op Nop) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "nop [" + strconv.Itoa(opBits) + "]	// No operation [" + strconv.Itoa(opBits) + "]\n"
	return result
}

func (op Nop) Op_get_instruction_len(arch *Arch) int {
	// "reference": {"length": "opBits"}
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
	// "reference": {"support_hdl":"ok"}
	// "reference": {"support_mt":"test"}
	tabsNum := 5

	result := ""
	result += tabs(tabsNum) + "NOP: begin\n"
	if th := ThreadInstructionStart(conf, arch, tabsNum+1); th != "" {
		result += th
		tabsNum++
	}
	result += tabs(tabsNum+1) + "$display(\"NOP\");\n"
	result += NextInstruction(conf, arch, tabsNum+1, "_pc + 1'b1")
	if th := ThreadInstructionEnd(conf, arch, tabsNum); th != "" {
		result += th
		tabsNum--
	}
	result += tabs(tabsNum) + "end\n"
	return result
}

func (op Nop) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Nop) Assembler(arch *Arch, words []string) (string, error) {
	// "reference": {"support_asm":"ok"}
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()
	result := ""
	for i := opBits; i < romWord; i++ {
		result += "0"
	}
	return result, nil
}

func (op Nop) Disassembler(arch *Arch, instr string) (string, error) {
	// "reference": {"support_disasm":"ok"}
	return "", nil
}

// The simulation does nothing
func (op Nop) Simulate(vm *VM, instr string) error {
	// "reference": {"support_sim":"ok"}
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
	return []UsageNotify{}, errors.New("obsolete")
}

func (Op Nop) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockName string, objects []string) string {
	result := ""
	switch blockName {
	default:
		result = ""
	}
	return result
}
func (Op Nop) HLAssemblerMatch(arch *Arch) []string {
	// "reference": {"support_hlasm":"ok"}
	result := make([]string, 2)
	result[0] = "nop"
	result[1] = "noop"
	return result
}
func (Op Nop) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "nop":
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "nop", Op: bmreqs.OpAdd})
		return line, nil
	case "noop":
		line.Operation.SetValue("nop")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "nop", Op: bmreqs.OpAdd})
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
