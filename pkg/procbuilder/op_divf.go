package procbuilder

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"strings"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Divf opcode is both a basic instruction and a template for other instructions.
type Divf struct{}

func (op Divf) Op_get_name() string {
	return "divf"
}

func (op Divf) Op_get_desc() string {
	return "Register divf"
}

func (op Divf) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "divf [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the division of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Divf) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Divf) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [31:0] divider_" + arch.Tag + "_input_a;\n"
	result += "\treg [31:0] divider_" + arch.Tag + "_input_b;\n"
	result += "\treg divider_" + arch.Tag + "_input_a_stb;\n"
	result += "\treg divider_" + arch.Tag + "_input_b_stb;\n"
	result += "\treg divider_" + arch.Tag + "_output_z_ack;\n\n"

	result += "\twire [31:0] divider_" + arch.Tag + "_output_z;\n"
	result += "\twire divider_" + arch.Tag + "_output_z_stb;\n"
	result += "\twire divider_" + arch.Tag + "_input_a_ack;\n"
	result += "\twire divider_" + arch.Tag + "_input_b_ack;\n\n"

	result += "\treg	[1:0] divider_" + arch.Tag + "_state;\n"
	result += "parameter divider_" + arch.Tag + "_put_a         = 2'd0,\n"
	result += "          divider_" + arch.Tag + "_put_b         = 2'd1,\n"
	result += "          divider_" + arch.Tag + "_get_z         = 2'd2;\n"

	result += "\tdivider_" + arch.Tag + " divider_" + arch.Tag + "_inst (divider_" + arch.Tag + "_input_a, divider_" + arch.Tag + "_input_b, divider_" + arch.Tag + "_input_a_stb, divider_" + arch.Tag + "_input_b_stb, divider_" + arch.Tag + "_output_z_ack, clock_signal, reset_signal, divider_" + arch.Tag + "_output_z, divider_" + arch.Tag + "_output_z_stb, divider_" + arch.Tag + "_input_a_ack, divider_" + arch.Tag + "_input_b_ack);\n\n"

	return result
}

func (op Divf) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					DIVF: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:divf", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:divf", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (divider_" + arch.Tag + "_state)\n"
			result += "							divider_" + arch.Tag + "_put_a : begin\n"
			result += "								if (divider_" + arch.Tag + "_input_a_ack) begin\n"
			result += "									divider_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "									divider_" + arch.Tag + "_input_a_stb <= #1 1;\n"
			result += "									divider_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									divider_" + arch.Tag + "_state <= #1 divider_" + arch.Tag + "_put_b;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							divider_" + arch.Tag + "_put_b : begin\n"
			result += "								if (divider_" + arch.Tag + "_input_b_ack) begin\n"
			result += "									divider_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "									divider_" + arch.Tag + "_input_b_stb <= #1 1;\n"
			result += "									divider_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									divider_" + arch.Tag + "_state <= #1 divider_" + arch.Tag + "_get_z;\n"
			result += "									divider_" + arch.Tag + "_input_a_stb <= #1 0;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							divider_" + arch.Tag + "_get_z : begin\n"
			result += "								if (divider_" + arch.Tag + "_output_z_stb) begin\n"
			result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 divider_" + arch.Tag + "_output_z;\n"
			result += "									divider_" + arch.Tag + "_output_z_ack <= #1 1;\n"
			result += "									divider_" + arch.Tag + "_state <= #1 divider_" + arch.Tag + "_put_a;\n"
			result += "									divider_" + arch.Tag + "_input_b_stb <= #1 0;\n"
			result += "									_pc <= #1 _pc + 1'b1 ;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"DIVF " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Divf) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Divf) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Divf) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The simulation does nothing
func (op Divf) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regDest := get_id(instr[:reg_bits])
	regSrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 32:
		var floatDest float32
		var floatSrc float32
		if v, ok := vm.Registers[regDest].(uint32); ok {
			floatDest = math.Float32frombits(v)
		} else {
			floatDest = float32(0.0)
		}
		if v, ok := vm.Registers[regSrc].(uint32); ok {
			floatSrc = math.Float32frombits(v)
		} else {
			floatSrc = float32(0.0)
		}
		vm.Registers[regDest] = math.Float32bits(floatDest / floatSrc)
	default:
		return errors.New("invalid register size, for float registers has to be 32 bits")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Divf) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Divf) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Divf) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Divf) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Divf) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Divf) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Divf) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Divf) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Divf) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	tmpl := arch.createBasicTemplateData()
	tmpl.ModuleName = "divider_" + arch.Tag
	var f bytes.Buffer
	t, _ := template.New("divider").Funcs(tmpl.funcmap).Parse(divf32)
	t.Execute(&f, *tmpl)
	result := f.String()
	return []string{"divider"}, []string{result}
}

func (Op Divf) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Divf) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Divf) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "divf::*--type=reg::*--type=reg")
	return result
}
func (Op Divf) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "divf":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "divf", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:divf", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:divf", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Divf) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
