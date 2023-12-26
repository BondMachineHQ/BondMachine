package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Cmpr struct {
	pipeline *bool
}

func (op Cmpr) Op_get_name() string {
	return "cmpr"
}

func (op Cmpr) Op_get_desc() string {
	return "Register comparison"
}

func (op Cmpr) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "cmpr [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Compare two registers value setting a flag [" + strconv.Itoa(opBits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Cmpr) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Cmpr) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pName string) string {
	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"cmpr", "cmpv", "jcmpl", "jcmpo", "jcmpa", "jcmprio", "jcmpria", "jncmpl", "jncmpo", "jncmpa", "jncmprio", "jncmpria"}) {
		result += "\treg cmpflag;\n"
	}

	return result
}

func (op Cmpr) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opBits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					CMPR: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opBits-1) + ":" + strconv.Itoa(rom_word-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:cmpr", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
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
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:cmpr", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "								if (_" + strings.ToLower(Get_register_name(j)) + " == _" + strings.ToLower(Get_register_name(i)) + ") begin\n"
			result += "									cmpflag <= 1'b1;\n"
			result += "								end else begin\n"
			result += "									cmpflag <= 1'b0;\n"
			result += "								end\n"
			result += "								$display(\"CMPR " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Cmpr) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpr) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Cmpr) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Cmpr) Simulate(vm *VM, instr string) error {
	// TODO
	return nil
}

// The random generation does nothing
func (op Cmpr) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Cmpr) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Cmpr) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Cmpr) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Cmpr) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpr) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpr) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpr) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cmpr) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Cmpr) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (op Cmpr) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Cmpr) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "cmpr::*--type=reg::*--type=reg")
	return result
}
func (op Cmpr) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "cmpr":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "cmpr", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:cmpr", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:cmpr", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Cmpr) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Cmpr) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
