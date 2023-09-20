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

type Jcmpo struct{}

func (op Jcmpo) Op_get_name() string {
	return "jcmpo"
}

func (op Jcmpo) Op_get_desc() string {
	return "Jump to a program location in the ROM conditioned to the comparison flag"
}

func (op Jcmpo) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := ""
	switch arch.Modes[0] {
	case "ha":
		result = "jcmpo [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.O)) + "]\n"
	case "hy":
		if arch.O > arch.L {
			result = "jcmpo [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location in ROM conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.O)) + "]\n"
		} else {
			result = "jcmpo [" + strconv.Itoa(int(arch.L)) + "(Location)]	// Jump to a program location in ROM conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.L)) + "]\n"
		}
	}
	return result
}

func (op Jcmpo) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	switch arch.Modes[0] {
	case "ha":
		return opBits + int(arch.O) // The bits for the opcode + bits for a location
	case "hy":
		if arch.O > arch.L {
			return opBits + int(arch.O) // The bits for the opcode + bits for a location
		} else {
			return opBits + int(arch.L) // The bits for the opcode + bits for a location
		}
	}
	return 0
}

func (op Jcmpo) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"cmpr", "jcmpl", "jcmpo", "jcmpa", "jcmprio", "jcmpria"}) {
		result += "\treg cmpflag;\n"
	}
	return result
}

func (op Jcmpo) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpo) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpo) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpo) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()

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
	result += "					JCMPO: begin\n"
	if locationBits == 1 {
		result += "						if (cmpflag == 1'b1) begin\n"
		if arch.Modes[0] == "hy" {
			result += "							exec_mode <= #1 1'b0;\n"
		}

		result += "							_pc <= current_instruction[" + strconv.Itoa(romWord-opBits-1) + "];\n"
		result += "						end\n"
		result += "						else begin\n"
		result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
		result += "						end\n"
		result += "						$display(\"JCMPL \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + "]);\n"
	} else {
		result += "						if (cmpflag == 1'b1) begin\n"
		if arch.Modes[0] == "hy" {
			result += "							exec_mode <= #1 1'b0;\n"
		}

		result += "							_pc <= current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "];\n"
		result += "						end\n"
		result += "						else begin\n"
		result += NextInstruction(conf, arch, 6, "_pc + 1'b1")
		result += "						end\n"
		result += "						$display(\"JCMPL \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "]);\n"
	}
	result += "					end\n"
	return result
}

func (op Jcmpo) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpo) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

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

	for i := opBits + int(locationBits); i < romWord; i++ {
		result += "0"
	}

	return result, nil
}

func (op Jcmpo) Disassembler(arch *Arch, instr string) (string, error) {

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

func (op Jcmpo) Simulate(vm *VM, instr string) error {
	// TODO
	return nil
}

func (op Jcmpo) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Jcmpo) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Jcmpo) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Jcmpo) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Jcmpo) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jcmpo) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "j", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op Jcmpo) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Jcmpo) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 3)
	result[0] = "jcmpo::*--type=number"
	result[1] = "jcmpo::*--type=rom--romaddressing=symbol"
	result[2] = "jcmp::*--type=rom--romaddressing=symbol"
	return result
}
func (op Jcmpo) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "jcmp":
		line.Operation.SetValue("jcmpo")
		return line, nil
	case "jcmpo":
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Jcmpo) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jcmpo) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
