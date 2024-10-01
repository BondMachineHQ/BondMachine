package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Cpy struct{}

func (op Cpy) Op_get_name() string {
	return "cpy"
}

func (op Cpy) Op_get_desc() string {
	return "Copy from a register to another"
}

func (op Cpy) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "chc [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Copy the value of a register to another [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Cpy) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Cpy) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (op Cpy) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	tabsNum := 5
	regNum := 1 << arch.R

	result := ""
	result += tabs(tabsNum) + "CPY: begin\n"
	if th := ThreadInstructionStart(conf, arch, tabsNum+1); th != "" {
		result += th
		tabsNum++
	}
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < regNum; j++ {
			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								$display(\"CPY " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"

		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += tabs(tabsNum+1) + "endcase\n"
	result += NextInstruction(conf, arch, tabsNum+1, "_pc + 1'b1")
	if th := ThreadInstructionEnd(conf, arch, tabsNum); th != "" {
		result += th
		tabsNum--
	}
	result += tabs(tabsNum) + "end\n"
	return result
}

func (op Cpy) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Cpy) Assembler(arch *Arch, words []string) (string, error) {
	// "reference": {"support_asm": "ok"}
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	regNum := 1 << arch.R

	if len(words) != 2 {
		return "", errors.New("wrong arguments number")
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", errors.New("unknown register name " + words[0])
	}

	partial := ""
	for i := 0; i < regNum; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", errors.New("unknown register name " + words[1])
	}

	result += partial

	for i := opBits + 2*int(arch.R); i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Cpy) Disassembler(arch *Arch, instr string) (string, error) {
	regId := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(regId)) + " "
	regId = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(regId))
	return result, nil
}

// The simulation does nothing
func (op Cpy) Simulate(vm *VM, instr string) error {
	regBits := vm.Mach.R
	regDest := get_id(instr[:regBits])
	regSrc := get_id(instr[regBits : regBits*2])
	vm.Registers[regDest] = vm.Registers[regSrc]
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Cpy) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Cpy) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Cpy) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Cpy) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Cpy) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cpy) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Cpy) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cpy) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cpy) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Cpy) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	return nil, errors.New("obsolete")
}

func (Op Cpy) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockName string, objects []string) string {
	result := ""
	switch blockName {
	default:
		result = ""
	}
	return result
}
func (Op Cpy) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "cpy::*--type=reg::*--type=reg")
	result = append(result, "mov::*--type=reg::*--type=reg")
	return result
}
func (Op Cpy) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "cpy":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		if regDst != "" && regSrc != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("cpy")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regDst)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(regSrc)
			newArg1.BasmMeta = newArg1.SetMeta("type", "reg")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Cpy) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Cpy) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	switch line.Operation.GetValue() {
	case "cpy":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		if regDst != "" && regSrc != "" {
			var meta *bmmeta.BasmMeta
			meta = meta.SetMeta("use", regSrc)
			meta = meta.SetMeta("inv", regDst)
			return meta, nil
		}
	case "mov":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		if regDst != "" && regSrc != "" {
			var meta *bmmeta.BasmMeta
			meta = meta.SetMeta("use", regSrc)
			meta = meta.SetMeta("inv", regDst)
			return meta, nil
		}
	}

	return nil, nil
}
