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

type Jcmpl struct{}

func (op Jcmpl) Op_get_name() string {
	return "jcmpl"
}

func (op Jcmpl) Op_get_desc() string {
	return "Jump to a program location conditioned to the comparison flag"
}

func (op Jcmpl) Op_show_assembler(arch *Arch) string {
	opBits := arch.Opcodes_bits()
	result := ""
	switch arch.Modes[0] {
	case "ha":
		result = "jcmpl [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.O)) + "]\n"
	case "vn":
		result = "jcmpl [" + strconv.Itoa(int(arch.L)) + "(Location)]	// Jump to a program location conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.L)) + "]\n"
	case "hy":
		if arch.O > arch.L {
			result = "jcmpl [" + strconv.Itoa(int(arch.O)) + "(Location)]	// Jump to a program location conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.O)) + "]\n"
		} else {
			result = "jcmpl [" + strconv.Itoa(int(arch.L)) + "(Location)]	// Jump to a program location conditioned to the comparison flag [" + strconv.Itoa(opBits+int(arch.L)) + "]\n"
		}
	}
	return result
}

func (op Jcmpl) Op_get_instruction_len(arch *Arch) int {
	opBits := arch.Opcodes_bits()
	switch arch.Modes[0] {
	case "ha":
		return opBits + int(arch.O) // The bits for the opcode + bits for a location
	case "vn":
		return opBits + int(arch.L) // The bits for the opcode + bits for a location
	case "hy":
		if arch.O > arch.L {
			return opBits + int(arch.O) // The bits for the opcode + bits for a location
		} else {
			return opBits + int(arch.L) // The bits for the opcode + bits for a location
		}
	}
	return 0
}

func (op Jcmpl) OpInstructionVerilogHeader(conf *Config, arch *Arch, flavor string, pname string) string {
	result := ""
	if arch.OnlyOne(op.Op_get_name(), []string{"cmpr", "cmprlt", "cmpv", "jcmpl", "jcmpo", "jcmpa", "jcmprio", "jcmpria", "jncmpl", "jncmpo", "jncmpa", "jncmprio", "jncmpria"}) {
		result += "\treg cmpflag;\n"
	}
	return result
}

func (op Jcmpl) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpl) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpl) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpl) Op_instruction_verilog_state_machine(conf *Config, arch *Arch, rg *bmreqs.ReqRoot, flavor string) string {
	romWord := arch.Max_word()
	opBits := arch.Opcodes_bits()

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
	case "hy":
		if arch.O > arch.L {
			locationBits = arch.O
		} else {
			locationBits = arch.L
		}
	}

	result := ""
	result += "					JCMPL: begin\n"
	if locationBits == 1 {
		result += "						if (cmpflag == 1'b1) begin\n"
		result += NextInstruction(conf, arch, 7, "current_instruction["+strconv.Itoa(romWord-opBits-1)+"]")
		result += "						end\n"
		result += "						else begin\n"
		result += NextInstruction(conf, arch, 7, "_pc + 1'b1")
		result += "							end\n"
		result += "						$display(\"JCMPL \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + "]);\n"
	} else {
		result += "						if (cmpflag == 1'b1) begin\n"
		result += NextInstruction(conf, arch, 7, "current_instruction["+strconv.Itoa(romWord-opBits-1)+":"+strconv.Itoa(romWord-opBits-int(locationBits))+"]")
		result += "						end\n"
		result += "						else begin\n"
		result += NextInstruction(conf, arch, 7, "_pc + 1'b1")
		result += "						end\n"
		result += "						$display(\"JCMPL \", current_instruction[" + strconv.Itoa(romWord-opBits-1) + ":" + strconv.Itoa(romWord-opBits-int(locationBits)) + "]);\n"
	}
	result += "					end\n"
	return result
}

func (op Jcmpl) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	return ""
}

func (op Jcmpl) Assembler(arch *Arch, words []string) (string, error) {
	opBits := arch.Opcodes_bits()
	romWord := arch.Max_word()

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
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

func (op Jcmpl) Disassembler(arch *Arch, instr string) (string, error) {

	locationBits := arch.O

	switch arch.Modes[0] {
	case "ha":
		locationBits = arch.O
	case "vn":
		locationBits = arch.L
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

func (op Jcmpl) Simulate(vm *VM, instr string) error {
	// TODO
	value := get_id(instr[:vm.Mach.O])
	if value < len(vm.Mach.Slocs) {
		vm.Pc = uint64(value)
	} else {
		vm.Pc = vm.Pc + 1
	}
	return nil
}

func (op Jcmpl) Generate(arch *Arch) string {
	max_value := 1 << arch.O
	value := rand.Intn(max_value)
	return zeros_prefix(int(arch.O), get_binary(value))
}

func (op Jcmpl) Required_shared() (bool, []string) {
	return false, []string{}
}

func (op Jcmpl) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Jcmpl) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Jcmpl) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jcmpl) AbstractAssembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "j", I_NIL}
	result[0] = newnot
	return result, nil
}

func (op Jcmpl) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
func (op Jcmpl) HLAssemblerMatch(arch *Arch) []string {
	result := make([]string, 4)
	result[0] = "jcmpl::*--type=number"
	result[1] = "jcmpl::*--type=symbol"
	result[2] = "jcmp::*--type=number"
	result[3] = "jcmp::*--type=symbol"
	return result
}
func (op Jcmpl) HLAssemblerNormalize(arch *Arch, rg *bmreqs.ReqRoot, node string, line *bmline.BasmLine) (*bmline.BasmLine, error) {
	switch line.Operation.GetValue() {
	case "jcmp":
		line.Operation.SetValue("jcmpl")
		return line, nil
	case "jcmpl":
		return line, nil
	}
	return nil, errors.New("HL Assembly normalize failed")
}
func (op Jcmpl) ExtraFiles(arch *Arch) ([]string, []string) {
	return []string{}, []string{}
}

func (op Jcmpl) HLAssemblerInstructionMetadata(arch *Arch, line *bmline.BasmLine) (*bmmeta.BasmMeta, error) {
	return nil, nil
}
