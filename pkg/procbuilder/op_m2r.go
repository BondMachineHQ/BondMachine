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

type M2r struct{}

func (op M2r) Op_get_name() string {
	return "m2r"
}

func (op M2r) Op_get_desc() string {
	return "Memory to register copy"
}

func (op M2r) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "m2r [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(int(arch.L)) + "(Location)]	// Set a register to the value of a memory location [" + strconv.Itoa(opbits+int(arch.R)+int(arch.L)) + "]\n"
	return result
}

func (op M2r) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.R) + int(arch.L) // The bits for the opcode + bits for a register + bits for memory location
}

func (op M2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\t//Internal Reg Wire for M2R opcode\n"
	result += "\treg state_read_mem;\n"
	result += "\twire [" + strconv.Itoa(int(arch.L)-1) + ":0] addr_ram_m2r;\n"
	result += "\n"
	return result
}

func (Op M2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\tstate_read_mem <= #1 1'b0;\n"
	return result
}

func (Op M2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (Op M2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	return result
}

func (op M2r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""

	result += "					M2R: begin\n"
	result += "						if (state_read_mem == 0) begin\n"
	result += "							state_read_mem <= #1 1'b1;\n"
	result += "						end\n"
	result += "						else begin\n"
	result += "							state_read_mem <= #1 1'b0;\n"
	result += NextInstruction(nil, arch, 7, "_pc + 1'b1")
	if arch.R == 1 {
		result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "							case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "								" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "									_" + strings.ToLower(Get_register_name(i)) + " <= #1 ram_dout;\n"
		result += "									$display(\"M2R " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "								end\n"
	}
	result += "							endcase\n"

	result += "						end\n"
	result += "					end\n"
	return result
}

func (op M2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
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

	result += "\t//logic code to control the address to read RAM\n"
	result += "\tassign addr_ram_m2r = current_instruction[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.L)) + "];\n"

	return result
}

func (op M2r) Assembler(arch *Arch, words []string) (string, error) {
	ramdepth := int(arch.L)
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

func (op M2r) Disassembler(arch *Arch, instr string) (string, error) {
	ramdepth := int(arch.L)
	reg_id := get_id(instr[:arch.R])
	result := strings.ToLower(Get_register_name(reg_id)) + " "
	value := get_id(instr[arch.R : int(arch.R)+ramdepth])
	result += strconv.Itoa(value)
	return result, nil
}

func (op M2r) Simulate(vm *VM, instr string) error {
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	memval := get_id(instr[reg_bits : reg_bits+8])
	vm.Registers[reg] = uint8(memval)
	return nil
}

func (op M2r) Generate(arch *Arch) string {
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	value := rand.Intn(256)
	return zeros_prefix(int(arch.R)+8, get_binary(reg)+get_binary(value))
}

func (op M2r) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op M2r) Required_modes() (bool, []string) {
	return true, []string{}
}

func (op M2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op M2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op M2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "m2r", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op M2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op M2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	result = append(result, "m2r::*--type=reg::*--type=number")
	result = append(result, "mov::*--type=reg::*--type=ram--ramaddressing=immediate")
	return result
}
func (Op M2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "m2r":
		regNeed := line.Elements[0].GetValue()
		location := line.Elements[1].GetMeta("location")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectMax, Name: "ram", Value: location, Op: bmreqs.OpAdd})
		return line, nil
	case "mov":
		regNeed := line.Elements[0].GetValue()
		location := line.Elements[1].GetMeta("location")
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectSet, Name: "registers", Value: regNeed, Op: bmreqs.OpAdd})
		rg.Requirement(bmreqs.ReqRequest{Node: node, T: bmreqs.ObjectMax, Name: "ram", Value: location, Op: bmreqs.OpAdd})
		if regNeed != "" && location != "" {
			newLine := new(bmline.BasmLine)
			newOp := new(bmline.BasmElement)
			newOp.SetValue("m2r")
			newLine.Operation = newOp
			newArgs := make([]*bmline.BasmElement, 2)
			newArg0 := new(bmline.BasmElement)
			newArg0.BasmMeta = newArg0.SetMeta("type", "reg")
			newArg0.SetValue(regNeed)
			newArgs[0] = newArg0
			newArg1 := new(bmline.BasmElement)
			newArg1.SetValue(location)
			newArg1.BasmMeta = newArg1.SetMeta("type", "number")
			newArgs[1] = newArg1
			newLine.Elements = newArgs
			return newLine, nil
		}
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op M2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op M2r) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
