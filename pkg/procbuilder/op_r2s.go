package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The R2s opcode is both a basic instruction and a template for other instructions.
type R2s struct{}

func (op R2s) Op_get_name() string {
	// TODO
	return "r2s"
}

func (op R2s) Op_get_desc() string {
	// TODO
	return "No operation"
}

func (op R2s) Op_show_assembler(arch *Arch) string {
	// TODO
	opbits := arch.Opcodes_bits()
	result := "r2s [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op R2s) Op_get_instruction_len(arch *Arch) int {
	// TODO
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op R2s) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	result += "\treg [" + strconv.Itoa(int(arch.L)-1) + ":0] sh_addr_r2s;\n"
	return result
}

func (Op R2s) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	result := ""

	sharemem_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "sharedmem" {
				sharemem_num++
			}
		}
	}

	result += "\t\t\tsh_wren_i <= #1 'b0;\n"
	for sh_num := 0; sh_num < sharemem_num; sh_num++ {
		result += "\t\t\tsh_din_i[" + strconv.Itoa(sh_num) + "] <= #1 'b0;\n"
		result += "\t\t\tsh_wren_i[" + strconv.Itoa(sh_num) + "] <= #1 'b0;\n"
	}
	result += "\t\t\tsh_addr_r2s <= #1 'b0;\n"
	return result
}

func (Op R2s) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	result := ""
	result += "\t\t\t\tsh_wren_i <= #1 'b0;\n"
	return result
}

func (op R2s) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {

	result := ""

	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	shso := Sharedmem{}
	shbits := arch.Shared_bits(shso.Shr_get_name()) //CHECK

	reg_num := 1 << arch.R

	result += "					R2S: begin\n"
	if arch.R == 1 {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "					case (rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}
	for i := 0; i < reg_num; i++ {

		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"
		result += "							sh_wren_i[rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits) + "]] <= #1 1'b1;\n"
		result += "							sh_addr_r2s <= #1 rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits-int(arch.Rsize)) + "];\n"
		result += "							sh_din_i[rom_value[" + strconv.Itoa(rom_word-opbits-int(arch.R)-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)-shbits) + "]] <= #1 _" + strings.ToLower(Get_register_name(i)) + ";\n"
		result += "							$display(\"R2S " + strings.ToUpper(Get_register_name(i)) + " \",_" + strings.ToLower(Get_register_name(i)) + ");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += "						_pc <= _pc + 1'b1 ;\n"
	result += "					end\n"

	return result
}

func (op R2s) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op R2s) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op R2s) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op R2s) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op R2s) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op R2s) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op R2s) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op R2s) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op R2s) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op R2s) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op R2s) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op R2s) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op R2s) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op R2s) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op R2s) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
