package procbuilder

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Multp struct{}

func (op Multp) Op_get_name() string {
	return "multp"
}

func (op Multp) Op_get_desc() string {
	return "Register pipelined multiplcation"
}

func (op Multp) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "multp [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the product of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op Multp) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op Multp) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] multp_" + arch.Tag + "_input_a;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] multp_" + arch.Tag + "_input_b;\n"
	result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] multp_" + arch.Tag + "_output_z;\n"

	result += "\treg	[1:0] multp_" + arch.Tag + "_state;\n"
	result += "parameter multp_" + arch.Tag + "_put         = 2'd0,\n"
	result += "          multp_" + arch.Tag + "_get         = 2'd1;\n"

	result += "\tmultp_" + arch.Tag + " multp_" + arch.Tag + "_inst (multp_" + arch.Tag + "_input_a, multp_" + arch.Tag + "_input_b,  multp_" + arch.Tag + "_output_z);\n\n"

	return result
}

func (op Multp) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					MULTP: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		if IsHwOptimizationSet(conf.HwOptimizations, HwOptimizations(OnlyDestRegs)) {
			cp := arch.Tag
			req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:multp", T: bmreqs.ObjectSet, Name: "destregs", Value: Get_register_name(i), Op: bmreqs.OpCheck})
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
				req := rg.Requirement(bmreqs.ReqRequest{Node: "/bm:cps/id:" + cp + "/opcodes:multp", T: bmreqs.ObjectSet, Name: "sourceregs", Value: Get_register_name(j), Op: bmreqs.OpCheck})
				if req.Value == "false" {
					continue
				}
			}

			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (multp_" + arch.Tag + "_state)\n"
			result += "							multp_" + arch.Tag + "_put : begin\n"
			result += "								multp_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								multp_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								multp_" + arch.Tag + "_state <= #1 multp_" + arch.Tag + "_get;\n"
			result += "							end\n"
			result += "							multp" + arch.Tag + "_get : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 multp_" + arch.Tag + "_output_z;\n"
			result += "								multp_" + arch.Tag + "_state <= #1 multp_" + arch.Tag + "_put;\n"
			result += "								_pc <= #1 _pc + 1'b1 ;\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"MULTP " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Multp) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Multp) Assembler(arch *Arch, words []string) (string, error) {
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

func (op Multp) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Multp) Simulate(vm *VM, instr string) error {
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
		vm.Registers[regDest] = math.Float32bits(floatDest * floatSrc)
	default:
		return errors.New("invalid register size, for float registers has to be 32 bits")
	}
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Multp) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Multp) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Multp) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Multp) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Multp) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multp) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Multp) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multp) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Multp) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	result := "\n\n"

	result += "module multp_" + arch.Tag + "(\n"
	result += "        input_a,\n"
	result += "        input_b,\n"
	result += "        output_z);\n"
	result += "\n"
	result += "  input     [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_a;\n"
	result += "  input     [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] input_b;\n"
	result += "  output    [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] output_z;\n"
	result += "  assign output_z = input_a * input_b;\n"
	result += "\n"
	result += "endmodule\n"

	return []string{"adder"}, []string{result}
}

func (Op Multp) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Multp) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Multp) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "multp::*--type=reg::*--type=reg")
	return result
}
func (Op Multp) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "multp":
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "opcodes", Value: "multp", Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:multp", T: bmreqs.ObjectSet, Name: "destregs", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node + "/opcodes:multp", T: bmreqs.ObjectSet, Name: "sourceregs", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Multp) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
