package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
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
	opbits := arch.Opcodes_bits()
	result := "chc [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Copy the value of a register to another [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Cpy) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Cpy) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (op Cpy) Op_instruction_verilog_state_machine(arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opbits := arch.Opcodes_bits()

	regNum := 1 << arch.R

	result := ""
	result += "					CPY: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(romWord-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(romWord-opbits-1) + ":" + strconv.Itoa(romWord-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(romWord-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(romWord-opbits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opbits-int(arch.R)-int(arch.R)) + "])\n"
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
	result += "						endcase\n"
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"
	return result
}

func (op Cpy) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Cpy) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Cpy) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The simulation does nothing
func (op Cpy) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regDest := get_id(instr[:reg_bits])
	regSrc := get_id(instr[reg_bits : reg_bits*2])
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
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Cpy) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
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
