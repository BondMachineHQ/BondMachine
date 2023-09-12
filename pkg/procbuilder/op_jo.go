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

type Jo struct{}

func (op Jo) Op_get_name() string {
	return "jo"
}

func (op Jo) Op_get_desc() string {
	return "Jump to a program location in the ROM"
}

func (op Jo) Op_show_assembler(arch *Arch) string {
	opbits := arch.Opcodes_bits()
	result := ""
	switch arch.Modes[0] {
	case "ha":
		result = "j [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location [" + strconv.Itoa(opbits+int(arch.O)) + "]\n"
	case "hy":
		if arch.O > arch.L {
			result = "j [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location in ROM [" + strconv.Itoa(opbits+int(arch.O)) + "]\n"
		} else {
			result = "j [" + strconv.Itoa(int(arch.L)) + "(Location)]	// Jump to a program location in ROM [" + strconv.Itoa(opbits+int(arch.L)) + "]\n"
		}
	}
	return result
}

func (op Jo) Op_get_instruction_len(arch *Arch) int {
	opbits := arch.Opcodes_bits()
	switch arch.Modes[0] {
	case "ha":
		return opbits + int(arch.O) // The bits for the opcode + bits for a location
	case "hy":
		if arch.O > arch.L {
			return opbits + int(arch.O) // The bits for the opcode + bits for a location
		} else {
			return opbits + int(arch.L) // The bits for the opcode + bits for a location
		}
	}
	return 0
}

func (op Jo) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	return ""
}

func (op Jo) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Jo) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jo) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jo) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	rom_word := arch.Max_word()
	opbits := arch.Opcodes_bits()

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	result := ""
	result += "					JO: begin\n"
	if arch.Modes[0] == "hy" {
		result += "						exec_mode <= #1 1'b0;\n"
	}

	if locationBits == 1 {
		result += "						_pc <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "];\n"
		result += "						$display(\"JO \", current_instruction[" + strconv.Itoa(rom_word-opbits-1) + "]);\n"
	} else {
		result += "						_pc <= #1 current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(locationBits)) + "];\n"
		result += "						$display(\"JO \", current_instruction[" + strconv.Itoa(rom_word-opbits-1) + ":" + strconv.Itoa(rom_word-opbits-int(locationBits)) + "]);\n"
	}
	result += "					end\n"
	return result
}

func (op Jo) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Jo) Assembler(arch *Arch, words []string) (string, error) {
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	if len(words) != 1 {
		return "", Prerror{"Wrong arguments number"}
	}

	result := ""
	if partial, err := Process_number(words[0]); err == nil {
		result += zeros_prefix(int(locationBits), partial)
	} else {
		return "", Prerror{err.Error()}
	}

	for i := opbits + int(locationBits); i < rom_word; i++ {
		result += "0"
	}

	return result, nil
}

func (op Jo) Disassembler(arch *Arch, instr string) (string, error) {

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	value := get_id(instr[:locationBits])
	result := strconv.Itoa(value)
	return result, nil
}

func (op Jo) Simulate(vm *VM, instr string) error {
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op Jo) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Jo) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Jo) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Jo) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Jo) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jo) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "j", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op Jo) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Jo) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 3)
	result[0] = "jo::*--type=number"
	result[1] = "jo::*--type=rom--romaddressing=symbol"
	result[2] = "jmp::*--type=rom--romaddressing=symbol"
	return result
}
func (op Jo) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "jmp":
		line.Operation.SetValue("jo")
		return line, nil
	case "jo":
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Jo) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jo) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
