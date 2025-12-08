package procbuilder

import (
	"errors"
	"math/rand"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Clc struct{}

func (op Clc) Op_get_name() string {
	return "clc"
}

// Op_get_desc returns the description of the Clc operation which is "Clear carry-bit".
func (op Clc) Op_get_desc() string {
	return "Clear carry-bit"
}

func (op Clc) Op_show_assembler(arch *Arch) string {
	result := "clc	// Set carry-bit to 0 \n"
	return result
}

func (op Clc) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits // The bits for the opcode + bits for a register
}

func (op Clc) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	setflag := conf.Runinfo.Check("carryflag")

	if setflag {
		result += "\treg carryflag;\n"
	}

	return result
}

func (Op Clc) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}
func (Op Clc) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Clc) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Clc) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	result := ""
	result += "					CLC: begin\n"
	result += "						carryflag <= #1 'b0;\n"
	result += "						$display(\"CLC\");\n"
	result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
	result += "					end\n"
	return result
}

func (op Clc) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Clc) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Clc) Disassembler(arch *Arch, instr string) (string, error) {
	return "", nil
}

func (op Clc) Simulate(vm *VM, instr string) error {
	//TODO FARE
	reg_bits := vm.Mach.R
	reg := get_id(instr[:reg_bits])
	vm.Registers[reg] = 0
	vm.Pc = vm.Pc + 1
	return nil
}

func (op Clc) Generate(arch *Arch) string {
	//TODO FARE
	reg_num := 1 << arch.R
	reg := rand.Intn(reg_num)
	return zeros_prefix(int(arch.R), get_binary(reg))
}

func (op Clc) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Clc) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Clc) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Clc) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Clc) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot0 := UsageNotify{C_OPCODE, "clc", I_NIL}
	result[0] = newnot0
	return result, nil
}

func (Op Clc) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Clc) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Clc) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Clc) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Clc) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
