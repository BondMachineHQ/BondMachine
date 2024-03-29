package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Addp struct {
	pipeline *bool
}

func (op Addp) Op_get_name() string {
	return "addp"
}

func (op Addp) Op_get_desc() string {
	return "Register pipelined addition"
}

func (op Addp) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "addp [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Addp) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Addp) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] addp_" + arch.Tag + "_input_a;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] addp_" + arch.Tag + "_input_b;\n"
	result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] addp_" + arch.Tag + "_output_z;\n"

	result += "\treg	[1:0] addp_" + arch.Tag + "_state;\n"
	result += "parameter addp_" + arch.Tag + "_put         = 2'd0,\n"
	result += "          addp_" + arch.Tag + "_get         = 2'd1;\n"

	result += "\taddp_" + arch.Tag + " addp_" + arch.Tag + "_inst (addp_" + arch.Tag + "_input_a, addp_" + arch.Tag + "_input_b,  addp_" + arch.Tag + "_output_z);\n\n"

	return result
}

func (op Addp) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					ADDP: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:addp", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:addp", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (addp_" + arch.Tag + "_state)\n"
			result += "							addp_" + arch.Tag + "_put : begin\n"
			result += "								addp_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								addp_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								addp_" + arch.Tag + "_state <= #1 addp_" + arch.Tag + "_get;\n"
			result += "							end\n"
			result += "							addp_" + arch.Tag + "_get : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 addp_" + arch.Tag + "_output_z;\n"
			result += "								addp_" + arch.Tag + "_state <= #1 addp_" + arch.Tag + "_put;\n"
			result += NextInstruction(conf, arch, 8, "_pc + 1'b1")
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"ADDP " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Addp) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Addp) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
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

	for i := opBits + 2*int(arch.R); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Addp) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Addp) Simulate(vm *VM, instr string) error {
	regBits := vm.Mach.R
	regDest := get_id(instr[:regBits])
	regSrc := get_id(instr[regBits : regBits*2])

	if *op.pipeline {
		switch vm.Mach.Rsize {
		case 8:
			vm.Registers[regDest] = vm.Registers[regDest].(uint8) + vm.Registers[regSrc].(uint8)
		case 16:
			vm.Registers[regDest] = vm.Registers[regDest].(uint16) + vm.Registers[regSrc].(uint16)
		case 32:
			vm.Registers[regDest] = vm.Registers[regDest].(uint32) + vm.Registers[regSrc].(uint32)
		case 64:
			vm.Registers[regDest] = vm.Registers[regDest].(uint64) + vm.Registers[regSrc].(uint64)
		default:
			return errors.New("invalid register size")
		}
		vm.Pc = vm.Pc + 1
		*op.pipeline = false
	} else {
		*op.pipeline = true
	}
	return nil
}

// The random generation does nothing
func (op Addp) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Addp) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Addp) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Addp) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Addp) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addp) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Addp) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addp) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addp) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	result := "\n\n"

	result += "module addp_" + arch.Tag + "(\n"
	result += "        input_a,\n"
	result += "        input_b,\n"
	result += "        output_z);\n"
	result += "\n"
	result += "  input     [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_a;\n"
	result += "  input     [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_b;\n"
	result += "  output    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] output_z;\n"
	result += "  assign output_z = input_a + input_b;\n"
	result += "\n"
	result += "endmodule\n"

	return []string{"multiplier"}, []string{result}
}

func (Op Addp) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Addp) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Addp) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "addp::*--type=reg::*--type=reg")
	return result
}
func (Op Addp) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "addp":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "addp", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:addp", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:addp", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Addp) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Addp) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
