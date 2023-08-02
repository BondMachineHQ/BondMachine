package procbuilder

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

// The Hit opcode is both a basic instruction and a template for other instructions.
type Hit struct{}

func (op Hit) Op_get_name() string {
	return "hit"
}

func (op Hit) Op_get_desc() string {
	return "Hit a barrier"
}

func (op Hit) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	barbits := arch.Shared_bits("barrier")
	result := "hit [" + strconv.Itoa(int(arch.R)) + "(Reg)] [" + strconv.Itoa(barbits) + "(Barrier)]       // Hit a barrier [" + strconv.Itoa(opbits+int(arch.R)+barbits) + "]\n"
	return result
}

func (op Hit) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	barbits := arch.Shared_bits("barrier")
	return opbits + int(arch.R) + int(barbits) // The bits for the opcode + bits for a register + bits for the barrier id
}

func (op Hit) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	// TODO
	return ""
}

func (Op Hit) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (Op Hit) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (Op Hit) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Hit) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()
	bar_num := arch.Shared_num("barrier")

	reg_num := 1 << arch.R

	result := ""
	result += "					HIT: begin\n"
	if arch.R == 1 {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "])\n"
	} else {
		result += "						case (current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.R)) + "])\n"
	}

	for i := 0; i < reg_num; i++ {
		result += "						" + strings.ToUpper(Get_register_name(i)) + " : begin\n"

		for j := 0; j < bar_num; j++ {
			result += "						" + strings.ToUpper(Get_register_name(j)) + " : begin\n"

			result += "							$display(\"CLR " + strings.ToUpper(Get_register_name(j)) + "\");\n"
			result += "						end\n"
		}

		result += "							$display(\"CLR " + strings.ToUpper(Get_register_name(i)) + "\");\n"
		result += "						end\n"
	}
	result += "						endcase\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Hit) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Hit) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Hit) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op Hit) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Hit) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Hit) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Hit) Required_modes() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Hit) Forbidden_modes() (bool, []string) {
	// TODO
	return false, []string{}
}

func (Op Hit) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Hit) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Hit) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Hit) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Hit) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Hit) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Hit) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
