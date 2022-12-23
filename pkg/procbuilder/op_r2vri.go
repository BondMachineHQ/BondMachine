package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Add opcode is both a basic instruction and a template for other instructions.
type R2vri struct{}

func (op R2vri) Op_get_name() string {
	return "r2vri"
}

func (op R2vri) Op_get_desc() string {
	return "Register indirect copy to video RAM"
}

func (op R2vri) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "r2vri [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.R)) + "(Reg)]	// Set a register to the sum of its value with another register [" + strconv.Itoa(opbits+int(arch.R)+int(arch.R)) + "]\n"
	return result
}

func (op R2vri) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.R) // The bits for the opcode + bits for a register + bits for another register
}

func (op R2vri) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (Op R2vri) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op R2vri) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	//rom_word := arch.Max_word()
	//opbits := arch.Opcodes_bits()

	//reg_num := 1 << arch.R

	result := ""
	result += "					R2VRI: begin\n"
	/*if arch.R == 1 {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							wr_int_ram <= #1 1'b1;\n"
		result += "							addr_ram_r2m <= #1 rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.Rsize)) + "];\n"
		result += "							ram_din_i <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "							$display(\"R2M " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "						end\n"
	}
	result += "						endcase\n"*/
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"

	return result
}

func (op R2vri) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	// This is always the last module if present
	firstModule := true
	lastModule := true

	// If r2v is also present, that will be the last opcode
	for _, currOp := range arch.Op {
		if currOp.Op_get_name() == "r2v" {
			firstModule = false
			break
		}
	}
	//	result += "\tassign vram_din = vram_din_i;\n"
	//	result += "\tassign vram_addr = addr_vram_r2v;\n"
	//	result += "\tassign vram_wren = wr_int_vram;\n"
	//	result += "\tassign vram_en = 1'b1;\n"

	// Check for differences
	//	result += "\talways @(rom_value"
	//	for i := 0; i < reg_num; i++ {
	//		result += ",_" + strings.ToLower(Get_register_name(i))
	//}
	//result += ")\n"
	if firstModule {
		result += "\talways @(posedge clock_signal)\n"
		result += "\tbegin\n"
		if opbits == 1 {
			result += "\t\tcase (rom_value[" + strconv.Itoa(rom_word-1) + "])\n"
		} else {
			result += "\t\tcase (rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "])\n"
		}
	}

	result += "		R2VRI: begin\n"

	if arch.R == 1 {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "				" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		if arch.R == 1 {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + "])\n"
		} else {
			result += "						case (rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.R)) + "])\n"
		}

		for j := 0; j < reg_num; j++ {
			result += "						" + strings.ToUpper(Get_register_name(j)) + " : begin\n"
			result += "							vtm0_wren_i <= 1'b1;\n"
			result += "							vtm0_addr_i[7:0] <= _" + strings.ToLower(Get_register_name(j)) + "[7:0];\n"
			result += "							vtm0_din_i[7:0] <= _" + strings.ToLower(Get_register_name(i)) + "[7:0];\n"
			result += "							$display(\"R2VRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
			result += "						end\n"

		}
		result += "							endcase\n"

		result += "				end\n"
	}
	result += "			endcase\n"
	result += "\t	end\n"
	if lastModule {
		result += "\t	default:\n"
		result += "			vtm0_wren_i <= 1'b0;\n"
		result += "\t	endcase\n"
		result += "\tend\n"
	}

	return result
}

func (op R2vri) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	//ramdepth := 8

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

func (op R2vri) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

// The simulation does nothing
func (op R2vri) Simulate(vm *VM, instr string) error {
	// TODO
	reg_bits := vm.Mach.R
	regPay := get_id(instr[:reg_bits])
	regPos := get_id(instr[reg_bits : reg_bits*2])

	pos := vm.Registers[regPos].(uint8)
	payload := vm.Registers[regPay].(uint8)

	cmd := make([]byte, 0)

	cmd = append(cmd, byte(vm.CpID))
	cmd = append(cmd, byte(pos))
	cmd = append(cmd, byte(payload))

	vm.CmdChan <- cmd

	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op R2vri) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2vri) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2vri) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2vri) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2vri) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2vri) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2vri) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2vri) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_REGISTER {

		result := make([]UsageNotify, 3)
		newnot0 := UsageNotify{C_OPCODE, "r2vri", I_NIL}
		result[0] = newnot0
		newnot1 := UsageNotify{C_REGSIZE, S_NIL, len(seq0)}
		result[1] = newnot1
		newnot2 := UsageNotify{C_REGSIZE, S_NIL, len(seq1)}
		result[2] = newnot2

		return result, nil

	}

	return []UsageNotify{}, errors.New("Wrong parameters")
}

func (Op R2vri) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2vri) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "r2vri::*--type=reg::*--type=reg")
	result = append(result, "mov::*--type=somov--sotype=vtm--soaddressing=register::*--type=reg")
	return result
}

func (Op R2vri) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2vri":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		regNeed = line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		regAddr := line.Elements[0].GetMeta("soregister")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regAddr, Op: bmreqs.OpAdd})
		if regAddr != "" && regVal != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2vri")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(regAddr)
			newArg1.BasmMeta = newArg1.SetMeta("type", "reg")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2vri) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
