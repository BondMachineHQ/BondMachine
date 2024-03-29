package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type I2r struct{}

func (op I2r) Op_get_name() string {
	return "i2r"
}

func (op I2r) Op_get_desc() string {
	return "Input to register"
}

func (op I2r) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	inBits := arch.Inputs_bits()
	result := "i2r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(inBits) + "(Input)]	// Set a register to the value of the given input [" + strconv.Itoa(opBits+int(arch.R)+inBits) + "]\n"
	return result
}

func (op I2r) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()
	return opbits + int(arch.R) + int(inpbits) // The bits for the opcode + bits for a register + bits for the input
}

func (op I2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""

	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	result += "\n"
	// Data received fro inputs
	for j := 0; j < int(arch.N); j++ {
		result += "\treg " + strings.ToLower(Get_input_name(j)) + "_recv;\n"
	}

	result += "\n"

	for j := 0; j < int(arch.N); j++ {

		objects := make([]string, 1)
		objects[0] = strings.ToLower(Get_input_name(j))

		// Process for data outputs data valid
		result += "\talways @(posedge clock_signal, posedge reset_signal)\n"
		result += "\tbegin\n"

		result += "\t\tif (reset_signal)\n"
		result += "\t\tbegin\n"
		result += "\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
		result += "\t\tend\n"
		result += "\t\telse\n"
		result += "\t\tbegin\n"

		if opbits == 1 {
			result += "\t\t\tcase(current_instruction[" + strconv.Itoa(rom_word-1) + "])\n"
		} else {
			result += "\t\t\tcase(current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "])\n"
		}

		for _, currop := range arch.Op {
			result += currop.Op_instruction_verilog_extra_block(arch, flavor, uint8(4), "input_data_received", objects)
		}

		result += "\t\t\t\tdefault: begin\n"
		result += "\t\t\t\t\tif (!" + strings.ToLower(Get_input_name(j)) + "_valid)\n"
		result += "\t\t\t\t\tbegin\n"
		result += "\t\t\t\t\t\t" + strings.ToLower(Get_input_name(j)) + "_recv <= #1 1'b0;\n"
		result += "\t\t\t\t\tend\n"
		result += "\t\t\t\tend\n"

		result += "\t\t\tendcase\n"

		result += "\t\tend\n"

		result += "\tend\n"
	}

	return result
}

func (op I2r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "					I2R: begin\n"
	if arch.N > 0 {
		if arch.R == 1 {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
		} else {
			result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
		}
		for i := 0; i < reg_num; i++ {
			result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

			if inpbits == 1 {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
			} else {
				result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(inpbits)) + "])\n"
			}

			for j := 0; j < int(arch.N); j++ {
				result += "							" + strings.ToUpper(Get_input_name(j)) + " : begin\n"
				result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 " + strings.ToLower(Get_input_name(j)) + ";\n"
				result += "								$display(\"I2R " + strings.ToUpper(Get_register_name(i)) + " " + strings.ToUpper(Get_input_name(j)) + "\");\n"
				result += "							end\n"

			}
			result += "							endcase\n"
			result += "						end\n"
		}
		result += "						endcase\n"
	} else {
		result += "						$display(\"NOP\");\n"
	}
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op I2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op I2r) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	inpbits := arch.Inputs_bits()
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

	if partial, err := Process_input(words[1], int(arch.N)); err == nil {
		result += zeros_prefix(inpbits, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + inpbits; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op I2r) Disassembler(arch *Arch, instr string) (string, error) {
	inpbits := arch.Inputs_bits()
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	inp_id := get_id(instr[arch.R : int(arch.R)+inpbits])
	result += strings.ToLower(Get_input_name(inp_id))
	return result, nil
}

func (op I2r) Simulate(vm *VM, instr string) error {
	inpbits := vm.Mach.Inputs_bits()
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	inp := get_id(instr[reg_bits : int(reg_bits)+inpbits])
	vm.Registers[reg] = vm.Inputs[inp]
	vm.Pc = vm.Pc + 1
	return nil
}

func (op I2r) Generate(arch *Arch) string {
	inpbits := arch.Inputs_bits()
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	inp := rand.Intn(int(arch.N))
	return zeros_prefix(int(arch.R), get_binary(reg)) + zeros_prefix(inpbits, get_binary(inp))
}

func (op I2r) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op I2r) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op I2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op I2r) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op I2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op I2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op I2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op I2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op I2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "i2r", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1

		for i, _ := range seq1 {
			result[i+2] = UsageNotify{C_INPUT, S_NIL, i + 1}
		}

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op I2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	// TODO: ADDI and I2R are to be fixed in order to coexist
	opbits := arch.Opcodes_bits()
	inbits := arch.Inputs_bits()
	rom_word := arch.Max_word()

	result := ""

	pref := strings.Repeat("\t", int(level))

	switch blockname {
	case "input_data_received":
		result += pref + "I2R: begin\n"
		if inbits == 1 {
			result += pref + "\tcase (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += pref + "\tcase (current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(inbits)) + "])\n"
		}

		result += pref + "\t" + strings.ToUpper(objects[0]) + " : begin\n"
		result += pref + "\t\tif (" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b1;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\tdefault: begin\n"
		result += pref + "\t\tif (!" + strings.ToLower(objects[0]) + "_valid)\n"
		result += pref + "\t\tbegin\n"
		result += pref + "\t\t\t" + strings.ToLower(objects[0]) + "_recv <= #1 1'b0;\n"
		result += pref + "\t\tend\n"
		result += pref + "\tend\n"

		result += pref + "\tendcase\n"

		result += pref + "end\n"
	default:
		result = ""
	}
	return result
}
func (Op I2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "i2r::*--type=reg::*--type=input")
	result = append(result, "mov--iomode=async::*--type=reg::*--type=input")
	return result
}
func (Op I2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "i2r":
		regNeed := line.Elements[0].GetValue()
		inNeed := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "inputs", Value: inNeed, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[0].GetValue()
		inVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "inputs", Value: inVal, Op: bmreqs.OpAdd})
		if regVal != "" && inVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("i2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(inVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "input")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op I2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op I2r) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
