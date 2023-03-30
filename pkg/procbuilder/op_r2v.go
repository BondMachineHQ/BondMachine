package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The R2v opcode
type R2v struct{}

func (op R2v) Op_get_name() string {
	return "r2v"
}

func (op R2v) Op_get_desc() string {
	return "Copy a register value to the video ram"
}

func (op R2v) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "r2v [" + strconv.Itoa(int(arch.R)) + "(Reg)] [ 8 (RAM address)]	// " + op.Op_get_desc() + " [" + strconv.Itoa(opbits+int(arch.R)+8) + "]\n"
	return result
}

func (op R2v) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + 8 // The bits for the opcode + bits for a register + bits for vRAM address
}

func (op R2v) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	//	result += "\treg [7:0] addr_vram_r2v;\n"
	//	result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] vram_din_i;\n"
	//	result += "\treg wr_int_vram;\n"
	return result
}

func (Op R2v) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\twr_int_ram <= #1 1'b0;\n"
	//result += "\t\t\taddr_ram_r2m <= #1 'b0;\n"
	//result += "\t\t\tram_din_i <= #1 'b0;\n"
	return result
}

func (op R2v) Op_instruction_verilog_state_machine(arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	//rom_word := arch.Max_word()
	//opbits := arch.Opcodes_bits()

	//reg_num := 1 << arch.R

	result := ""
	result += "					R2V: begin\n"
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

func (op R2v) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	// This is always the first module if present
	firstModule := true
	lastModule := true

	// If r2vri is also present, that will be the last opcode
	for _, currOp := range arch.Op {
		if currOp.Op_get_name() == "r2vri" {
			lastModule = false
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

	result += "		R2V: begin\n"

	if arch.R == 1 {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "			case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "				" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "					vtm0_wren_i <= 1'b1;\n"
		result += "					vtm0_addr_i[7:0] <= rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-8) + "];\n"
		result += "					vtm0_din_i[7:0] <= _" + strings.ToLower(Get_register_name(i)) + "[7:0];\n"
		result += "					$display(\"R2V " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
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

func (op R2v) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	ramdepth := 8

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

func (op R2v) Disassembler(arch *Arch, instr string) (string, error) {
	ramdepth := 8
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	value := get_id(instr[arch.R : int(arch.R)+ramdepth])
	result += strconv.Itoa(value)
	return result, nil
}

// The simulation does nothing
func (op R2v) Simulate(vm *VM, instr string) error {
	// TODO

	reg_bits := vm.Mach.R
	regPay := get_id(instr[:reg_bits])
	posS := instr[reg_bits : reg_bits+8]

	pos := uint8(get_id(posS))
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
func (op R2v) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2v) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2v) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2v) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2v) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	//result += "\t\t\t\twr_int_ram <= #1 1'b0;\n"
	return result
}

func (Op R2v) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2v) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2v) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "r2v", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op R2v) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2v) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "r2v::*--type=reg::*--type=number")
	result = append(result, "mov::*--type=somov--sotype=vtm--soaddressing=immediate::*--type=reg")
	return result
}
func (Op R2v) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "r2v":
		regNeed := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regVal := line.Elements[1].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regVal, Op: bmreqs.OpAdd})
		port := line.Elements[0].GetMeta("soport")
		if regVal != "" && port != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("r2v")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regVal)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(port)
			newArg1.BasmMeta = newArg1.SetMeta("type", "number")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2v) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
