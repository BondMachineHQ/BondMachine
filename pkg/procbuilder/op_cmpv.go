package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Cmpv struct {
}

func (op Cmpv) Op_get_name() string {
	return "cmpv"
}

func (op Cmpv) Op_get_desc() string {
	// "reference": {"desc":"The CMPV instruction check if the input given as argument has a valid signal and sets the cmpflag accordingly."}
	// "reference": {"desc1": "If the input has valid data, cmpflag is set to 1, otherwise it is set to 0."}
	// "reference": {"desc2": "The flag cmpflag can be used by other instructions to take decisions based on the result of this comparison."}
	return "Check if the given input has valid data and set cmpflag accordingly"
}

func (op Cmpv) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	result := "cmpv [" + strconv.Itoa(inBits) + "(Input)]	// Check if the given input has valid data and set cmpflag accordingly [" + strconv.Itoa(opBits+inBits) + "]\n"
	return result
}

func (op Cmpv) Op_get_instruction_len(arch *Arch) int {
	// "reference": {"length": "opBits + inBits"}
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	return opBits + inBits // The bits for the opcode + bits for an input
}

func (op Cmpv) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""
	if arch.OnlyOne(op.Op_get_name(), unique["cmpflag"]) {
		result += "\treg cmpflag;\n"
	}

	return result
}

func (op Cmpv) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()

	tabsNum := 5

	result := ""
	result += tabs(tabsNum) + "CMPV: begin\n"
	if arch.N > 0 {
		if inBits == 1 {
			result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
		} else {
			result += tabs(tabsNum+1) + "case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(inBits)) + "])\n"
		}

		for i := 0; i < int(arch.N); i++ {
			result += tabs(tabsNum+1) + strings.ToUpper(Get_input_name(i)) + " : begin\n"
			result += tabs(tabsNum+2) + "if (" + strings.ToLower(Get_input_name(i)) + "_valid) begin\n"
			result += tabs(tabsNum+3) + "cmpflag <= 1'b1;\n"
			result += tabs(tabsNum+2) + "end else begin\n"
			result += tabs(tabsNum+3) + "cmpflag <= 1'b0;\n"
			result += tabs(tabsNum+2) + "end\n"
			result += tabs(tabsNum+1) + "end\n"
		}
	} else {
		result += tabs(tabsNum+1) + "$display(\"NOP\");\n"
	}

	result += tabs(tabsNum+1) + "endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += tabs(tabsNum) + "end\n"
	return result
}

func (op Cmpv) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpv) Assembler(arch *Arch, words []string) (string, error) {
	// "reference": {"support_asm": "ok"}
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	romWord := arch.Max_word()

	if len(words) != 1 {
		return "", errors.New("wrong arguments number")
	}

	result := ""

	if partial, err := Process_input(words[0], int(arch.N)); err == nil {
		result += zeros_prefix(inBits, partial)
	} else {
		return "", err
	}

	for i := opBits + inBits; i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Cmpv) Disassembler(arch *Arch, instr string) (string, error) {
	// "reference": {"support_disasm": "ok"}
	inBits := arch.Inputs_bits()
	inId := get_id(instr[:inBits])
	result := strings.ToLower(Get_input_name(inId))
	return result, nil
}

func (op Cmpv) Simulate(vm *VM, instr string) error {
	// TODO
	return nil
}

// The random generation does nothing
func (op Cmpv) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Cmpv) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Cmpv) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Cmpv) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Cmpv) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpv) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpv) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpv) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpv) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Cmpv) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (op Cmpv) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Cmpv) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "cmpv::*--type=input")
	result = append(result, "chk::*--type=input")
	return result
}
func (op Cmpv) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "cmpv":
		inNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "inputs", Value: inNeed, Op: bmreqs.OpAdd})
		return line, nil
	case "chk":
		inNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "inputs", Value: inNeed, Op: bmreqs.OpAdd})
		newLine := new(bmline.BasmLine)
		newOp := new(bmline.BasmElement)
		newOp.SetValue("cmpv")
		newLine.Operation = newOp
		newArgs := make([]*bmline.BasmElement, 1)
		newArg1 := new(bmline.BasmElement)
		newArg1.SetValue(inNeed)
		newArg1.BasmMeta = newArg1.SetMeta("type", "input")
		newArgs[0] = newArg1
		newLine.Elements = newArgs
		return newLine, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Cmpv) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Cmpv) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
