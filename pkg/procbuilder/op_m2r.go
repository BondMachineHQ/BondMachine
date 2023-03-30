package procbuilder

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
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
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	reg_num := 1 << arch.R

	result := ""
	result += "\t\t\tif(state_read_mem) begin\n"

	if arch.R == 1 {
		result += "				case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "				case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "					" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "						_" + strings.ToLower(Get_register_name(i)) + " <= #1 ram_dout;\n"
		result += "						state_read_mem <= #1 1'b0;\n"
		result += "						$display(\"M2R " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "					end\n"
	}
	result += "				endcase\n"
	result += "\t\t\t\t_pc <= #1 _pc + 1'b1;\n"
	result += "\t\t\tend\n"

	return result
}

func (op M2r) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	result := ""
	result += "					M2R: begin\n"
	result += "						state_read_mem <= #1 1'b1;\n"
	result += "					end\n"
	return result
}

func (op M2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	result := ""
	result += "\t//logic code to control the address to read RAM\n"
	result += "\tassign addr_ram_m2r = rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-int(arch.L)) + "];\n"

	setflag := true
	for _, currop := range arch.Op {
		if currop.Op_get_name() == "r2m" {
			setflag = false
			break
		} else if currop.Op_get_name() == "m2r" {
			break
		}
	}
	if setflag {
		result += "\tassign ram_din = ram_din_i;\n"
		result += "\tassign ram_addr = (rom_value[" + strconv.Itoa(rom_word-1) + ":" + strconv.Itoa(rom_word-opbits) + "]==M2R) ? addr_ram_m2r : addr_ram_r2m;\n"
		result += "\tassign ram_wren = wr_int_ram;\n"
		result += "\tassign ram_en = 1'b1;\n"
	}

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

func (Op M2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\t\tstate_read_mem <= #1 1'b0;\n"
	return result
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
	return result
}
func (Op M2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op M2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
