package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Ro2r struct{}

func (op Ro2r) Op_get_name() string {
	return "ro2r"
}

func (op Ro2r) Op_get_desc() string {
	return "ROM to register"
}

func (op Ro2r) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "ro2r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Set a register to the value of the given ROM location [" + strconv.Itoa(opbits+int(arch.R)+int(arch.O)) + "]\n"
	return result
}

func (op Ro2r) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.O) // The bits for the opcode + bits for a register + bits for the location
}

func (op Ro2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {

	result := ""

	// Check if the romread facility has already been included
	romreadflag := conf.Runinfo.Check("romread")

	// If not, include it
	if romreadflag {

		romWord := arch.Max_word()

		result += "\twire [" + strconv.Itoa(int(romWord-1)) + ":0] romread_value;\n"
		result += "\treg [" + strconv.Itoa(int(arch.O)-1) + ":0] romread_bus;\n"
		result += "\treg romread_ready;\n"
		result += "\n"
		result += "\t" + pname + "rom romread_instance(romread_bus,romread_value);\n"

	}

	return result
}

func (op Ro2r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	// TODO Incomplete
	romWord := arch.Max_word()
	opbits := arch.Opcodes_bits()

	regNum := 1 << arch.R

	result := ""
	result += "					RO2R: begin\n"
	if arch.R == 1 {
		result += "						case (rom_value[" + strconv.Itoa(romWord-opbits-1) + "])\n"
	} else {
		result += "						case (rom_value[" + strconv.Itoa(romWord-opbits-1) + ":" + strconv.Itoa(romWord-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < regNum; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							if (romread_ready == 1'b1) begin\n"
		result += "								_" + strings.ToLower(Get_register_name(i)) + " <= #1 romread_value[" + strconv.Itoa(romWord-1) + ":0];\n"
		result += "								romread_ready <= 1'b0;\n"
		result += "								_pc <= #1 _pc + 1'b1 ;\n"
		result += "							end\n"
		result += "							else begin\n"
		result += "								romread_bus[" + strconv.Itoa(int(arch.O)-1) + ":0] <= rom_value[" + strconv.Itoa(romWord-opbits-int(arch.R)-1) + ":" + strconv.Itoa(romWord-opbits-int(arch.R)-int(arch.O)) + "];\n"
		result += "								romread_ready <= 1'b1;\n"
		result += "							end\n"
		result += "							$display(\"RO2R " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result
}

func (op Ro2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Ro2r) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
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

	if partial, err := Process_number(words[1]); err == nil {
		result += zeros_prefix(int(arch.O), partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + int(arch.O); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Ro2r) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	value := get_id(instr[arch.R : arch.R+arch.O])
	result += strconv.Itoa(value)
	return result, nil
}

func (op Ro2r) Simulate(vm *VM, instr string) error {
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Ro2r) Generate(arch *Arch) string {
	return ""
}

func (op Ro2r) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Ro2r) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Ro2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Ro2r) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Ro2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Ro2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "ro2r", I_NIL}
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

func (Op Ro2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Ro2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "ro2r::*--type=rom::*--type=reg")
	result = append(result, "mov::*--type=reg::*--type=rom--romaddressing=immediate")
	return result
}
func (Op Ro2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "ro2r":
		regNeed := line.Elements[0].GetValue()
		location := line.Elements[1].GetMeta("location")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectMax, Name: "rom", Value: location, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regNeed := line.Elements[0].GetValue()
		location := line.Elements[1].GetMeta("location")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectMax, Name: "rom", Value: location, Op: bmreqs.OpAdd})
		if regNeed != "" && location != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("ro2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regNeed)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(location)
			newArg1.BasmMeta = newArg1.SetMeta("type", "rom")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Ro2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
