package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The S2r opcode is both a basic instruction and a template for other instructions.
type S2r struct{}

func (op S2r) Op_get_name() string {
	// TODO
	return "s2r"
}

func (op S2r) Op_get_desc() string {
	// TODO
	return "No operation"
}

func (op S2r) Op_show_assembler(arch *Arch) string {
	// TODO
	opbits := arch.Opcodes_bits()
	result := "s2r [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op S2r) Op_get_instruction_len(arch *Arch) int {
	// TODO
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op S2r) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\t//Internal Reg Wire for S2R opcode\n"
	result += "\treg state_sh_read_mem;\n"
	result += "\twire [" + strconv.Itoa(int(arch.L)-1) + ":0] sh_addr_s2r;\n"
	result += "\n"
	return result
}

func (Op S2r) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\tstate_sh_read_mem <= #1 1'b0;\n"
	return result
}

func (Op S2r) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {

	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	shso := Sharedmem{}
	shbits := arch.Shared_bits(shso.Shr_get_name()) //CHECK

	reg_num := 1 << arch.R

	result += "\t\t\tif(state_sh_read_mem) begin\n"

	if arch.R == 1 {
		result += "				case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "				case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {
		result += "					" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "						_" + strings.ToLower(Get_register_name(i)) + " <= #1 sh_dout_i[rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits) + "]];\n"
		result += "						state_sh_read_mem <= #1 1'b0;\n"
		result += "						$display(\"S2R " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "					end\n"
	}
	result += "				endcase\n"
	result += "\t\t\t\t_pc <= #1 _pc + 1'b1;\n"
	result += "\t\t\tend\n"

	return result
}

func (op S2r) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	result := ""
	result += "					S2R: begin\n"
	result += "						state_sh_read_mem <= #1 1'b1;\n"
	result += "					end\n"
	return result
}

func (op S2r) Op_instruction_verilog_footer(arch *Arch, flavor string) string {

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	shso := Sharedmem{}
	shbits := arch.Shared_bits(shso.Shr_get_name()) //CHECK

	result := ""
	result += "\t//logic code to control the address to read RAM\n"
	result += "\tassign sh_addr_s2r = rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits-int(arch.Rsize)) + "];\n"

	return result
}

func (op S2r) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op S2r) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op S2r) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op S2r) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op S2r) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op S2r) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op S2r) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op S2r) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op S2r) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op S2r) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op S2r) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op S2r) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op S2r) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op S2r) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
