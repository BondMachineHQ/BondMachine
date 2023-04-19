package procbuilder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The LinearQuantizer opcode is both a basic instruction and a template for other instructions.
type LinearQuantizer struct {
	lqName string
	max    float64
	s      int
	t      int
	opType uint8
}

func (op LinearQuantizer) Op_get_name() string {
	return op.lqName
}

func (op LinearQuantizer) Op_get_desc() string {
	return "LinearQuantizer dynamical instruction " + op.lqName
}

func (op LinearQuantizer) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := op.lqName + " [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op LinearQuantizer) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op LinearQuantizer) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.lqName + "_" + arch.Tag + "_input_a;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.lqName + "_" + arch.Tag + "_input_b;\n"
	result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.lqName + "_" + arch.Tag + "_output_z;\n"

	if op.opType == LQMULT || op.opType == LQDIV {
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.lqName + "_" + arch.Tag + "_input_corr;\n"
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + op.lqName + "_" + arch.Tag + "_output_corr;\n"
	}

	result += "\treg	[1:0] " + op.lqName + "_" + arch.Tag + "_state;\n"
	result += "parameter " + op.lqName + "_" + arch.Tag + "_put         = 2'd0,\n"
	result += "          " + op.lqName + "_" + arch.Tag + "_corr        = 2'd1,\n"
	result += "          " + op.lqName + "_" + arch.Tag + "_get         = 2'd2;\n"

	result += "\t" + op.lqName + "_" + arch.Tag + " " + op.lqName + "_" + arch.Tag + "_inst (" + op.lqName + "_" + arch.Tag + "_input_a, " + op.lqName + "_" + arch.Tag + "_input_b,  " + op.lqName + "_" + arch.Tag + "_output_z);\n\n"

	if op.opType == LQMULT || op.opType == LQDIV {
		result += "\t" + op.lqName + "_correction_" + arch.Tag + " " + op.lqName + "_correction_" + arch.Tag + "_inst (" + op.lqName + "_" + arch.Tag + "_input_corr, " + op.lqName + "_" + arch.Tag + "_output_corr);\n\n"
	}

	return result
}

func (op LinearQuantizer) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					" + strings.ToUpper(op.lqName) + ": begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:" + op.lqName, T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opBits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:" + op.lqName, T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (" + op.lqName + "_" + arch.Tag + "_state)\n"
			result += "							" + op.lqName + "_" + arch.Tag + "_put : begin\n"
			result += "								" + op.lqName + "_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								" + op.lqName + "_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			switch op.opType {
			case LQADD:
				result += "								" + op.lqName + "_" + arch.Tag + "_state <= #1 " + op.lqName + "_" + arch.Tag + "_get;\n"
			case LQMULT, LQDIV:
				result += "								" + op.lqName + "_" + arch.Tag + "_state <= #1 " + op.lqName + "_" + arch.Tag + "_corr;\n"
			}
			result += "							end\n"
			if op.opType == LQMULT || op.opType == LQDIV {
				result += "							" + op.lqName + "_" + arch.Tag + "_corr : begin\n"
				result += "								" + op.lqName + "_" + arch.Tag + "_input_corr <= #1 " + op.lqName + "_" + arch.Tag + "_output_z;\n"
				result += "								" + op.lqName + "_" + arch.Tag + "_state <= #1 " + op.lqName + "_" + arch.Tag + "_get;\n"
				result += "							end\n"
			}
			result += "							" + op.lqName + "_" + arch.Tag + "_get : begin\n"
			if op.opType == LQADD {
				result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + op.lqName + "_" + arch.Tag + "_output_z;\n"
			} else {
				result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + op.lqName + "_" + arch.Tag + "_output_corr;\n"
			}
			result += "								" + op.lqName + "_" + arch.Tag + "_state <= #1 " + op.lqName + "_" + arch.Tag + "_put;\n"
			result += "								_pc <= #1 _pc + 1'b1 ;\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"" + strings.ToUpper(op.lqName) + " " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op LinearQuantizer) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op LinearQuantizer) Assembler(arch *Arch, words []string) (string, error) {
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

func (op LinearQuantizer) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op LinearQuantizer) Simulate(vm *VM, instr string) error {
	return errors.New("unimplemented LinearQuantizer simulation")
}

// The random genaration does nothing
func (op LinearQuantizer) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op LinearQuantizer) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op LinearQuantizer) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op LinearQuantizer) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op LinearQuantizer) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op LinearQuantizer) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op LinearQuantizer) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op LinearQuantizer) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op LinearQuantizer) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	result := "\n\n"
	var moduleName string

	result += "module " + op.lqName + "_" + arch.Tag + "(\n"
	result += "        input_a,\n"
	result += "        input_b,\n"
	result += "        output_z);\n"
	result += "\n"
	result += "  input signed    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_a;\n"
	result += "  input signed    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_b;\n"
	result += "  output signed   [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] output_z;\n"
	switch op.opType {
	case LQADD:
		result += "  assign output_z = input_a + input_b;\n"
		moduleName = "adder"
	case LQMULT:
		result += "  assign output_z = input_a * input_b;\n"
		moduleName = "multiplier"
	case LQDIV:
		result += "  assign output_z = input_a / input_b;\n"
		moduleName = "divider"
	}
	result += "\n"
	result += "endmodule\n"

	moduleNames := []string{moduleName}
	moduleCodes := []string{result}

	if op.opType == LQDIV || op.opType == LQMULT {
		var correctionName string
		correction := "\n\n"
		correction += "module " + op.lqName + "_correction_" + arch.Tag + "(\n"
		correction += "        input_corr,\n"
		correction += "        output_corr);\n"
		correction += "\n"
		correction += "  input signed    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_corr;\n"
		correction += "  output signed   [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] output_corr;\n"
		sd := float64(uint64(1) << (arch.Rsize - 1))
		sn := float64(op.max)
		s := sd / sn
		corr := fmt.Sprint(uint(s))
		correction += "  parameter signed [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] CORRECTION = " + strconv.Itoa(int(arch.Rsize)) + "'d" + corr + ";\n"
		switch op.opType {
		case LQMULT:
			correction += "  assign output_corr = input_corr / CORRECTION;\n"
			correctionName = "multiplier_correction"
		case LQDIV:
			correction += "  assign output_corr = input_corr * CORRECTION;\n"
			correctionName = "divider_correction"
		}
		correction += "\n"
		correction += "endmodule\n"
		moduleNames = append(moduleNames, correctionName)
		moduleCodes = append(moduleCodes, correction)
	}

	return moduleNames, moduleCodes
}

func (Op LinearQuantizer) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op LinearQuantizer) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op LinearQuantizer) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, Op.lqName+"::*--type=reg::*--type=reg")
	return result
}
func (Op LinearQuantizer) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case Op.lqName:
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: Op.lqName, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:" + Op.lqName, T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:" + Op.lqName, T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}

func (Op LinearQuantizer) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
