package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The FloPoCo opcode is both a basic instruction and a template for other instructions.
type FloPoCo struct {
	floPoCoName string
	regSize     int
	vHDL        string
	entities    []string
	topEntity   string
	pipeline    int
}

func (op FloPoCo) Op_get_name() string {
	return op.floPoCoName
}

func (op FloPoCo) Op_get_desc() string {
	return "FloPoCo dynamical instruction " + op.floPoCoName
}

func (op FloPoCo) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := op.floPoCoName + " [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op FloPoCo) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op FloPoCo) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	regSizeLS := strconv.Itoa(op.regSize - 1)
	result := ""
	result += "\treg [" + regSizeLS + ":0] " + op.floPoCoName + "_" + arch.Tag + "_input_a;\n"
	result += "\treg [" + regSizeLS + ":0] " + op.floPoCoName + "_" + arch.Tag + "_input_b;\n"
	result += "\n"

	result += "\twire [" + regSizeLS + ":0] " + op.floPoCoName + "_" + arch.Tag + "_output_z;\n"
	result += "\n"

	result += "\treg	[1:0] " + op.floPoCoName + "_" + arch.Tag + "_state;\n"
	result += "\tparameter " + op.floPoCoName + "_" + arch.Tag + "_put_inputs      = 2'd0,\n"
	result += "\t          " + op.floPoCoName + "_" + arch.Tag + "_wait_pipeline   = 2'd1,\n"
	result += "\t          " + op.floPoCoName + "_" + arch.Tag + "_get_out         = 2'd2;\n"
	result += "\n"

	pipBits := Needed_bits(op.pipeline + 1)
	result += "\treg	[" + strconv.Itoa(pipBits-1) + ":0] " + op.floPoCoName + "_" + arch.Tag + "_pipeline;\n"
	result += "\n"

	result += "\tcp" + arch.Tag + "_" + op.floPoCoName + " cp" + arch.Tag + "_" + op.floPoCoName + "_inst (\n"
	result += "\t\t.clk(clock_signal),\n"
	result += "\t\t.rst(reset_signal),\n"
	result += "\t\t.X(" + op.floPoCoName + "_" + arch.Tag + "_input_a" + "),\n"
	result += "\t\t.Y(" + op.floPoCoName + "_" + arch.Tag + "_input_b" + "),\n"
	result += "\t\t.R(" + op.floPoCoName + "_" + arch.Tag + "_output_z" + ")\n"
	result += "\t);"

	return result
}

func (op FloPoCo) Op_instruction_verilog_state_machine(arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	pipBits := Needed_bits(op.pipeline + 1)
	pipelineS := strconv.Itoa(op.pipeline)

	reg_num := 1 << arch.R

	result := ""
	result += "					" + strings.ToUpper(op.floPoCoName) + ": begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {
			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							case (" + op.floPoCoName + "_" + arch.Tag + "_state)\n"
			result += "							" + op.floPoCoName + "_" + arch.Tag + "_put_inputs : begin\n"
			result += "								" + op.floPoCoName + "_" + arch.Tag + "_input_a <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
			result += "								" + op.floPoCoName + "_" + arch.Tag + "_input_b <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "								" + op.floPoCoName + "_" + arch.Tag + "_state <= #1 " + op.floPoCoName + "_" + arch.Tag + "_wait_pipeline;\n"
			result += "								" + op.floPoCoName + "_" + arch.Tag + "_pipeline <= #1 " + strconv.Itoa(pipBits) + "'d" + pipelineS + ";\n"
			result += "							end\n"
			result += "							" + op.floPoCoName + "_" + arch.Tag + "_wait_pipeline : begin\n"
			result += "								if (" + op.floPoCoName + "_" + arch.Tag + "_pipeline == " + strconv.Itoa(pipBits) + "'d0) begin\n"
			result += "									" + op.floPoCoName + "_" + arch.Tag + "_state <= #1 " + op.floPoCoName + "_" + arch.Tag + "_get_out;\n"
			result += "								end\n"
			result += "								else begin\n"
			result += "									" + op.floPoCoName + "_" + arch.Tag + "_pipeline <= #1 " + op.floPoCoName + "_" + arch.Tag + "_pipeline - 1;\n"
			result += "								end\n"
			result += "							end\n"
			result += "							" + op.floPoCoName + "_" + arch.Tag + "_get_out : begin\n"
			result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + op.floPoCoName + "_" + arch.Tag + "_output_z;\n"
			result += "								" + op.floPoCoName + "_" + arch.Tag + "_state <= #1 " + op.floPoCoName + "_" + arch.Tag + "_put_inputs;\n"
			result += "								_pc <= #1 _pc + 1'b1 ;\n"
			result += "							end\n"
			result += "							endcase\n"
			result += "								$display(\"" + op.floPoCoName + " " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "							end\n"
		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op FloPoCo) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op FloPoCo) Assembler(arch *Arch, words []string) (string, error) {
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

func (op FloPoCo) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op FloPoCo) Simulate(vm *VM, instr string) error {
	return errors.New("unimplemented FloPoCo simulation")
}

// The random genaration does nothing
func (op FloPoCo) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op FloPoCo) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op FloPoCo) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op FloPoCo) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op FloPoCo) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FloPoCo) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op FloPoCo) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FloPoCo) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op FloPoCo) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op FloPoCo) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op FloPoCo) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op FloPoCo) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, Op.floPoCoName+"::*--type=reg::*--type=reg")
	return result
}
func (Op FloPoCo) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case Op.floPoCoName:
		regDst := line.Elements[0].GetValue()
		regSrc := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDst, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regSrc, Op: bmreqs.OpAdd})
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}

func (Op FloPoCo) ExtraFiles(arch *Arch) ([]string, []string) {
	vHDL := Op.vHDL
	for _, ent := range Op.entities {
		if ent == Op.topEntity {
			vHDL = strings.ReplaceAll(vHDL, ent, "cp"+arch.Tag+"_"+Op.floPoCoName)
		} else {
			vHDL = strings.ReplaceAll(vHDL, ent, "cp"+arch.Tag+"_"+ent)
		}
	}
	return []string{"cp" + arch.Tag + "_" + Op.floPoCoName + ".vhd"}, []string{vHDL}
}
