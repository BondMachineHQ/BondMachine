package procbuilder

import (
	"errors"
	"math/rand"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Cset struct{}

func (op Cset) Op_get_name() string {
	return "cset"
}

func (op Cset) Op_get_desc() string {
	return "Set carry-bit to 1"
}

func (op Cset) Op_show_assembler(arch *Arch) string {
	result := "cset	// Set carry-bit to 1 \n"
	return result
}

func (op Cset) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits // The bits for the opcode + bits for a register
}

func (op Cset) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	setflag := conf.Runinfo.Check("carryflag")

	if setflag {
		result += "\treg carryflag;\n"
	}

	return result
}

func (Op Cset) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}
func (Op Cset) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Cset) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Cset) Op_instruction_verilog_state_machine(arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	result := ""
	result += "					CSET: begin\n"
	result += "						carryflag <= #1 'b1;\n"
	result += "						$display(\"CSET\");\n"
	result += "						_pc <= #1 _pc + 1'b1 ;\n"
	result += "					end\n"
	return result
}

func (op Cset) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Cset) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Cset) Disassembler(arch *Arch, instr string) (string, error) {
	return "", nil
}

func (op Cset) Simulate(vm *VM, instr string) error {
	//TODO FARE
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	vm.Registers[reg] = 0
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Cset) Generate(arch *Arch) string {
	//TODO FARE
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	return zeros_prefix(int(arch.R), get_binary(reg))
}

func (op Cset) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Cset) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Cset) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Cset) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Cset) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot0 := UsageNotify{C_OPCODE, "cset", I_NIL}
	result[0] = newnot0
	return result, nil
}

func (Op Cset) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Cset) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Cset) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Cset) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}
