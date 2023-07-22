package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The R2mri opcode is both a basic instruction and a template for other instructions.
type R2mri struct{}

func (op R2mri) Op_get_name() string {
	return "r2mri"
}

func (op R2mri) Op_get_desc() string {
	return "Copy a register value to the ram"
}

func (op R2mri) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "r2mri [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.L)) + "(RAM address)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opbits+int(arch.R)+int(arch.L)) + "]\n"
	return result
}

func (op R2mri) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.L) // The bits for the opcode + bits for a register + bits for RAM address
}

func (op R2mri) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.L)-1) + ":0] addr_ram_r2mri;\n"
	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] ram_din_i;\n"
	result += "\treg wr_int_ram;\n"
	return result
}

func (Op R2mri) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\twr_int_ram <= #1 1'b0;\n"
	//result += "\t\t\taddr_ram_r2mri <= #1 'b0;\n"
	//result += "\t\t\tram_din_i <= #1 'b0;\n"
	return result
}

func (op R2mri) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	//rom_word := arch.Max_word()
	//opbits := arch.Opcodes_bits()

	//reg_num := 1 << arch.R

	result := ""
	result += "					R2MRI: begin\n"
	/*if arch.R == 1 {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							wr_int_ram <= #1 1'b1;\n"
		result += "							addr_ram_r2mri <= #1 rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.Rsize)) + "];\n"
		result += "							ram_din_i <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "							$display(\"R2MRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "						end\n"
	}
	result += "						endcase\n"*/
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"

	return result
}

func (op R2mri) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "m2r" {
			setflag = false
			break
		} else if currop.Op_get_name() == "r2mri" {
			break
		}
	}
	if setflag {
		result += "\tassign ram_din = ram_din_i;\n"
		result += "\tassign ram_addr = (rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2R) ? addr_ram_m2r : addr_ram_r2mri;\n"
		result += "\tassign ram_wren = wr_int_ram;\n"
		result += "\tassign ram_en = 1'b1;\n"
	}

	result += "\talways @(rom_value"
	for i := 0; i < reg_num; i++ {
		result += ",_" + strings.ToLower(Get_register_name(i))
	}
	result += ")\n"

	result += "\tbegin\n"

	if opbits == 1 {
		result += "		if(rom_value[" + strconv.Itoa(rom_word-1) + "] == R2MRI) begin\n"
	} else {
		result += "		if(rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "] == R2MRI) begin\n"
	}

	if arch.R == 1 {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "				" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "					wr_int_ram <= 1'b1;\n"
		result += "					addr_ram_r2mri <= rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.L)) + "];\n"
		result += "					ram_din_i <= _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "					$display(\"R2MRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "				end\n"
	}
	result += "			endcase\n"
	result += "\t	end\n"
	result += "\t	else\n"
	result += "			wr_int_ram <= 1'b0;\n"
	result += "\tend\n"

	return result
}

func (op R2mri) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	ramdepth := int(arch.L)

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

	if partial, err := Process_number(words[1]); err == nil {
		result += zeros_prefix(ramdepth, partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.R) + ramdepth; i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op R2mri) Disassembler(arch *Arch, instr string) (string, error) {
	ramdepth := int(arch.L)
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	value := get_id(instr[arch.R : int(arch.R)+ramdepth])
	result += strconv.Itoa(value)
	return result, nil
}

// The simulation does nothing
func (op R2mri) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op R2mri) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2mri) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2mri) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2mri) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2mri) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\t\twr_int_ram <= #1 1'b0;\n"
	return result
}

func (Op R2mri) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2mri) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2mri) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "r2mri", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op R2mri) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2mri) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "r2mri::*--type=reg::*--type=ram--ramaddressing=register")
	result = append(result, "mov::*--type=ram--ramaddressing=register::*--type=reg")
	return result
}
func (Op R2mri) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2mri":
		regVal := line.Elements[0].GetValue()
		ramVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: ramVal, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[1].GetValue()
		ramVal := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: ramVal, Op: bmreqs.OpAdd})
		if regVal != "" && ramVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2mri")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(ramVal)
			newArg1.BasmMeta = newArg1.SetMeta("type", "reg")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2mri) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2mri) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
