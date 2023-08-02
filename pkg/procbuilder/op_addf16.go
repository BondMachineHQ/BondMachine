package procbuilder

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"

	"github.com/x448/float16"
)

// The Addf16 opcode is both a basic instruction and a template for other instructions.
type Addf16 struct{}

func (op Addf16) Op_get_name() string {
	return "addf16"
}

func (op Addf16) Op_get_desc() string {
	return "Register addf16"
}

func (op Addf16) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "addf16 [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Addf16) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Addf16) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [15:0] adder_" + arch.Tag + "_input_a;\n"
	result += "\treg [15:0] adder_" + arch.Tag + "_input_b;\n"
	result += "\treg adder_" + arch.Tag + "_input_a_stb;\n"
	result += "\treg adder_" + arch.Tag + "_input_b_stb;\n"
	result += "\treg adder_" + arch.Tag + "_output_z_ack;\n\n"

	result += "\twire [15:0] adder_" + arch.Tag + "_output_z;\n"
	result += "\twire adder_" + arch.Tag + "_output_z_stb;\n"
	result += "\twire adder_" + arch.Tag + "_input_a_ack;\n"
	result += "\twire adder_" + arch.Tag + "_input_b_ack;\n\n"

	result += "\treg	[1:0] adder_" + arch.Tag + "_state;\n"
	result += "parameter adder_" + arch.Tag + "_put_a         = 2'd0,\n"
	result += "          adder_" + arch.Tag + "_put_b         = 2'd1,\n"
	result += "          adder_" + arch.Tag + "_get_z         = 2'd2;\n"

	result += "\tadder_" + arch.Tag + " adder_" + arch.Tag + "_inst (adder_" + arch.Tag + "_input_a, adder_" + arch.Tag + "_input_b, adder_" + arch.Tag + "_input_a_stb, adder_" + arch.Tag + "_input_b_stb, adder_" + arch.Tag + "_output_z_ack, clock_signal, reset_signal, adder_" + arch.Tag + "_output_z, adder_" + arch.Tag + "_output_z_stb, adder_" + arch.Tag + "_input_a_ack, adder_" + arch.Tag + "_input_b_ack);\n\n"

	return result
}

func (op Addf16) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					ADDF16: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:addf16", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
			if req.Value == "false" {
				continue
			}
		}

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {

			if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlySrcRegs)) {
				cp := arch.Tag
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:addf16", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (adder_" + arch.Tag + "_state)\n"
			result += "							adder_" + arch.Tag + "_put_a : begin\n"
			result += "								if (adder_" + arch.Tag + "_input_a_ack) begin\n"
			result += "									adder_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "									adder_" + arch.Tag + "_input_a_stb <= #1 1;\n"
			result += "									adder_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									adder_" + arch.Tag + "_state <= #1 adder_" + arch.Tag + "_put_b;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							adder_" + arch.Tag + "_put_b : begin\n"
			result += "								if (adder_" + arch.Tag + "_input_b_ack) begin\n"
			result += "									adder_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "									adder_" + arch.Tag + "_input_b_stb <= #1 1;\n"
			result += "									adder_" + arch.Tag + "_output_z_ack <= #1 0;\n"
			result += "									adder_" + arch.Tag + "_state <= #1 adder_" + arch.Tag + "_get_z;\n"
			result += "									adder_" + arch.Tag + "_input_a_stb <= #1 0;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							adder_" + arch.Tag + "_get_z : begin\n"
			result += "								if (adder_" + arch.Tag + "_output_z_stb) begin\n"
			result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 adder_" + arch.Tag + "_output_z;\n"
			result += "									adder_" + arch.Tag + "_output_z_ack <= #1 1;\n"
			result += "									adder_" + arch.Tag + "_state <= #1 adder_" + arch.Tag + "_put_a;\n"
			result += "									adder_" + arch.Tag + "_input_b_stb <= #1 0;\n"
			result += NextInstruction(conf, arch, 9, "_pc + 1'b1")
			result += "								end\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"ADDF16 " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Addf16) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Addf16) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Addf16) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Addf16) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regDest := get_id(instr[:reg_bits])
	regSrc := get_id(instr[reg_bits : reg_bits*2])
	switch vm.Mach.Rsize {
	case 16:
		var floatDest float16.Float16
		var floatSrc float16.Float16
		if v, ok := vm.Registers[regDest].(uint16); ok {
			floatDest = float16.Frombits(v)
		} else {
			floatDest = float16.Float16(0)
		}
		if v, ok := vm.Registers[regSrc].(uint16); ok {
			floatSrc = float16.Frombits(v)
		} else {
			floatSrc = float16.Float16(0)
		}
		operation := floatDest.Float32() + floatSrc.Float32()
		vm.Registers[regDest] = float16.Fromfloat32(operation).Bits()
	default:
		return errors.New("invalid register size, for float registers has to be 16 bits")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Addf16) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Addf16) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Addf16) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Addf16) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Addf16) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addf16) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Addf16) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addf16) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Addf16) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	tmpl := arch.createBasicTemplateData()
	tmpl.ModuleName = "adder_" + arch.Tag
	var f bytes.Buffer
	t, _ := template.New("adder").Funcs(tmpl.funcmap).Parse(addf16)
	t.Execute(&f, *tmpl)
	result := f.String()
	return []string{"adder"}, []string{result}
}

func (Op Addf16) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Addf16) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Addf16) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "addf16::*--type=reg::*--type=reg")
	return result
}
func (Op Addf16) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "addf16":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "addf16", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:addf16", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:addf16", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Addf16) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Addf16) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
