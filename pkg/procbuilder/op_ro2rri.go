package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Ro2rri struct{}

func (op Ro2rri) Op_get_name() string {
	return "ro2rri"
}

func (op Ro2rri) Op_get_desc() string {
	return "ROM to register"
}

func (op Ro2rri) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := "ro2rri [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Set a register to the value of the given ROM location [" + strconv.Itoa(opBits+int(arch.R)+int(arch.O)) + "]\n"
	return result
}

func (op Ro2rri) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	return opBits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for a register
}

func (op Ro2rri) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {

	result := ""

	// Check if the romread facility has already been included
	romReadFlag := conf.Runinfo.Check("romread")

	// If not, include it
	if romReadFlag {

		romWord := arch.Max_word()

		result += "\twire [" + strconv.Itoa(int(romWord-1)) + ":0] romread_value;\n"
		result += "\treg [" + strconv.Itoa(int(arch.O)-1) + ":0] romread_bus;\n"
		result += "\treg romread_ready;\n"
		result += "\n"
		result += "\t" + pname + "rom romread_instance(romread_bus,romread_value);\n"

	}

	return result
}

func (op Ro2rri) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()

	regNum := 1 << arch.R
	rSize := int(arch.Rsize)

	// If the WORD size is bigger than the register size, we need to slice the ROM value (that is left padded with zeros)
	sliceROM := "romread_value[" + strconv.Itoa(int(arch.Rsize)-1) + ":0]"
	sliceReg := "[" + strconv.Itoa(int(arch.Rsize)-1) + ":0]"

	// If the register size is bigger than the WORD size, we need to add zeros to the ROM value
	// TODO: unchecked code, check if it works
	if int(arch.Rsize) > romWord {
		sliceROM = "{ " + strconv.Itoa(int(arch.Rsize)-romWord) + "'b0 ,romread_value[" + strconv.Itoa(romWord-1) + ":0] }"
	}

	result := ""
	result += "					RO2RRI: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + "])\n"
		} else {
			result += "							case (current_instruction[" + strconv.Itoa(romWord-opBits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opBits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < regNum; j++ {
			result += "							" + strings.ToUpper(Get_register_name(j)) + " : begin\n"

			result += "								if (romread_ready == 1'b1) begin\n"
			result += "									_" + strings.ToLower(Get_register_name(i)) + sliceReg + " <= #1 " + sliceROM + ";\n"
			result += "									romread_ready <= 1'b0;\n"
			result += NextInstruction(conf, arch, 9, "_pc + 1'b1")
			result += "								end\n"
			result += "								else begin\n"
			if arch.O <= arch.Rsize {
				result += "									romread_bus[" + strconv.Itoa(int(arch.O)-1) + ":0] <= _" + strings.ToLower(Get_register_name(j)) + "[" + strconv.Itoa(int(arch.O)-1) + ":0];\n"
			} else {
				// TODO: unchecked code, check if it works
				result += "									romread_bus[" + strconv.Itoa(int(arch.O)-1) + ":0] <= { " + strconv.Itoa(int(arch.O)-rSize) + "'b0 ,_" + strings.ToLower(Get_register_name(j)) + "[" + strconv.Itoa(int(arch.O)-1) + ":0] };\n"
			}
			result += "									romread_ready <= 1'b1;\n"
			result += "								end\n"
			result += "								$display(\"RO2RRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"

			result += "							end\n"

		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result

}

func (op Ro2rri) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Ro2rri) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	regNum := 1 << arch.R

	if len(words) != 2 {
		return "", Prerror{"wrong arguments number"}
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"unknown register name " + words[0]}
	}

	partial := ""
	for i := 0; i < regNum; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", Prerror{"unknown register name " + words[1]}
	}

	result += partial

	for i := opBits + 2*int(arch.R); i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Ro2rri) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op Ro2rri) Simulate(vm *VM, instr string) error {
	regBits := vm.Mach.R
	regDest := get_id(instr[:regBits])
	regSrc := get_id(instr[regBits : regBits*2])
	sliceSize := int(vm.Mach.Rsize)
	rSize := int(vm.Mach.Rsize)
	wordSize := len(vm.Mach.Program.Slocs[0])
	slocsLen := len(vm.Mach.Program.Slocs)

	// The slice size is the minimum between the register size and the word size
	if sliceSize > wordSize {
		sliceSize = wordSize
	}

	if rSize <= 8 {
		loc := int(vm.Registers[regSrc].(uint8))
		if loc < slocsLen {
			vm.Registers[regDest] = uint8(get_id(vm.Mach.Program.Slocs[loc][wordSize-sliceSize : wordSize]))
		} else {
			vm.Registers[regDest] = uint8(get_id(vm.Mach.Data.Vars[loc-slocsLen][wordSize-sliceSize : wordSize]))
		}
	} else if rSize <= 16 {
		loc := int(vm.Registers[regSrc].(uint16))
		if loc < slocsLen {
			vm.Registers[regDest] = uint16(get_id(vm.Mach.Program.Slocs[loc][wordSize-sliceSize : wordSize]))
		} else {
			vm.Registers[regDest] = uint16(get_id(vm.Mach.Data.Vars[loc-slocsLen][wordSize-sliceSize : wordSize]))
		}
	} else if rSize <= 32 {
		loc := int(vm.Registers[regSrc].(uint32))
		if loc < slocsLen {
			vm.Registers[regDest] = uint32(get_id(vm.Mach.Program.Slocs[loc][wordSize-sliceSize : wordSize]))
		} else {
			vm.Registers[regDest] = uint32(get_id(vm.Mach.Data.Vars[loc-slocsLen][wordSize-sliceSize : wordSize]))
		}
	} else if rSize <= 64 {
		loc := int(vm.Registers[regSrc].(uint64))
		if loc < slocsLen {
			vm.Registers[regDest] = uint64(get_id(vm.Mach.Program.Slocs[loc][wordSize-sliceSize : wordSize]))
		} else {
			vm.Registers[regDest] = uint64(get_id(vm.Mach.Data.Vars[loc-slocsLen][wordSize-sliceSize : wordSize]))
		}
	} else {
		return errors.New("invalid register size, must be <= 64")
	}

	vm.Pc = vm.Pc + 1
	return nil
}

func (op Ro2rri) Generate(arch *Arch) string {
	return ""
}

func (op Ro2rri) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Ro2rri) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Ro2rri) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Ro2rri) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2rri) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2rri) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2rri) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2rri) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Ro2rri) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "ro2rri", I_NIL}
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

func (Op Ro2rri) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Ro2rri) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "ro2rri::*--type=reg::*--type=reg")
	result = append(result, "mov::*--type=reg::*--type=rom--romaddressing=register")
	return result
}
func (Op Ro2rri) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "ro2rri":
		regNeed := line.Elements[1].GetValue()
		regDest := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDest, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regDest := line.Elements[0].GetValue()
		regNeed := line.Elements[1].GetMeta("romregister")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDest, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		if regDest != "" && regNeed != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("ro2rri")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regDest)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(regNeed)
			newArg1.BasmMeta = newArg1.SetMeta("type", "reg")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Ro2rri) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Ro2rri) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
