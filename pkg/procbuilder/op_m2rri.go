package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type M2rri struct{}

func (op M2rri) Op_get_name() string {
	return "m2rri"
}

func (op M2rri) Op_get_desc() string {
	return "ROM to register"
}

func (op M2rri) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "m2rri [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Set a register to the value of the given ROM location [" + strconv.Itoa(opbits+int(arch.R)+int(arch.O)) + "]\n"
	return result
}

func (op M2rri) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.O) // The bits for the opcode + bits for a register + bits for the location
}

func (op M2rri) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\t//Internal Reg Wire for M2R opcode\n"
	result += "\treg state_read_mem_m2rri;\n"
	result += "\treg wait_read_mem;\n"
	result += "\treg [" + strconv.Itoa(int(arch.L)-1) + ":0] addr_ram_m2rri;\n"
	result += "\n"
	return result
}

func (Op M2rri) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	// result += "\t\t\tstate_read_mem_m2rri <= #1 1'b0;\n"
	return result
}

func (Op M2rri) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	// rom_word := arch.Max_word()
	// opbits := arch.Opcodes_bits()

	// reg_num := 1 << arch.R

	result := ""
	// result += "\t\t\tif(state_read_mem_m2rri) begin\n"

	// if arch.R == 1 {
	// 	result += "				case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	// } else {
	// 	result += "				case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	// }
	// for i := 0; i < reg_num; i++ {
	// 	result += "					" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
	// 	result += "						_" + strings.ToLower(Get_register_name(i)) + " <= #1 ram_dout;\n"
	// 	result += "						state_read_mem_m2rri <= #1 1'b0;\n"
	// 	result += "						$display(\"M2RRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
	// 	result += "					end\n"
	// }
	// result += "				endcase\n"
	// result += NextInstruction(conf, arch, 4, "_pc + 1'b1")
	// result += "\t\t\tend\n"

	return result
}

func (Op M2rri) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}
func (op M2rri) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()
	regNum := 1 << arch.R

	result := ""
	result += "					M2RRI: begin\n"
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

			result += "								if (state_read_mem_m2rri == 1'b1) begin\n"
			result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 ram_dout;\n"
			result += "									state_read_mem_m2rri <= 1'b0;\n"
			result += NextInstruction(conf, arch, 9, "_pc + 1'b1")
			result += "								end\n"
			result += "								else begin\n"
			result += "									if (wait_read_mem == 1'b1) begin\n"
			result += "										state_read_mem_m2rri <= 1'b1;\n"
			result += "										wait_read_mem <= 1'b0;\n"
			result += "									end\n"
			result += "									else begin\n"
			result += "										wait_read_mem <= 1'b1;\n"
			result += "										addr_ram_m2rri <= #1 _" + strings.ToLower(Get_register_name(j)) + ";\n"
			result += "									end\n"
			result += "								end\n"
			result += "								$display(\"M2RRI " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(j)) + ");\n"

			result += "							end\n"

		}
		result += "							endcase\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "					end\n"
	return result

}

func (op M2rri) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	// The ram is enabled if any of the opcodes is active
	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m", "m2r", "m2rri"}) {
		result += "\tassign ram_en = 1'b1;\n"
	}

	// The ram write signals
	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m"}) {
		result += "\tassign ram_din = ram_din_i;\n"
		result += "\tassign ram_wren = wr_int_ram;\n"
	}

	ramAddr := ""
	if arch.HasAny([]string{"r2mri", "r2m"}) {
		ramAddr += "addr_ram_to_mem"
	}

	if arch.HasOp("m2r") {
		ramAddr = " (current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2R) ? addr_ram_m2r : " + ramAddr
	}

	if arch.HasOp("m2rri") {
		ramAddr = " (current_instruction[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2RRI) ? addr_ram_m2rri: " + ramAddr
	}

	if arch.Modes[0] == "hy" || arch.Modes[0] == "vn" {
		ramAddr = " (exec_mode == 1'b1 && vn_state == FETCH) ? _pc : " + ramAddr
	}

	if arch.OnlyOne(op.Op_get_name(), []string{"r2mri", "r2m", "m2r", "m2rri"}) {
		result += "\tassign ram_addr = " + ramAddr + ";\n"
	}

	return result
}

func (op M2rri) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	regNum := 1 << arch.R

	if len(words) != 2 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	for i := 0; i < regNum; i++ {
		if words[0] == strings.ToLower(Get_register_name(i)) {
			result += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if result == "" {
		return "", Prerror{"Unknown register name " + words[0]}
	}

	partial := ""
	for i := 0; i < regNum; i++ {
		if words[1] == strings.ToLower(Get_register_name(i)) {
			partial += zeros_prefix(int(arch.R), get_binary(i))
			break
		}
	}

	if partial == "" {
		return "", Prerror{"Unknown register name " + words[1]}
	}

	result += partial

	for i := opbits + 2*int(arch.R); i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op M2rri) Disassembler(arch *Arch, instr string) (string, error) {
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	reg_id = get_id(instr[arch.R : 2*int(arch.R)])
	result += strings.ToLower(Get_register_name(reg_id))
	return result, nil
}

func (op M2rri) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	regDest := get_id(instr[:reg_bits])
	regSrc := get_id(instr[reg_bits : reg_bits*2])

	if vm.Mach.Rsize <= 8 {
		loc := int(vm.Registers[regSrc].(uint8))
		if loc < len(vm.Mach.Program.Slocs) {
			vm.Registers[regDest] = uint8(get_id(vm.Mach.Program.Slocs[loc]))
		} else {
			vm.Registers[regDest] = uint8(get_id(vm.Mach.Data.Vars[loc-len(vm.Mach.Program.Slocs)]))
		}
	} else if vm.Mach.Rsize <= 16 {
		loc := int(vm.Registers[regSrc].(uint16))
		if loc < len(vm.Mach.Program.Slocs) {
			vm.Registers[regDest] = uint16(get_id(vm.Mach.Program.Slocs[loc]))
		} else {
			vm.Registers[regDest] = uint16(get_id(vm.Mach.Data.Vars[loc-len(vm.Mach.Program.Slocs)]))
		}
	} else if vm.Mach.Rsize <= 32 {
		loc := int(vm.Registers[regSrc].(uint32))
		if loc < len(vm.Mach.Program.Slocs) {
			vm.Registers[regDest] = uint32(get_id(vm.Mach.Program.Slocs[loc]))
		} else {
			vm.Registers[regDest] = uint32(get_id(vm.Mach.Data.Vars[loc-len(vm.Mach.Program.Slocs)]))
		}
	} else if vm.Mach.Rsize <= 64 {
		loc := int(vm.Registers[regSrc].(uint64))
		if loc < len(vm.Mach.Program.Slocs) {
			vm.Registers[regDest] = uint64(get_id(vm.Mach.Program.Slocs[loc]))
		} else {
			vm.Registers[regDest] = uint64(get_id(vm.Mach.Data.Vars[loc-len(vm.Mach.Program.Slocs)]))
		}
	} else {
		return errors.New("invalid register size, must be <= 64")
	}

	vm.Pc = vm.Pc + 1
	return nil
}

func (op M2rri) Generate(arch *Arch) string {
	return ""
}

func (op M2rri) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op M2rri) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op M2rri) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op M2rri) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op M2rri) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	seq0, types0 := Sequence_to_0(words[0])
	seq1, types1 := Sequence_to_0(words[1])

	if len(seq0) > 0 && types0 == O_REGISTER && len(seq1) > 0 && types1 == O_INPUT {

		result := make([]UsageNotify, 2+len(seq1))
		newnot0 := UsageNotify{C_OPCODE, "m2rri", I_NIL}
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

func (Op M2rri) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op M2rri) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "m2rri::*--type=reg::*--type=reg")
	result = append(result, "mov::*--type=reg::*--type=ram--ramaddressing=register")
	return result
}
func (Op M2rri) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "m2rri":
		regNeed := line.Elements[1].GetValue()
		regDest := line.Elements[0].GetValue()
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDest, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regDest := line.Elements[0].GetValue()
		regNeed := line.Elements[1].GetMeta("ramregister")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regDest, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		if regDest != "" && regNeed != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("m2rri")
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
func (Op M2rri) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op M2rri) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
