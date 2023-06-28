package procbuilder

// TODO This is the ROM, change it to halndle also the RAM case

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmmeta"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

type Jc struct{}

func (op Jc) Op_get_name() string {
	return "jc"
}

func (op Jc) Op_get_desc() string {
	return "Jump to a program location if carry-bit is 0"
}

func (op Jc) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := "jc [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jcump to a program location [" + strconv.Itoa(opbits+int(arch.O)) + "]\n"
	return result
}

func (op Jc) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	return opbits + int(arch.O) // The bits for the opcode + bits for a location
}

func (op Jc) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	setflag := conf.Runinfo.Check("carryflag")

	if setflag {
		result += "\treg carryflag;\n"
	}

	return result
}

func (Op Jc) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Jc) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Jc) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jc) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	result := ""
	result += "					JC: begin\n"
	if arch.O == 1 {
		result += "					if(carryflag == 'b0) begin\n"
		result += "						_pc <= #1 rom_value[" + strconv.Itoa(rom_word-opbits-1) + "];\n"
		result += "						$display(\"JC \", rom_value[" + strconv.Itoa(rom_word-opbits-1) + "]);\n"
		result += " end \n"
	} else {
		result += "					if(carryflag == 'b0) begin\n"
		result += "						_pc <= #1 rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.O)) + "];\n"
		result += "						$display(\"JC \", rom_value[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(arch.O)) + "]);\n"
		result += " end \n"
	}
	result += "					end\n"
	return result
}

func (op Jc) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Jc) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	if len(words) != 1 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	if partial, err := Process_number(words[0]); err == nil {
		result += zeros_prefix(int(arch.O), partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(arch.O); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Jc) Disassembler(arch *Arch, instr string) (string, error) {
	value := get_id(instr[:arch.O])
	result := strconv.Itoa(value)
	return result, nil
}

func (op Jc) Simulate(vm *VM, instr string) error {
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op Jc) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Jc) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Jc) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Jc) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (Op Jc) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Jc) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "jc", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op Jc) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (Op Jc) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 0)
	return result
}
func (Op Jc) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	return nil, errors.New("HL Assembly normalize failed")
}
func (Op Jc) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Jc) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
