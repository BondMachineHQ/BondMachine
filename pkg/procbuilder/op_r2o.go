package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type R2o struct{}

func (op R2o) Op_get_name() string {
	return "r2o"
}

func (op R2o) Op_get_desc() string {
	return "Register to output"
}

func (op R2o) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	outbits := arch.Outputs_bits()
	result := "r2o [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(outbits) + "(Output)]	// Set an output to the value of a given register [" + strconv.Itoa(opbits+int(arch.R)+outbits) + "]\n"
	return result
}

func (op R2o) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	outbits := arch.Outputs_bits()
	return opbits + int(arch.R) + int(outbits) // The bits for the opcode + bits for a register + bits for the output
}

func (op R2o) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

	// The extra process will be created and handled by the first waiting operator present in the architecture
	setflag := true

	if setflag {
		// It the first operator

		opbits := arch.Opcodes_bits()
		rom_word := arch.Max_word()

		result += "\n"
		// Data valid for outputs
		for j := 0; j < int(arch.M); j++ {
			result += "\treg " + strings.ToLower(Get_output_name(j)) + "_val;\n"
		}

		result += "\n"

		for j := 0; j < int(arch.M); j++ {

			objects := make([]string, 1)
			objects[0] = strings.ToLower(Get_output_name(j))

			// Process for data outputs data valid
			result += "\talways @(posedge clock_signal, posedge reset_signal)\n"
			result += "\tbegin\n"

			result += "\t\tif (reset_signal)\n"
			result += "\t\tbegin\n"
			result += "\t\t\t" + strings.ToLower(Get_output_name(j)) + "_val <= #1 1'b0;\n"
			result += "\t\tend\n"
			result += "\t\telse\n"
			result += "\t\tbegin\n"

			if opbits == 1 {
				result += "\t\t\tcase(rom_value[" + strconv.Itoa(rom_word-1) + "])\n"
			} else {
				result += "\t\t\tcase(rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "])\n"
			}

			for _, currop := range arch.Op {
				result += currop.Op_instruction_verilog_extra_block(arch, flavor, uint8(4), "output_data_valid", objects)
			}

			result += "\t\t\t\tdefault: begin\n"
			result += "\t\t\t\t\tif (" + strings.ToLower(Get_output_name(j)) + "_received)\n"
			result += "\t\t\t\t\tbegin\n"
			result += "\t\t\t\t\t\t" + strings.ToLower(Get_output_name(j)) + "_val <= #1 1'b0;\n"
			result += "\t\t\t\t\tend\n"
			result += "\t\t\t\tend\n"

			result += "\t\t\tendcase\n"

			result += "\t\tend\n"

			result += "\tend\n"
		}
	}

	return result
}

func (op R2o) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	outbits := arch.Outputs_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					R2O: begin\n"
	if arch.M > 0 {
		if arch.R == 1 {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
		} else {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if outbits == 1 {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
			} else {
				result += "							case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(outbits)) + "])\n"
			}

			for j := 0; j < int(arch.M); j++ {
				result += "							" + strings.ToUpper(Get_output_name(j)) + " : begin\n"
				result += "								_aux" + strings.ToLower(Get_output_name(j)) + " <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
				result += "								$display(\"R2O " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_output_name(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
	}
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"
	return result
}

func (op R2o) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op R2o) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	outbits := arch.Outputs_bits()
	rom_word := arch.Max_word()

	reg_num := 2
	reg_num = reg_num << (arch.R - 1)

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

	if partial, err := Process_output(words[1], int(arch.M)); err == nil {
		result += zeros_prefix(outbits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + outbits; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op R2o) Disassembler(arch *Arch, instr string) (string, error) {
	outbits := arch.Outputs_bits()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	inp_id := get_id(instr[arch.R : int(arch.R)+outbits])
	result += strings.ToLower(Get_output_name(inp_id))
	return result, nil
}

func (op R2o) Simulate(vm *VM, instr string) error {
	outbits := vm.Mach.Outputs_bits()
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	inp := get_id(instr[reg_bits : int(reg_bits)+outbits])
	vm.Outputs[inp] = vm.Registers[reg]
	vm.Pc = vm.Pc + 1
	return nil
}

func (op R2o) Generate(arch *Arch) string {
	outbits := arch.Outputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.M))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(outbits, get_binary(inp))
}

func (op R2o) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op R2o) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2o) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op R2o) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2o) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op R2o) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2o) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2o) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2o) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_OUTPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "r2o", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1

		for i, _ := range seq1 {
			result[i+2] = UsageNotify{C_OUTPUT, S_NIL, i + 1}
		}

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")

}

func (Op R2o) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	opbits := arch.Opcodes_bits()
	outbits := arch.Outputs_bits()
	rom_word := arch.Max_word()

	result := ""

	pref := strings.Repeat("\t", int(level))

	switch blockname {
	case "output_data_valid":
		result += pref + "R2O: begin\n"
		if outbits == 1 {
			result += pref + "\tcase (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += pref + "\tcase (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(outbits)) + "])\n"
		}

		result += pref + "\t" + strings.ToUpper(objects[0]) + " : begin\n"
		result += pref + "\t\t" + strings.ToLower(objects[0]) + "_val <= 1'b1;\n"
		result += pref + "\tend\n"

		result += pref + "\tdefault: begin\n"
		result += pref + "\t\tif (" + strings.ToLower(objects[0]) + "_received)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_val <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\tendcase\n"

		result += pref + "end\n"
	default:
		result = ""
	}
	return result
}
func (Op R2o) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "r2o::*--type=reg::*--type=output")
	result = append(result, "mov--iomode=async::*--type=output::*--type=reg")
	return result
}
func (Op R2o) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2o":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		outNeed := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "outputs", Value: outNeed, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		outVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "outputs", Value: outVal, Op: bmreqs.OpAdd})
		if regVal != "" && outVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2o")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(outVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "output")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2o) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
